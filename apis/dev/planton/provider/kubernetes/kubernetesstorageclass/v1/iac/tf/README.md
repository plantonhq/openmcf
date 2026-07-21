# Kubernetes Storage Class - Terraform Module

## Overview

This Terraform module creates and manages a Kubernetes `storage.k8s.io/v1` StorageClass. It supports the StorageClass surface: provisioner, provisioner-specific parameters, reclaim policy, volume binding mode, volume expansion, mount options, a single allowed-topologies term, and the default-class marker.

## Architecture

```
iac/tf/
├── provider.tf     # Terraform and Kubernetes provider requirements
├── variables.tf    # Input variables mirroring spec.proto
├── locals.tf       # Derived values: labels, default-class annotation, resolved policies
├── main.tf         # Creates kubernetes_storage_class_v1 resource
├── outputs.tf      # Exports storage_class_name, provisioner, is_default_class
└── README.md       # This file
```

## How It Works

1. **Variable Input**: The `spec` variable mirrors the protobuf schema; enum fields arrive as the proto enum value names (e.g. `"delete"`, `"retain"`, `"immediate"`, `"wait_for_first_consumer"`)
2. **Label Merging**: Standard Planton labels are merged with user labels (identity keys cannot be overridden)
3. **Default-Class Annotation**: `locals.tf` renders `is_default_class: true` as the `storageclass.kubernetes.io/is-default-class: "true"` annotation — the spec field is the source of truth, identical to the Pulumi module
4. **Policy Resolution**: `locals.tf` maps the spec's enum values to the API strings (`delete` → `Delete`, `wait_for_first_consumer` → `WaitForFirstConsumer`) with the API server's own defaults applied. The resolved values are ALWAYS sent explicitly, mirroring the Pulumi module, so both engines submit byte-identical objects for the same manifest
5. **Resource Creation**: `main.tf` creates a single `kubernetes_storage_class_v1` resource; the provider forces replacement when the immutable `provisioner` or `parameters` change, matching the Pulumi module's delete-before-replace semantics (StorageClass names are cluster-unique, so delete-then-create is the only safe order)
6. **Output Export**: Class name, provisioner, and the default-class flag are exported

## Semantics Preserved by the Module

- **Policies are always sent explicitly** — resolved with the API defaults in locals rather than left to the provider, keeping both engines' submitted objects identical
- **The default-class marker is the annotation wire form** — never a hand-written annotation in the spec's annotation map
- **Empty collections are omitted** — `parameters` and `mount_options` are sent only when non-empty

> **PARITY-EXCEPTION**: the Terraform Kubernetes provider models `allowed_topologies` as a SINGLE selector term (`max_items = 1`), while the API — and the Pulumi module — accept multiple OR'd terms. One term passes through intact; this module **fails the plan with a precondition** when the spec lists several: deploy that class with the Pulumi provisioner, or combine the zone values into a single term (values within one requirement already OR together). Failing loudly beats silently dropping terms.

## Usage

```hcl
module "storage_class" {
  source = "./iac/tf"

  metadata = {
    name = "fast-ssd"
  }

  spec = {
    name        = "fast-ssd"
    provisioner = "ebs.csi.aws.com"

    parameters = {
      type      = "gp3"
      encrypted = "true"
    }

    volume_binding_mode    = "wait_for_first_consumer"
    allow_volume_expansion = true
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `metadata` | Resource metadata (name, org, env) | object | yes |
| `spec` | StorageClass specification | object | yes |

## Outputs

| Name | Description |
|------|-------------|
| `storage_class_name` | Name of the StorageClass object as created in the cluster |
| `provisioner` | The provisioner (CSI driver) backing this class |
| `is_default_class` | Whether this class is annotated as the cluster's default StorageClass |

> **Note**: The class object is creatable on any cluster — parameters are consumed by the CSI driver at provision time, not validated at class creation. A claim under a `wait_for_first_consumer` class stays Pending until a pod uses it — correct behavior, not an error.
