---
title: "Horizontal Pod Autoscaler"
description: "Horizontal Pod Autoscaler deployment documentation"
icon: "package"
order: 100
componentName: "kuberneteshorizontalpodautoscaler"
---

# Kubernetes Horizontal Pod Autoscaler

Deploys a Kubernetes HorizontalPodAutoscaler — automatic replica scaling driven by observed metrics — to a target cluster through a single declarative manifest, covering the complete `autoscaling/v2` surface: resource and per-container resource metrics, custom per-pod metrics, object and external metrics, and per-direction scaling behavior with stabilization windows and velocity policies. The IaC module handles label merging, namespace resolution, and default resolution automatically.

## What Gets Created

When you deploy a KubernetesHorizontalPodAutoscaler resource, Planton provisions:

- **HorizontalPodAutoscaler** — an `autoscaling/v2` HPA owning the scale target's replica count between `min_replicas` and `max_replicas`
- **Labels** — standard Planton tracking labels merged with any user-provided labels
- **Annotations** — user-provided annotations applied to the HorizontalPodAutoscaler metadata

**Each configured metric proposes a replica count and the HIGHEST wins** — metrics OR together toward scale-out. Once the HPA governs a workload, the workload's own `replicas` field becomes advisory: the autoscaler owns the count.

## Prerequisites

- **metrics-server** in the cluster for CPU/memory metrics; **a custom/external-metrics adapter** (prometheus-adapter, KEDA, a cloud adapter) for `pods`, `object`, and `external` metrics
- **Resource requests declared on the target's pods** for `utilization` targets — utilization is measured against requests
- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, or a `KubernetesNamespace` resource referenced from `spec.namespace` — the autoscaler must live in the scale target's own namespace

## Quick Start

Create a file `hpa.yaml` — this holds a workload's average CPU at 60% of requests, between 2 and 10 replicas:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesHorizontalPodAutoscaler
metadata:
  name: checkout-hpa
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesHorizontalPodAutoscaler.checkout-hpa
spec:
  name: checkout-hpa
  namespace:
    value: backend
  scale_target:
    name:
      value: checkout
  min_replicas: 2
  max_replicas: 10
  metrics:
    - type: resource
      resource:
        name: cpu
        target:
          type: utilization
          average_utilization: 60
```

Deploy:

```shell
planton apply -f hpa.yaml
```

The scale target defaults to an `apps/v1` `Deployment`; only the name is required. The target does not need to exist when the HPA is created — the controller reports it unresolved until the workload appears, the correct steady state for infra deployed ahead of the app.

> For simple CPU/memory autoscaling of a Planton Deployment's OWN replicas, prefer the workload's built-in `availability.horizontal_pod_autoscaling` block. Use this standalone kind for operator-managed and non-Planton targets, and for the advanced v2 surface (custom/object/external metrics, per-container metrics, behavior tuning). Never point both at the same target — two controllers fighting over one replica count flaps the fleet.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the HorizontalPodAutoscaler (`metadata.name` in the cluster). | 1–253 characters, matches `^[a-z0-9]([a-z0-9.-]{0,251}[a-z0-9])?$` |
| `spec.scale_target.name` | `StringValueOrRef` | The target workload's name, in the autoscaler's own namespace. Accepts a literal or a reference to a `KubernetesDeployment`'s exported name. | Required |
| `spec.max_replicas` | `int32` | The replica ceiling — the honest capacity conversation. | ≥ 1, and ≥ `min_replicas` |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespace` | `StringValueOrRef` | `default` | Namespace the autoscaler lives in — MUST be the scale target's own namespace (an HPA cannot scale across namespaces). |
| `spec.scale_target.api_version` | `string` | `apps/v1` | API version of the target. |
| `spec.scale_target.kind` | `string` | `Deployment` | Kind of the target — anything exposing the `scale` subresource. DaemonSets are rejected (one pod per node by definition). |
| `spec.min_replicas` | `int32` | `1` | The replica floor; must be ≥ 1 (scale-to-zero is feature-gated upstream and not modeled). |
| `spec.metrics` | `list(metric)` | `[]` | The metrics driving the decision; each proposes a count and the highest wins. **Empty applies the Kubernetes default: 80% average CPU utilization** (requires CPU requests on the pods). |
| `spec.behavior` | `behavior` | Kubernetes defaults | Per-direction tuning. Defaults: scale up fast (double or +4 pods per 15s, no stabilization); scale down to the recommendation's 300-second high-water mark. |
| `spec.labels` / `spec.annotations` | `map<string, string>` | `{}` | Merged with standard Planton labels / applied to the object. |

