# Kubernetes Resource Quota - Terraform Module

## Overview

This Terraform module creates and manages a namespace-governance pair: a `core/v1` **ResourceQuota** carrying aggregate caps, scopes, and a scope selector, plus an optional companion `core/v1` **LimitRange** created when the spec sets `limit_defaults`. It supports the complete ResourceQuotaSpec and LimitRangeSpec surfaces: compute, storage, and object-count caps; all seven quota scopes; scope selectors with `In`/`NotIn`/`Exists`/`DoesNotExist`; and container, pod, and persistent-volume-claim limit items with defaults, bounds, and burst ratios.

## Architecture

```
iac/tf/
├── provider.tf     # Terraform and Kubernetes provider requirements
├── variables.tf    # Input variables mirroring spec.proto
├── locals.tf       # Derived values: labels, annotations, namespace default, enum→API-string maps, LimitRange condition
├── main.tf         # Creates kubernetes_resource_quota_v1 and count-gated kubernetes_limit_range_v1
├── outputs.tf      # Exports resource_quota_name, namespace, limit_range_name
└── README.md       # This file
```

## How It Works

1. **Variable Input**: The `spec` variable mirrors the protobuf schema; the namespace `StringValueOrRef` arrives flattened to a plain string, and enum fields arrive as the proto enum value names (e.g. `"best_effort"`, `"container"`)
2. **Namespace Default**: `locals.tf` falls back to the `default` namespace when none is provided
3. **Label Merging**: Standard Planton labels are merged with user labels (identity keys cannot be overridden); both created objects receive the same labels and annotations
4. **Enum Mapping**: `locals.tf` maps the proto enum value names to the Kubernetes API strings (`best_effort` → `BestEffort`, `persistent_volume_claim` → `PersistentVolumeClaim`), identical to the Pulumi module's mapping
5. **Resource Creation**: `main.tf` creates the `kubernetes_resource_quota_v1` resource and, gated on `limit_defaults` being non-empty, the `kubernetes_limit_range_v1` resource
6. **Output Export**: Quota name, namespace, and the LimitRange name (empty when none exists) are exported

## The Two-Object Creation

This module manages **two Kubernetes resources as one governance pair**:

- The **ResourceQuota** is always created
- The **LimitRange** is created exactly when the spec carries `limit_defaults` (a `count`-gated resource), and it **shares the quota's name and namespace** — one governance pair, one identity
- The pairing is what keeps a compute quota livable: once a quota caps `requests.*`/`limits.*`, the API rejects pods that omit them; the LimitRange's container defaults let naive workloads inherit values instead of being rejected

> For simple T-shirt-size governance on namespaces Planton creates, the `KubernetesNamespace` kind's resource profiles (which manage a quota and limit range internally) are the simpler alternative. This module is the full-fidelity instrument.

## Usage

```hcl
module "resource_quota" {
  source = "./iac/tf"

  metadata = {
    name = "team-alpha-quota"
  }

  spec = {
    name      = "team-alpha-quota"
    namespace = "team-alpha"

    hard = {
      "requests.cpu"    = "10"
      "requests.memory" = "20Gi"
      "limits.cpu"      = "20"
      "limits.memory"   = "40Gi"
    }

    limit_defaults = [{
      type = "container"
      default_request = {
        cpu    = "100m"
        memory = "128Mi"
      }
      default_limit = {
        cpu    = "500m"
        memory = "512Mi"
      }
    }]
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `metadata` | Resource metadata (name, org, env) | object | yes |
| `spec` | ResourceQuota + LimitRange specification | object | yes |

## Outputs

| Name | Description |
|------|-------------|
| `resource_quota_name` | Name of the ResourceQuota object as created in the cluster |
| `namespace` | Namespace the quota governs |
| `limit_range_name` | Name of the companion LimitRange; empty when no `limit_defaults` were configured |

> **Note**: Quota is enforced at admission time — objects created before the quota existed count against usage but are never evicted, and a compute-capping quota without `limit_defaults` rejects pods that omit requests/limits.
