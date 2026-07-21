# Kubernetes Network Policy - Terraform Module

## Overview

This Terraform module creates and manages a Kubernetes `networking/v1` NetworkPolicy. It supports the complete NetworkPolicy surface: pod selector, policy types, ingress and egress rules with all three peer forms (pod selector, namespace selector, IP block with exceptions), TCP/UDP/SCTP ports, named ports, and port ranges.

## Architecture

```
iac/tf/
├── provider.tf     # Terraform and Kubernetes provider requirements
├── variables.tf    # Input variables mirroring spec.proto
├── locals.tf       # Derived values: labels, annotations, namespace default, resolved policy types
├── main.tf         # Creates kubernetes_network_policy_v1 resource
├── outputs.tf      # Exports network_policy_name, namespace, policy_types
└── README.md       # This file
```

## How It Works

1. **Variable Input**: The `spec` variable mirrors the protobuf schema; the namespace `StringValueOrRef` arrives flattened to a plain string, and enum fields arrive as the proto enum value names (e.g. `"ingress"`, `"egress"`, `"TCP"`)
2. **Namespace Default**: `locals.tf` falls back to the `default` namespace when none is provided
3. **Label Merging**: Standard Planton labels are merged with user labels (identity keys cannot be overridden)
4. **Policy Type Resolution**: `locals.tf` maps the spec's directions to the API strings (`ingress` → `Ingress`) and, when the spec omits `policy_types`, applies the Kubernetes inference rule (ingress always; egress only when egress rules exist). The resolved set is ALWAYS sent explicitly, mirroring the Pulumi module, so both engines submit byte-identical direction sets for the same manifest
5. **Resource Creation**: `main.tf` creates a single `kubernetes_network_policy_v1` resource
6. **Output Export**: Policy name, namespace, and the governed directions (as `"Ingress"`, `"Egress"`, or `"Ingress,Egress"`) are exported

## Semantics Preserved by the Module

- **Each spec peer maps to exactly one API peer** — a peer with both `pod_selector` and `namespace_selector` stays one AND'd peer (pods matching the selector in the selected namespaces) and is never split into two OR'd peers
- **An absent pod selector renders as the empty selector block** — "all pods in the namespace", the wire form default-deny policies depend on
- **An empty port string matches all ports for the protocol**; a numeric string matches a port number, anything else a named container port; `end_port > 0` expresses a contiguous range

## Usage

```hcl
module "network_policy" {
  source = "./iac/tf"

  metadata = {
    name = "backend-api-ingress"
  }

  spec = {
    name      = "backend-api-ingress"
    namespace = "backend"

    pod_selector = {
      match_labels = { app = "backend-api" }
    }

    policy_types = ["ingress"]

    ingress_rules = [{
      from = [{
        pod_selector = {
          match_labels = { app = "frontend" }
        }
      }]
      ports = [{
        protocol = "TCP"
        port     = "8080"
      }]
    }]
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `metadata` | Resource metadata (name, org, env) | object | yes |
| `spec` | NetworkPolicy specification | object | yes |

## Outputs

| Name | Description |
|------|-------------|
| `network_policy_name` | Name of the NetworkPolicy object as created in the cluster |
| `namespace` | Namespace the NetworkPolicy was created in |
| `policy_types` | Governed directions as deployed, including inferred types |

> **Note**: The created object is only enforced by a NetworkPolicy-implementing CNI (Calico, Cilium, cloud CNIs with enforcement enabled). On clusters whose CNI ignores NetworkPolicies, the apply succeeds but traffic is unaffected.
