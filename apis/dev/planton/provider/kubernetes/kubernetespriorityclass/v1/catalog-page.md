# Kubernetes Priority Class

Deploys a Kubernetes PriorityClass — one rung of the cluster's workload importance ladder — through a single declarative manifest, covering the complete `scheduling.k8s.io/v1` surface: the priority value, the global-default flag, the human description, and the preemption policy. The IaC module handles label merging, explicit preemption-policy defaults, and safe delete-before-replace on value changes automatically.

## What Gets Created

When you deploy a KubernetesPriorityClass resource, Planton provisions:

- **PriorityClass** — a cluster-scoped `scheduling.k8s.io/v1` PriorityClass; pods reference it by name in `priority_class_name` to receive its priority value
- **Labels** — standard Planton tracking labels merged with any user-provided labels
- **Annotations** — user-provided annotations applied to the PriorityClass metadata

**Higher priority schedules first — and, by default, evicts.** Among Pending pods, higher priority goes first; when a higher-priority pod cannot schedule, the scheduler preempts (evicts) lower-priority pods to make room, unless the class sets `preemption_policy: never`.

## Prerequisites

- **A Kubernetes cluster** — PriorityClass is GA (`scheduling.k8s.io/v1`) since Kubernetes 1.14, available everywhere
- **Kubernetes credentials** configured via environment variables or Planton provider config
- **Cluster-scoped permissions** — PriorityClasses are cluster-scoped; the credentials must be allowed to create them

## Quick Start

Create a file `priority-class.yaml` — the critical tier of the ladder:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPriorityClass
metadata:
  name: critical-services
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPriorityClass.critical-services
spec:
  name: critical-services
  value: 1000000
  description: Revenue-path services; preempts lower tiers under capacity pressure.
```

Deploy:

```shell
planton apply -f priority-class.yaml
```

Any pod that sets `priority_class_name: critical-services` now schedules ahead of lower-priority pods and, under pressure, evicts them to make room.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the PriorityClass (`metadata.name` in the cluster) — the value pods reference in `priority_class_name`. | 1–253 characters, matches `^[a-z0-9]([a-z0-9.-]{0,251}[a-z0-9])?$`; must not start with the reserved `system-` prefix |
| `spec.value` | `int32` | The priority integer pods of this class receive — higher schedules (and preempts) ahead of lower. Negative values are valid and useful for always-preemptable batch tiers. **IMMUTABLE**: changing it replaces the class (delete-before-replace). | At most 1,000,000,000; the range above is reserved for Kubernetes system classes |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.global_default` | `bool` | `false` | Makes this class the cluster-wide default for pods that name NO priority class. Only one class should be the global default — when several claim it, Kubernetes uses the SMALLEST such value. Never re-prioritizes existing pods. |
| `spec.description` | `string` | `""` | Human guidance on when to use this class (surfaced by `kubectl describe priorityclass`). Write it for the next engineer choosing a class. |
| `spec.preemption_policy` | `preempt_lower_priority \| never` | `preempt_lower_priority` | Whether pending pods of this class evict lower-priority running pods. `never` = jump the queue but never evict — the right policy for high-priority batch work. |
| `spec.labels` / `spec.annotations` | `map<string, string>` | `{}` | Merged with standard Planton labels / applied to the object. |

## Examples

### The Standard Default Tier

The class every unmarked pod lands in — exactly one class per cluster should carry `global_default: true`:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPriorityClass
metadata:
  name: standard-default
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPriorityClass.standard-default
spec:
  name: standard-default
  value: 1000
  global_default: true
  description: Cluster default for workloads that do not choose a class.
```

### The Preemptable Batch Tier

Negative value (yields to everything) and preemption disabled (never evicts running pods):

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPriorityClass
metadata:
  name: preemptable-batch
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPriorityClass.preemptable-batch
spec:
  name: preemptable-batch
  value: -100
  description: Batch workloads; always preemptable, never preempts.
  preemption_policy: never
```

### Referencing a Class from a Workload

Pods opt in by name — the shared workload pod spec's `priority_class_name` field:

```yaml
# In a KubernetesDeployment (or any Planton workload kind) pod spec:
priority_class_name: critical-services
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `priorityClassName` | `string` | Name of the PriorityClass object as created in the cluster — the composition handle workload pod specs reference |
| `value` | `int32` | The priority integer pods of this class receive |

## Related Components

- [KubernetesDeployment](/docs/catalog/kubernetes/kubernetesdeployment) — workload whose pod spec references the class via `priority_class_name`
- [KubernetesResourceQuota](/docs/catalog/kubernetes/kubernetesresourcequota) — budgets pods of specific priority classes with its `priority_class` scope selector
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) — namespace-level governance that composes with cluster-level priority tiers