### Metrics

Each entry sets `type` and exactly the matching source field:

| Type | Source | Needs | Typical Target |
|------|--------|-------|----------------|
| `resource` | cpu/memory averaged across pods | metrics-server | `utilization` (percent of requests) |
| `container_resource` | cpu/memory of ONE named container | metrics-server | `utilization` — isolates the app from sidecar skew |
| `pods` | a custom per-pod metric, averaged | custom-metrics adapter | `average_value` |
| `object` | a metric describing one other object (e.g. an Ingress) | custom-metrics adapter | `raw_value` or `average_value` |
| `external` | a metric from outside the cluster (queue depth, LB QPS) | external-metrics adapter (e.g. KEDA, prometheus-adapter) | `average_value` (e.g. "30 messages per pod") |

Target types: `utilization` (percent of requests), `average_value` (quantity per pod), `raw_value` (raw quantity — the upstream API's `Value` type).

### Behavior

| Field | Description |
|-------|-------------|
| `scale_up` / `scale_down` | Per-direction tuning blocks |
| `stabilization_window_seconds` | 0–3600; the flap damper — the safest recommendation in the window wins |
| `select_policy` | `max_change` (default), `min_change`, or `disabled` (turns the direction off entirely) |
| `policies[]` | Velocity caps: `type` (`pods`/`percent`), `value`, `period_seconds` (1–1800) |

## Examples

### Queue-Driven Worker

Scale a worker fleet on external queue depth — roughly one pod per 30 ready messages. Requires an external-metrics adapter (e.g. KEDA or prometheus-adapter) serving the metric:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesHorizontalPodAutoscaler
metadata:
  name: worker-hpa
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesHorizontalPodAutoscaler.worker-hpa
spec:
  name: worker-hpa
  namespace:
    value: backend
  scale_target:
    name:
      value: order-worker
  min_replicas: 1
  max_replicas: 30
  metrics:
    - type: external
      external:
        metric:
          name: queue_messages_ready
          match_labels:
            queue: orders
        target:
          type: average_value
          average_value: "30"
```

### Sidecar-Isolated CPU Scaling with Conservative Scale-Down

Scale on the app container's CPU only (a hot proxy sidecar cannot mask an idle app), and bleed capacity down at most 10% per minute after a 10-minute stabilization window:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesHorizontalPodAutoscaler
metadata:
  name: api-hpa
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesHorizontalPodAutoscaler.api-hpa
spec:
  name: api-hpa
  namespace:
    value: backend
  scale_target:
    name:
      value: api-server
  min_replicas: 3
  max_replicas: 20
  metrics:
    - type: container_resource
      container_resource:
        name: cpu
        container: app
        target:
          type: utilization
          average_utilization: 60
  behavior:
    scale_down:
      stabilization_window_seconds: 600
      policies:
        - type: percent
          value: 10
          period_seconds: 60
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `horizontalPodAutoscalerName` | `string` | Name of the HorizontalPodAutoscaler object as created in the cluster |
| `namespace` | `string` | Namespace the autoscaler was created in |
| `scaleTarget` | `string` | The scale target as `"Kind/name"` (e.g. `"Deployment/checkout"`) |
| `minReplicas` | `int32` | The configured replica floor |
| `maxReplicas` | `int32` | The configured replica ceiling |

## Related Components

- [KubernetesDeployment](/docs/catalog/kubernetes/deployment) — the usual scale target; reference its exported name from `scale_target.name`. For simple CPU/memory scaling of its own replicas, use its built-in `availability.horizontal_pod_autoscaling` block instead
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) — provides the target namespace; reference it from `spec.namespace` to deploy both in one run
- [KubernetesPodDisruptionBudget](/docs/catalog/kubernetes/pod-disruption-budget) — the availability counterpart; the budget bounds how fast the fleet can shrink during drains while the autoscaler owns how large it runs
