# AKS Workload Identity

This preset trusts a Kubernetes service account in an AKS cluster to
authenticate as a managed identity -- the modern way for pods to reach Key
Vault, Storage, or any RBAC-granted Azure resource without node-level
credentials or secrets in the cluster.

The cluster's OIDC issuer signs a projected service-account token into
annotated pods; the workload-identity webhook and Azure SDK exchange it for
the identity's credentials. The subject format is fixed by Kubernetes:
`system:serviceaccount:{namespace}:{serviceaccount}`.

The Kubernetes side that consumes it -- the service account annotation:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: <serviceaccount>
  namespace: <namespace>
  annotations:
    azure.workload.identity/client-id: <the identity's status.outputs.client_id>
```

and `azure.workload.identity/use: "true"` on the pod template.

## When to Use

- Pods that read Key Vault secrets, blob data, or Service Bus queues as a
  managed identity
- Replacing the deprecated AAD Pod Identity add-on
- Any in-cluster workload that should hold Azure permissions scoped to its
  own service account, not shared node credentials

## Key Configuration Choices

- **One credential per service account** -- each namespace/service-account
  pair that acts as this identity gets its own credential; the trust is
  auditable per workload
- **Identity per workload** -- prefer a dedicated `AzureUserAssignedIdentity`
  per application over one shared identity with broad grants

## Key Configuration Choices (continued)

- **Issuer by reference** -- the preset wires the issuer to the cluster's
  `oidc_issuer_url` output, so the trust always matches the cluster it is
  deployed beside and never needs a hand-copied URL. For a cluster managed
  outside the environment, replace the reference with a literal
  `issuer: { value: <url> }` (find the URL with
  `az aks show --query oidcIssuerProfile.issuerUrl`).

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<user-assigned-identity-arm-id>` | The parent identity's ARM ID (or use `valueFrom` against an `AzureUserAssignedIdentity`) | The identity's `status.outputs.identity_id` |
| `<aks-cluster-name>` | The `AzureAksCluster` resource whose OIDC issuer signs the tokens (OIDC issuer must be enabled -- it is by default) | Your environment's cluster resource |
| `<namespace>` / `<serviceaccount>` | The Kubernetes workload allowed to authenticate | Your workload manifests |
