# Kubernetes PodDisruptionBudget

Deploys a Kubernetes PodDisruptionBudget — insurance against maintenance. The budget limits how many of a set of pods may be taken down VOLUNTARILY at once (node drains, cluster upgrades, autoscaler consolidation); the eviction API refuses to breach it. Manages availability policy declaratively through a Kubernetes Provider Connection with full audit trail and versioning.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kubernetes PodDisruptionBudget** -- a policy/v1 budget in the specified namespace carrying the pod selector, exactly one availability bound (min available OR max unavailable), and the unhealthy-pod eviction policy
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- The target namespace must already exist (the module does not create it).
- For a Planton Deployment or StatefulSet's OWN pods, prefer the workload's built-in disruption-budget block — it derives the selector automatically. This standalone kind is for pods a Planton workload does not manage.

## Deploy

### Console

Open the deployment store, find **PodDisruptionBudget on Kubernetes**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Protect Workload** preset — the one-at-a-time rolling-maintenance contract — in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPodDisruptionBudget
metadata:
  name: checkout-pdb
  org: acme-corp
  env: prod
spec:
  name: checkout-pdb
  namespace:
    value: backend-services
  selector:
    match_labels:
      app: checkout
  max_unavailable: "1"
```

```shell
planton apply -f pdb.yaml
```

This lets cluster maintenance move at most ONE `checkout` pod at a time — drains negotiate, they never surprise.

## Key Configuration

These are the most important decisions when configuring a Kubernetes PodDisruptionBudget. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Selection by the label contract** -- Every Planton workload stamps the `app` label (set to its name) on its pods, so `match_labels: {app: checkout}` protects exactly that workload. An EMPTY selector protects ALL pods in the namespace — legal, explicit, and rarely what you want.

**Exactly one bound** -- `max_unavailable` (preferred for workloads that scale — it tracks replica count) or `min_available` (an absolute floor that goes stale as replicas change). The API rejects both together, and a budget with neither protects nothing.

**The unsatisfiable-budget trap** -- `max_unavailable: "0"` or `min_available: "100%"` blocks ALL voluntary evictions including node drains; a single-replica workload behind such a budget wedges every cluster upgrade until someone relaxes it.

**The crash-loop wedge** -- A budget counts HEALTHY (Ready) pods. Under the conservative default, a crash-looping app can block drains forever; **Always Allow** lets not-ready pods be evicted regardless — the practical choice for workloads that can crash-loop.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | The namespace whose pods the budget protects; omitted means the cluster's `default` namespace |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `pod_disruption_budget_name` | The name of the created budget | `kubectl get pdb` allowed-disruptions checks |
| `namespace` | The namespace the budget protects | Verifying scope |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Rolling-maintenance contract** -- One workload, one pod movable at a time. Start from the **Protect Workload** preset.

**Percentage tier** -- At most 25% of a scaled fleet down at once. Start from the **Tier Percentage** preset.

**Crash-loop tolerant** -- Always Allow eviction of not-ready pods so a broken app never wedges upgrades. Start from the **Crashloop Tolerant** preset.

## Works With

- **Kubernetes Namespace** -- reference the namespace so infra charts create it and this budget in dependency order.
- **Kubernetes Deployment and StatefulSet** -- their built-in disruption-budget block covers their own pods; this kind covers operator-managed replicas and everything else, matched via their `app` label.
