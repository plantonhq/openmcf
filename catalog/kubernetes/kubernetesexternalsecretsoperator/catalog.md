# External Secrets Operator

Installs the External Secrets Operator (ESO) — the controller that syncs secrets FROM external stores (AWS Secrets Manager, GCP Secret Manager, Azure Key Vault, Vault/OpenBao, and more) INTO Kubernetes Secrets — from the official Helm chart. ESO runs three components: the controller (reconciles stores and external secrets), the webhook (validates ESO resources at admission), and the cert-controller (bootstraps the webhook's serving certificate). This component installs the OPERATOR MACHINERY only: which stores exist and which secrets sync are separate first-class resources (Kubernetes Cluster Secret Store, Kubernetes Secret Store, Kubernetes External Secret). One installation per cluster serves all of them.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Helm Release** -- the `external-secrets` chart from `https://charts.external-secrets.io`, pinned to the chart version you choose (the release name is fixed to `external-secrets`; the CRDs and webhook configuration are cluster-global)
- **Controller, Webhook, and cert-controller** -- the three ESO deployments, each tunable (replicas, resources)
- **CRDs** -- ExternalSecret, SecretStore, ClusterSecretStore, and companions (installed with the release by default, kept on uninstall by default via `helm.sh/resource-policy: keep`)
- **ServiceAccount** -- the controller identity, optionally annotated for keyless cloud access (GKE Workload Identity, EKS IRSA, AKS Workload Identity) — the ambient identity stores fall back to when their auth block is empty
- **Namespace** (optional) -- created with standard governance labels when `create_namespace` is true

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### Kubernetes Cluster

- Any conformant cluster (GKE, EKS, AKS, or self-managed).
- When enabling the Prometheus ServiceMonitor, the Prometheus operator CRDs (e.g. kube-prometheus-stack) must already be on the cluster — the release fails to install without them.
- For the ambient keyless posture, the cluster's workload-identity mechanism must be active (GKE Workload Identity, EKS OIDC provider for IRSA, or AKS Workload Identity) and the cloud-side trust policy must name the controller ServiceAccount.

## Deploy

### Console

Open the deployment store, find **External Secrets Operator**, and click **Deploy**. The creation wizard walks you through installation, controller tuning, multi-tenancy fencing, the ambient workload identity, companion components, observability, and scheduling. Start from the **Minimal** preset for a chart-defaults install in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesExternalSecretsOperator
metadata:
  name: external-secrets
  org: acme-corp
  env: prod
spec:
  namespace:
    value: external-secrets
  createNamespace: true
  chartVersion: 2.8.0
```

```shell
planton apply -f external-secrets-operator.yaml
```

This installs the operator into a new `external-secrets` namespace with CRDs installed and kept on uninstall — the standard single-installation path.

## Key Configuration

These are the most important decisions when configuring the operator. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The namespace is permanent** -- Kept CRDs (the default) retain the Helm release's namespace in their ownership metadata, so re-installing into a DIFFERENT namespace fails on the surviving CRDs. Moving an install means first deleting the kept CRDs — which cascades to every ExternalSecret and SecretStore object cluster-wide. `external-secrets` is the convention.

**Pin the chart version deliberately** -- `chart_version` is also the operator version (chart and app versions are aligned upstream; `2.8.0` ships operator v2.8.0). Pick versions from the chart repository's index (`helm search repo`): the served chart is the contract.

**CRD lifecycle** -- `crds.install` defaults to true (one component managing both halves is strictly simpler); disable only when another installation already owns the CRDs. `crds.keep_on_uninstall` defaults to true, rendered as the `helm.sh/resource-policy: keep` annotation — the chart itself has no keep knob and would DELETE the CRDs, cascading to every ESO object cluster-wide.

**Replicas require leader election** -- one controller replica is standard; with more than one, `leader_elect` is required so exactly one reconciles (the API enforces the coupling). Sync-latency scaling lives in `concurrent` (reconciles in parallel), not replicas.

**The ambient identity** -- `workload_identity` binds the controller ServiceAccount to a cloud identity. Every store whose auth block is empty authenticates through it — the simplest posture when one identity may read every synced secret. Per-store identities (finer-grained, recommended for multi-team clusters) reference dedicated ServiceAccounts in each store's auth block and need nothing here.

**Multi-tenancy fences** -- `controller_class` shards stores between several isolated installations; `scoped_namespace` fences the operator into ONE namespace (ClusterSecretStores become unreachable); `scoped_rbac` swaps the ClusterRole for a namespace Role and requires the scoped namespace. Most clusters leave all three empty.

**The webhook is the apply path** -- an unreachable webhook blocks creating and editing ESO resources (existing secrets keep refreshing). Production clusters often run `webhook.replicas: 2`. Disabling the webhook moves misconfiguration failures from apply time to reconcile time.

**The escape hatch** -- `helm_values` carries additional chart values as a YAML document, merged LAST over everything the typed fields render. For the chart surface beyond the typed fields (e.g. the cert-manager webhook-certificate integration via `webhook.certManager`) — never the substitute for them, and never a place for secrets.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | The namespace the operator is installed into |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `namespace` | Namespace the operator was installed into | Kubernetes Cluster Secret Store's `secrets_namespace` (where cluster-scoped stores' credential Secrets are materialized) |
| `release_name` | Helm release name (always `external-secrets`) | Debugging the release (`helm status`) |
| `controller_service_account` | Controller ServiceAccount name | The identity to bind on the cloud side for the ambient keyless posture (IRSA trust policy subject, GKE Workload Identity member, Azure federated credential subject) |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard install** -- Conventional `external-secrets` namespace, pinned chart, CRDs managed by the release. Start from the **Minimal** preset.

**EKS with ambient IRSA** -- IRSA binds the controller to an IAM role so stores with empty auth blocks read AWS Secrets Manager with zero stored keys. Start from the **EKS Ambient Identity** preset.

**Many teams, many secrets** -- raised reconcile concurrency, two webhook replicas, explicit component sizing. Start from the **Tuned Multi-Team** preset.

## Works With

- **Kubernetes Cluster Secret Store / Kubernetes Secret Store** -- the backend connections; deploy them after the operator (cluster-scoped stores reference its namespace output).
- **Kubernetes External Secret** -- declares each secret to sync through those stores.
- **Kubernetes Cert Manager** -- can own the webhook's serving certificate via Helm values (`webhook.certManager`) in place of the cert-controller.
- **Workloads (Kubernetes Deployment, StatefulSet, ...)** -- consume the materialized Secrets exactly like any other (env valueFrom, volume mounts).
