# Kubernetes Persistent Volume Claim - Terraform Module

## Overview

This Terraform module creates and manages a Kubernetes `core/v1` PersistentVolumeClaim. It supports the PersistentVolumeClaim surface expressible in the Terraform Kubernetes provider: access modes, storage requests and limits, StorageClass selection (including the empty-vs-absent distinction), volume mode, static binding to a named volume, and volume selectors. Data sources (clone/snapshot-restore) are not expressible in the provider and are rejected at plan time (see parity note).

## Architecture

```
iac/tf/
├── provider.tf     # Terraform and Kubernetes provider requirements
├── variables.tf    # Input variables mirroring spec.proto
├── locals.tf       # Derived values: labels, namespace default, access modes, resolved class name
├── main.tf         # Creates kubernetes_persistent_volume_claim_v1 resource
├── outputs.tf      # Exports pvc_name, namespace, storage_request
└── README.md       # This file
```

## How It Works

1. **Variable Input**: The `spec` variable mirrors the protobuf schema; the namespace and `storage_class_name` `StringValueOrRef` fields arrive flattened to plain strings, and enum fields arrive as the proto enum value names (e.g. `"filesystem"`, `"block"`)
2. **Namespace Default**: `locals.tf` falls back to the `default` namespace when none is provided
3. **Label Merging**: Standard Planton labels are merged with user labels (identity keys cannot be overridden)
4. **Access Modes and Volume Mode**: the Kubernetes-default access mode (`ReadWriteOnce`) is applied module-side when the spec omits the field — the API itself REQUIRES `accessModes` — and `volume_mode` is mapped to its API string (`Filesystem` default). Both are ALWAYS sent explicitly, mirroring the Pulumi module, so both engines submit identical claims for the same manifest
5. **Class Name Resolution**: `locals.tf` computes the three-valued `storageClassName`: the class name, `""` when `disable_dynamic_provisioning` is set, or `null` (unset — cluster default applies). The empty-vs-absent distinction is load-bearing upstream
6. **Resource Creation**: `main.tf` creates a single `kubernetes_persistent_volume_claim_v1` resource with `wait_until_bound = false`
7. **Output Export**: Claim name, namespace, and storage request are exported

## Semantics Preserved by the Module

- **The apply NEVER waits for the claim to bind** — `wait_until_bound = false` is set deliberately (mirroring the Pulumi module's skipAwait annotation): under a `WaitForFirstConsumer` StorageClass a claim is correctly Pending until a pod consumes it, and the provider's default wait would hang every such apply. The attribute exists only in configuration — it is declared config-only in the provider import catalog
- **Empty-vs-absent class name is preserved** — `null` means absent (cluster default class); `""` means explicitly no dynamic provisioning (bind only pre-provisioned volumes)
- **Defaults are sent explicitly** — access modes and volume mode resolved in locals, keeping both engines' submitted claims identical

> **PARITY-EXCEPTION**: the Terraform Kubernetes provider's PVC resource cannot express `spec.dataSource`/`dataSourceRef` (clone a PVC / restore a VolumeSnapshot); the Pulumi module sends them natively. This module **fails the plan with a precondition** when `data_source` is set: deploy that claim with the Pulumi provisioner, or drop `spec.data_source`. Failing loudly beats silently provisioning an EMPTY volume where the user asked for restored data.

## Usage

```hcl
module "persistent_volume_claim" {
  source = "./iac/tf"

  metadata = {
    name = "app-data"
  }

  spec = {
    name      = "app-data"
    namespace = "backend"

    access_modes = ["ReadWriteOnce"]

    storage_request    = "10Gi"
    storage_class_name = "fast-ssd"
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `metadata` | Resource metadata (name, org, env) | object | yes |
| `spec` | PersistentVolumeClaim specification | object | yes |

## Outputs

| Name | Description |
|------|-------------|
| `pvc_name` | Name of the PersistentVolumeClaim object as created in the cluster |
| `namespace` | The namespace the claim was created in |
| `storage_request` | The requested storage size as a Kubernetes quantity |

> **Note**: A claim under a `wait_for_first_consumer` StorageClass stays Pending until a pod uses it — correct behavior, not an error. The outputs deliberately avoid bind-time status (bound volume name, phase) for the same reason.
