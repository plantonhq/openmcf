# Planton Operator

Installs the Planton operator — the lifecycle manager that reconciles `PlantonPlatform` declarations into running self-hosted Planton platforms (control plane, console, identity, databases, secrets manager, in-cluster runner) — from the official `planton-operator` Helm chart (OCI, ghcr.io/plantonhq/charts). This component installs the MANAGER; the platforms themselves are declared with KubernetesPlantonPlatform resources, one per platform, all served by this one operator. One installation per cluster: the operator enforces that itself at startup, and the release name is fixed to `planton-operator`.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm Release** (`planton-operator`) -- the operator's manager Deployment, its ServiceAccount, the cluster-wide reconciliation ClusterRole/ClusterRoleBinding, and the namespaced leader-election Role/RoleBinding
- **CRD** (`plantonplatforms.planton.ai`) -- module-owned: applied from a copy staged at the pinned chart version, KEPT on uninstall so removing the operator never cascade-deletes platform declarations, and upgraded deliberately with every `chart_version` bump
- **Namespace** (optional) -- created with standard governance labels when `create_namespace` is true (`planton-operator` is the convention)

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with credentials for the target cluster (installing an operator is a cluster-admin act — the chart grants the manager its cluster-wide reconciliation RBAC). Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### Cluster Prerequisites

- None beyond the connection — the operator is self-contained. The platforms it deploys bring their own dependencies (declared per KubernetesPlantonPlatform).

## Deploy

### Console

Open the deployment store, find **Planton Operator**, and click **Deploy**. The creation wizard walks you through placement, the chart pin, operator runtime, image sourcing, and scheduling. Start from the **Standard** preset in the Presets tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPlantonOperator
metadata:
  name: planton-operator
  org: acme-corp
  env: platform
spec:
  namespace:
    value: planton-operator
  create_namespace: true
```

```bash
planton apply -f planton-operator.yaml
```

Then declare each platform with a KubernetesPlantonPlatform resource — the operator reconciles it to Ready and the console's first visitor becomes the admin.

## Destroy

Destroying this resource removes the manager and its RBAC; the `plantonplatforms.planton.ai` CRD is deliberately KEPT, so platform declarations — and the running platforms behind them — survive, unmanaged, until an operator is reinstalled and adopts them. Platform deletion completes even without the operator (teardown is Kubernetes garbage collection of each declaration's owner-referenced objects).
