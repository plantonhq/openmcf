# Kubernetes Pod Disruption Budget - Terraform Module

## Overview

This Terraform module creates and manages a Kubernetes `policy/v1` PodDisruptionBudget. It supports exact-match and set-based label selectors and absolute and percentage availability bounds (`min_available` / `max_unavailable`). One field is deliberately out of reach — see the parity exception below.

## Architecture

```
iac/tf/
├── provider.tf     # Terraform and Kubernetes provider requirements
├── variables.tf    # Input variables mirroring spec.proto
├── locals.tf       # Derived values: labels, annotations, namespace default, resolved eviction policy
├── main.tf         # Creates kubernetes_pod_disruption_budget_v1 resource
├── outputs.tf      # Exports pod_disruption_budget_name, namespace
└── README.md       # This file
```

## How It Works

1. **Variable Input**: The `spec` variable mirrors the protobuf schema; the namespace `StringValueOrRef` arrives flattened to a plain string, and the `unhealthy_pod_eviction_policy` enum arrives as the proto enum value name (e.g. `"if_healthy_budget"`, `"always_allow"`)
2. **Namespace Default**: `locals.tf` falls back to the `default` namespace when none is provided
3. **Label Merging**: Standard Planton labels are merged with user labels (identity keys cannot be overridden)
4. **Policy Resolution**: `locals.tf` maps the eviction-policy enum to the API string with the server default (`IfHealthyBudget`) applied — used by the parity-exception precondition below
5. **Resource Creation**: `main.tf` creates a single `kubernetes_pod_disruption_budget_v1` resource with the selector always rendered and the chosen bound passed through as IntOrString
6. **Output Export**: Budget name and namespace are exported

## PARITY-EXCEPTION: unhealthy_pod_eviction_policy

The Terraform kubernetes provider's PDB resource **cannot express `spec.unhealthyPodEvictionPolicy` at all** — the field simply does not exist in the resource schema. The Pulumi module always sends it (with the server default `IfHealthyBudget` applied).

This module therefore carries a plan-time `precondition` that **fails any plan where the spec requests `always_allow`**, with an error message directing the user to deploy via the Pulumi provisioner or drop the field. Failing loudly is the design: the server default is exactly what a non-default value would override, so silently deploying `IfHealthyBudget` where `always_allow` was requested would deploy the opposite of the user's intent. For every spec this module accepts, both engines produce identical deployed objects.

## Semantics Preserved by the Module

- **The selector block is always rendered** (it is required in the spec) — an empty selector block is the "all pods in the namespace" wire form, deliberately explicit, because a `policy/v1` budget with a NULL selector matches no pods at all
- **Bounds pass through as IntOrString** — a numeric string (`"2"`) is an absolute count, a `%`-suffixed one (`"25%"`) a percentage; the spec enforces exactly one bound is set

## Usage

```hcl
module "pod_disruption_budget" {
  source = "./iac/tf"

  metadata = {
    name = "checkout-pdb"
  }

  spec = {
    name      = "checkout-pdb"
    namespace = "backend"

    selector = {
      match_labels = { app = "checkout" }
    }

    min_available = "1"
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `metadata` | Resource metadata (name, org, env) | object | yes |
| `spec` | PodDisruptionBudget specification | object | yes |

## Outputs

| Name | Description |
|------|-------------|
| `pod_disruption_budget_name` | Name of the PodDisruptionBudget object as created in the cluster |
| `namespace` | Namespace the budget was created in |

> **Note**: The budget governs only voluntary disruptions (the eviction API). A budget can only refuse evictions — the selected workload must run enough replicas above the declared floor for node drains to proceed.
