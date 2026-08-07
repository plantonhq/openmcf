# Kubernetes Priority Class - Terraform Module

## Overview

This Terraform module creates and manages a Kubernetes `scheduling.k8s.io/v1` PriorityClass. It supports the complete PriorityClass surface: the priority value, the global-default flag, the human description, and the preemption policy.

## Architecture

```
iac/tf/
├── provider.tf     # Terraform and Kubernetes provider requirements
├── variables.tf    # Input variables mirroring spec.proto
├── locals.tf       # Derived values: labels, annotations, resolved preemption policy
├── main.tf         # Creates kubernetes_priority_class_v1 resource
├── outputs.tf      # Exports priority_class_name, value
└── README.md       # This file
```

## How It Works

1. **Variable Input**: The `spec` variable mirrors the protobuf schema; the `preemption_policy` enum arrives as the proto enum value name (e.g. `"preempt_lower_priority"`, `"never"`)
2. **Label Merging**: Standard Planton labels are merged with user labels (identity keys cannot be overridden)
3. **Preemption Policy Resolution**: `locals.tf` maps the enum value name to the API string (`preempt_lower_priority` → `PreemptLowerPriority`) and applies the server default when the field is omitted. The resolved value is ALWAYS sent explicitly, mirroring the Pulumi module, so both engines submit byte-identical objects for the same manifest
4. **Resource Creation**: `main.tf` creates a single `kubernetes_priority_class_v1` resource; `global_default` is also sent explicitly (default `false`)
5. **Output Export**: Class name and value are exported

## Semantics Preserved by the Module

- **Value changes force replacement** — the priority value is immutable upstream; the provider forces replacement on change, matching the Pulumi module's delete-before-replace semantics (PriorityClass names are cluster-unique, so delete-then-create is the only safe order)
- **PriorityClasses are cluster-scoped** — no namespace is involved anywhere in the module
- **All optional fields resolve to the server defaults module-side** — a spec that omits them deploys exactly what the API server would have defaulted, with no engine drift

## Usage

```hcl
module "priority_class" {
  source = "./iac/tf"

  metadata = {
    name = "critical-services"
  }

  spec = {
    name        = "critical-services"
    value       = 1000000
    description = "Revenue-path services; preempts lower tiers under pressure."

    preemption_policy = "preempt_lower_priority"
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `metadata` | Resource metadata (name, org, env) | object | yes |
| `spec` | PriorityClass specification | object | yes |

## Outputs

| Name | Description |
|------|-------------|
| `priority_class_name` | Name of the PriorityClass object as created in the cluster |
| `value` | The priority integer pods of this class receive |

> **Note**: `global_default: true` affects every pod in the cluster that names no priority class. Only one class should carry the flag — when several do, Kubernetes uses the smallest such value.
