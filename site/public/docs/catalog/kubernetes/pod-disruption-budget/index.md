---
title: "Pod Disruption Budget"
description: "Pod Disruption Budget deployment documentation"
icon: "package"
order: 100
componentName: "kubernetespoddisruptionbudget"
---

# Kubernetes Pod Disruption Budget

Deploys a Kubernetes PodDisruptionBudget — the availability floor for voluntary disruptions — to a target cluster through a single declarative manifest, covering the complete `policy/v1` surface: exact-match and set-based selectors, absolute and percentage bounds, and the unhealthy-pod eviction policy. The IaC module handles label merging, namespace resolution, and int-or-percent bound mapping automatically.

## What Gets Created

When you deploy a KubernetesPodDisruptionBudget resource, Planton provisions:

- **PodDisruptionBudget** — a `policy/v1` PodDisruptionBudget selecting pods via `selector` and declaring how many may be taken down voluntarily at once
- **Labels** — standard Planton tracking labels merged with any user-provided labels
- **Annotations** — user-provided annotations applied to the PodDisruptionBudget metadata

**Budgets govern only voluntary disruptions** — node drains, cluster upgrades, autoscaler consolidation, all of which go through the eviction API. Involuntary failures (node crashes, OOM kills) never consult the budget.

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, or a `KubernetesNamespace` resource referenced from `spec.namespace` so both deploy in one run
- **Replica headroom** — a budget can only refuse evictions; the selected workload must run enough replicas above the floor for drains to proceed

## Quick Start

Create a file `pod-disruption-budget.yaml` — this keeps at least one pod of a workload alive through any drain:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPodDisruptionBudget
metadata:
  name: checkout-pdb
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPodDisruptionBudget.checkout-pdb
spec:
  name: checkout-pdb
  namespace:
    value: backend
  selector:
    match_labels:
      app: checkout
  min_available: "1"
```

Deploy:

```shell
planton apply -f pod-disruption-budget.yaml
```

The selector targets the `checkout` workload's pods via the workload label contract — every Planton workload kind stamps `app: <workload-metadata-name>` on its pods. With `min_available: "1"`, the eviction API refuses any drain step that would leave zero `checkout` pods available.

> For a Planton Deployment's or StatefulSet's OWN pods, prefer the workload's built-in `availability.pod_disruption_budget` block — it derives the selector automatically. Use this standalone kind for operator-managed pods, non-Planton workloads, and selector-level coverage across workloads. Never point both at the same pods: Kubernetes fails evictions for pods covered by more than one budget.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the PodDisruptionBudget (`metadata.name` in the cluster). | 1–253 characters, matches `^[a-z0-9]([a-z0-9.-]{0,251}[a-z0-9])?$` |
| `spec.selector` | `LabelSelector` | The pods this budget protects. **An explicitly empty selector protects ALL pods in the namespace.** Target one Planton workload with `match_labels: {app: <workload-name>}`. | Required — a `policy/v1` budget with no selector protects nothing |
| one of `spec.min_available` / `spec.max_unavailable` | `string` | The availability bound (see below). | Exactly one must be set |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespace` | `StringValueOrRef` | `default` | Namespace the budget lives in — it governs only pods in its own namespace. Accepts a literal name (`{ value: my-namespace }`) or a reference to a `KubernetesNamespace` resource. |
| `spec.unhealthy_pod_eviction_policy` | `if_healthy_budget \| always_allow` | `if_healthy_budget` | How running-but-not-ready pods are treated. `always_allow` prevents crash-looping pods from wedging drains — and deploys via the Pulumi provisioner only (the Terraform provider cannot express the field and fails the plan by design). |
| `spec.labels` / `spec.annotations` | `map<string, string>` | `{}` | Merged with standard Planton labels / applied to the object. |

### The Availability Bound

| Field | Meaning | Notes |
|-------|---------|-------|
| `min_available` | At least this many selected pods must stay available — absolute (`"2"`) or percentage (`"50%"` of desired replicas) | Percentages round UP (stricter). `"100%"` blocks ALL voluntary evictions, including node drains |
| `max_unavailable` | At most this many selected pods may be down — absolute or percentage | Percentages round UP (more eviction room). `"0"` blocks all voluntary evictions. Prefer this form for workloads that scale — it tracks replica count |

## Examples

### Tier-Wide Budget with Set-Based Selection

One budget spanning every workload labelled `tier: web` or `tier: api`, allowing at most a quarter of them down at once:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPodDisruptionBudget
metadata:
  name: frontline-tier-pdb
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPodDisruptionBudget.frontline-tier-pdb
spec:
  name: frontline-tier-pdb
  namespace:
    value: backend
  selector:
    match_expressions:
      - key: tier
        operator: In
        values:
          - web
          - api
  max_unavailable: "25%"
```

### Crash-Loop-Tolerant Budget

A budget over an operator-managed database that keeps two replicas available but never lets a crash-looping replica wedge a node drain. `always_allow` requires the Pulumi provisioner:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPodDisruptionBudget
metadata:
  name: postgres-pdb
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPodDisruptionBudget.postgres-pdb
spec:
  name: postgres-pdb
  namespace:
    value: databases
  selector:
    match_labels:
      cluster-name: main-postgres
  min_available: "2"
  unhealthy_pod_eviction_policy: always_allow
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `podDisruptionBudgetName` | `string` | Name of the PodDisruptionBudget object as created in the cluster |
| `namespace` | `string` | Namespace the budget was created in |

## Related Components

- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) — provides the target namespace; reference it from `spec.namespace` to deploy both in one run
- [KubernetesDeployment](/docs/catalog/kubernetes/deployment) — for its OWN pods, use the built-in `availability.pod_disruption_budget` block instead; its pods carry the `app: <workload-name>` label this kind's selector targets for composition
- [KubernetesHorizontalPodAutoscaler](/docs/catalog/kubernetes/horizontal-pod-autoscaler) — the scaling counterpart; a budget bounds how fast the fleet can shrink during drains while the autoscaler owns how large it runs
