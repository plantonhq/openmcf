# Kubernetes Horizontal Pod Autoscaler - Terraform Module

## Overview

This Terraform module creates and manages a Kubernetes `autoscaling/v2` HorizontalPodAutoscaler. It supports the complete v2 surface: the scale target reference, replica floor and ceiling, all five metric source families (resource, container_resource, pods, object, external) with all three target value forms (utilization, value, average value), and per-direction scaling behavior with stabilization windows and velocity policies.

## Architecture

```
iac/tf/
├── provider.tf     # Terraform and Kubernetes provider requirements
├── variables.tf    # Input variables mirroring spec.proto
├── locals.tf       # Derived values: labels, annotations, namespace default, resolved target and floor, enum→API-string maps
├── main.tf         # Creates kubernetes_horizontal_pod_autoscaler_v2 resource
├── outputs.tf      # Exports horizontal_pod_autoscaler_name, namespace, scale_target, min_replicas, max_replicas
└── README.md       # This file
```

## How It Works

1. **Variable Input**: The `spec` variable mirrors the protobuf schema; the namespace and `scale_target.name` `StringValueOrRef` fields arrive flattened to plain strings, and enum fields arrive as the proto enum value names (e.g. `"resource"`, `"utilization"`, `"max_change"`, `"pods"`)
2. **Namespace Default**: `locals.tf` falls back to the `default` namespace when none is provided
3. **Label Merging**: Standard Planton labels are merged with user labels (identity keys cannot be overridden)
4. **Default Resolution**: `locals.tf` applies the spec defaults — `apps/v1` `Deployment` scale target, replica floor 1 — identically to the Pulumi module, and both are ALWAYS sent explicitly so both engines submit byte-identical objects for the same manifest
5. **Enum Mapping**: `locals.tf` maps proto enum value names to the API strings: metric types to `Resource`/`ContainerResource`/`Pods`/`Object`/`External`, target types to `Utilization`/`Value`/`AverageValue` (the spec's `raw_value` is the API's `Value`), select policies to `Max`/`Min`/`Disabled`, policy types to `Pods`/`Percent`
6. **Resource Creation**: `main.tf` creates a single `kubernetes_horizontal_pod_autoscaler_v2` resource; dynamic blocks render at most one source per metric — the spec's validations guarantee exactly the source matching each metric's type is present
7. **Output Export**: Autoscaler name, namespace, scale target (`Kind/name`), and replica bounds are exported

## Semantics Preserved by the Module

- **`min_replicas` is always sent explicitly** with the Kubernetes default (1) applied in locals
- **When the spec lists no metrics, the metric blocks are OMITTED** — the API server then applies its own default (80% average CPU utilization); an empty list would instead disable metric-driven scaling
- **`select_policy` is always sent in behavior blocks** with the API default (`Max`) applied, mirroring the Pulumi module
- **Metric identifier `match_labels` render as the metric selector** — scoping which series the metrics adapter reads

## Usage

```hcl
module "horizontal_pod_autoscaler" {
  source = "./iac/tf"

  metadata = {
    name = "checkout-hpa"
  }

  spec = {
    name      = "checkout-hpa"
    namespace = "backend"

    scale_target = {
      name = "checkout"
    }

    min_replicas = 2
    max_replicas = 10

    metrics = [{
      type = "resource"
      resource = {
        name = "cpu"
        target = {
          type                = "utilization"
          average_utilization = 60
        }
      }
    }]
  }
}
```

## Inputs

| Name | Description | Type | Required |
|------|-------------|------|----------|
| `metadata` | Resource metadata (name, org, env) | object | yes |
| `spec` | HorizontalPodAutoscaler specification | object | yes |

## Outputs

| Name | Description |
|------|-------------|
| `horizontal_pod_autoscaler_name` | Name of the HorizontalPodAutoscaler object as created in the cluster |
| `namespace` | Namespace the autoscaler was created in |
| `scale_target` | The scale target as `Kind/name` |
| `min_replicas` | The configured replica floor |
| `max_replicas` | The configured replica ceiling |

> **Note**: The HPA reads metrics; it does not produce them. Resource metrics require metrics-server; `pods`/`object`/`external` metrics require a custom/external-metrics adapter (e.g. prometheus-adapter, KEDA). The scale target's existence is not required at creation — the controller reports it unresolved until the workload appears.
