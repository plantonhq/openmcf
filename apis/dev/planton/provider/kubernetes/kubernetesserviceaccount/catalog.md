# Kubernetes ServiceAccount

Deploys a Kubernetes ServiceAccount -- the identity pods run as -- with inherited image pull secrets, a tri-state automatic API token mount, and optional keyless federation to a GCP service account (GKE Workload Identity), AWS IAM role (EKS IRSA), or Azure managed identity (AKS Workload Identity). Manages workload identity declaratively through a Kubernetes Provider Connection with full audit trail and versioning.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Kubernetes ServiceAccount** -- a single ServiceAccount in the specified namespace, with `imagePullSecrets` attached and `automountServiceAccountToken` applied only when a position was taken
- **Cloud identity annotation** -- written automatically when a workload-identity arm is configured: `iam.gke.io/gcp-service-account` (GKE), `eks.amazonaws.com/role-arn` (EKS), or `azure.workload.identity/client-id` (+ optional tenant) for AKS
- **Kubernetes Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking

## Before You Deploy

### Planton Setup

- **Kubernetes Provider Connection** -- an active connection in the Connect module with kubeconfig credentials for the target cluster. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline kubeconfig authentication.

### Kubernetes Cluster

- The target namespace must already exist (the module does not create it). Use the Kubernetes Namespace component to manage namespaces declaratively.
- For workload identity, the CLOUD side must trust this exact namespace/ServiceAccount pair: a `workloadIdentityUser` IAM binding (GKE), an OIDC trust policy on the role (EKS), or a federated credential on the managed identity (AKS).

## Deploy

### Console

Open the deployment store, find **ServiceAccount on Kubernetes**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Basic** preset for a plain per-app identity or a **Workload Identity** preset for keyless cloud access in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesServiceAccount
metadata:
  name: checkout-identity
  org: acme-corp
  env: prod
spec:
  name: checkout-identity
  namespace:
    value: backend-services
  automountServiceAccountToken: false
```

```shell
planton apply -f service-account.yaml
```

This creates a hardened per-app identity in the `backend-services` namespace: pods running as it get no automatic API server token.

## Key Configuration

These are the most important decisions when configuring a Kubernetes ServiceAccount. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**One identity per app** -- Workloads reference it via `serviceAccountName`, RBAC grants bind permissions to it, and cloud federation attaches to it. Separate identities per app keep least privilege auditable.

**Image pull secrets are inherited** -- Secret NAMES listed in `imagePullSecrets` (never the credential material) are attached to every pod running as this identity -- private-registry access configured once instead of per workload.

**API token mount is a tri-state** -- Unset means the cluster default (the token mounts); `false` hardens every pod running as this identity (most app workloads never call the Kubernetes API); `true` forces the mount. The unset state is a real position -- the module writes the field only when one was taken.

**Workload identity is keyless** -- Exactly one arm (GKE / EKS / AKS) federates this Kubernetes identity to a cloud identity: pods call cloud APIs with its permissions and zero stored keys. The AKS arm optionally carries a tenant ID for cross-tenant identities.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Purpose |
|-------|-----------|---------|
| `spec.namespace` | KubernetesNamespace (`spec.name`) | The namespace the identity is scoped to; omitted means the cluster's `default` namespace |
| `spec.imagePullSecrets[]` | KubernetesSecret (`spec.name`) | Docker-registry Secret names every pod with this identity inherits |
| `spec.workloadIdentity.gke.serviceAccountEmail` | GcpServiceAccount (`status.outputs.email`) | The Google service account this identity federates to |
| `spec.workloadIdentity.eks.roleArn` | AwsIamRole (`status.outputs.role_arn`) | The IAM role this identity federates to |
| `spec.workloadIdentity.aks.clientId` | AzureUserAssignedIdentity (`status.outputs.client_id`) | The managed identity this identity federates to |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `service_account_name` | The name of the created ServiceAccount | Workload `serviceAccountName` references |
| `namespace` | The namespace the identity lives in | Verifying workload co-location |
| `rbac_subject` | The ready-to-use subject string (`system:serviceaccount:<namespace>:<name>`) | RBAC bindings inside and outside the platform |
| `workload_identity_handle` | The federated cloud identity (GSA email, role ARN, or client ID); empty when not configured | Cloud-side trust verification and audit |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Per-app identity** -- A plain ServiceAccount named for the app, referenced by its workload and its RBAC grant. Start from the **Basic** preset.

**Keyless cloud access** -- The GKE / EKS-IRSA / AKS arm federates the identity to a cloud principal; pods call cloud APIs with zero stored keys. Start from the matching **Workload Identity** preset.

**Private registry access** -- Registry credentials attached once via `imagePullSecrets`; every pod running as this identity pulls without per-pod configuration. Start from the **Image Pull Secrets** preset.

## Works With

- **Kubernetes Namespace** -- reference the namespace so infra charts create it and this identity in dependency order.
- **Kubernetes Secret** -- docker-registry Secrets attach via `imagePullSecrets`; the Secret kind's `serviceAccountToken` variant mints long-lived tokens FOR this identity.
- **Kubernetes RBAC** -- grants bind Kubernetes permissions to this identity as a ServiceAccount subject.
- **Kubernetes Deployment and the other workload kinds** -- run as this identity via `serviceAccountName`.
