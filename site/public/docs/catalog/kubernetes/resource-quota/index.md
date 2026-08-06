---
title: "Resource Quota"
description: "Resource Quota deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesresourcequota"
---

# Kubernetes Resource Quota

Governs resource consumption in one namespace through a single declarative manifest: aggregate caps on what the namespace may use in total (a `core/v1` ResourceQuota) plus optional per-object defaults and bounds (a companion `core/v1` LimitRange), covering the complete upstream surface — compute, storage, and object-count caps, scope filters, and container/pod/claim limit items. The IaC module handles label merging, namespace resolution, and the two-object lifecycle automatically.

## What Gets Created

When you deploy a KubernetesResourceQuota resource, Planton provisions:

- **ResourceQuota** — a `core/v1` ResourceQuota carrying the `hard` caps, scopes, and scope selector; the API server rejects any object whose creation would exceed a cap
- **LimitRange** (optional) — a companion `core/v1` LimitRange created only when `spec.limit_defaults` is set, sharing the quota's name; applies per-container defaults and per-object bounds
- **Labels** — standard Planton tracking labels merged with any user-provided labels, on both objects
- **Annotations** — user-provided annotations applied to both objects' metadata

**Capping compute makes the API reject pods that omit requests/limits.** Once a quota caps `requests.cpu` or `limits.memory`, a pod that doesn't declare them fails admission. Pairing `hard` with `limit_defaults` is the safe pattern: naive workloads inherit the defaults instead of being rejected.

## Prerequisites

- **A Kubernetes cluster** — the ResourceQuota and LimitRange admission plugins are enabled by default on all mainstream distributions
- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, or a `KubernetesNamespace` resource referenced from `spec.namespace` so both deploy in one run

> For simple T-shirt-size governance on namespaces Planton creates, prefer `KubernetesNamespace`'s resource profiles — they manage a quota and limit range internally. This kind is the full-fidelity instrument.

## Quick Start

Create a file `resource-quota.yaml` — compute caps paired with container defaults, the safe pattern:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesResourceQuota
metadata:
  name: team-alpha-quota
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesResourceQuota.team-alpha-quota
spec:
  namespace:
    value: team-alpha
  name: team-alpha-quota
  hard:
    requests.cpu: "10"
    requests.memory: 20Gi
    limits.cpu: "20"
    limits.memory: 40Gi
  limit_defaults:
    - type: container
      default_request:
        cpu: 100m
        memory: 128Mi
      default_limit:
        cpu: 500m
        memory: 512Mi
```

Deploy:

```shell
planton apply -f resource-quota.yaml
```

The namespace may now consume at most 10 CPUs / 20Gi of requests in aggregate, and containers that omit requests/limits inherit the defaults instead of being rejected by the quota's admission check.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the ResourceQuota (`metadata.name` in the cluster); the companion LimitRange shares this name. | 1–253 characters, matches `^[a-z0-9]([a-z0-9.-]{0,251}[a-z0-9])?$` |
| `spec.hard` | `map<string, string>` | The aggregate caps, resource name → quantity. Compute (`requests.cpu`, `limits.memory`, ...), storage (`requests.storage`, `persistentvolumeclaims`, per-class variants), object counts (`pods`, `services`, `services.loadbalancers`, `count/<resource>.<group>`, ...). | At least one entry; values non-empty |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespace` | `StringValueOrRef` | `default` | Namespace to govern. Accepts a literal name (`{ value: my-namespace }`) or a reference to a `KubernetesNamespace` resource. |
| `spec.scopes` | `list(scope)` | `[]` | Coarse filters on which objects the quota tracks — an object must match ALL listed scopes. Values: `terminating`, `not_terminating`, `best_effort`, `not_best_effort`, `priority_class`, `cross_namespace_pod_affinity`, `volume_attributes_class`. A `best_effort` quota may cap only `pods`; conflicting pairs are rejected. |
| `spec.scope_selector` | `list(requirement)` | `[]` | Fine-grained scope filters — most usefully `priority_class` with `In`/`NotIn` values. ANDs with `scopes`. Pod-behavior scopes accept only `Exists`. |
| `spec.limit_defaults` | `list(item)` | `[]` | Per-object defaults and bounds, managed as the companion LimitRange. Omit entirely to manage only the quota — but a compute quota without defaults rejects pods that omit requests/limits. |
| `spec.labels` / `spec.annotations` | `map<string, string>` | `{}` | Merged with standard Planton labels / applied to both objects. |

### Limit Defaults Items

Each `limit_defaults` item governs one object type and must set at least one constraint:

| Field | Applies To | Description |
|-------|-----------|-------------|
| `type` | — | `container`, `pod`, or `persistent_volume_claim` |
| `max` / `min` | all types | Bounds any single object of this type must respect, per resource (claims use `{storage: ...}`) |
| `default_limit` | `container` only | LIMIT applied to containers that omit their own; when set alone, Kubernetes copies it into the request (Guaranteed QoS) |
| `default_request` | `container` only | REQUEST applied to containers that omit their own — the field that keeps a requests-capping quota from rejecting naive pods |
| `max_limit_request_ratio` | all types | Maximum limit-to-request burst ratio per resource (e.g. `{cpu: "4"}` allows at most 4x) |

## Examples

### Object-Count Caps Only

The safest quota to introduce on a live namespace — it constrains nothing that pods must declare, so naive pod creation keeps working:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesResourceQuota
metadata:
  name: team-alpha-object-caps
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesResourceQuota.team-alpha-object-caps
spec:
  namespace:
    value: team-alpha
  name: team-alpha-object-caps
  hard:
    pods: "100"
    services: "20"
    services.loadbalancers: "2"
    persistentvolumeclaims: "10"
```

### Cap Only BestEffort Pods

The quota tracks only pods with no requests or limits at all — containing unbounded naive pods without touching well-behaved workloads. A `best_effort`-scoped quota may cap only `pods`:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesResourceQuota
metadata:
  name: besteffort-guard
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesResourceQuota.besteffort-guard
spec:
  namespace:
    value: team-alpha
  name: besteffort-guard
  hard:
    pods: "10"
  scopes:
    - best_effort
```

### Separate Budget for High-Priority Pods

A scope selector budgets pods of a specific priority class independently:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesResourceQuota
metadata:
  name: critical-tier-budget
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesResourceQuota.critical-tier-budget
spec:
  namespace:
    value: team-alpha
  name: critical-tier-budget
  hard:
    requests.cpu: "4"
    requests.memory: 8Gi
    pods: "20"
  scope_selector:
    - scope_name: priority_class
      operator: In
      values:
        - critical-services
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `resourceQuotaName` | `string` | Name of the ResourceQuota object as created in the cluster |
| `namespace` | `string` | Namespace the quota governs |
| `limitRangeName` | `string` | Name of the companion LimitRange object; empty when the spec set no `limit_defaults` (no LimitRange is created) |

## Related Components

- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) — provides the target namespace; its resource profiles are the simpler alternative for T-shirt-size governance on namespaces Planton creates
- [KubernetesPriorityClass](/docs/catalog/kubernetes/priority-class) — defines the priority classes a `priority_class` scope selector budgets against
- [KubernetesNetworkPolicy](/docs/catalog/kubernetes/network-policy) — the other per-namespace governance instrument: network segmentation alongside resource caps
