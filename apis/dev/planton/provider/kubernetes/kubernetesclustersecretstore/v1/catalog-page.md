# KubernetesClusterSecretStore

Creates one External Secrets Operator ClusterSecretStore — a cluster-wide backend connection that ExternalSecrets in ANY namespace can sync from. Full backend surface: AWS Secrets Manager, GCP Secret Manager, Azure Key Vault, Vault/OpenBao, another Kubernetes cluster, and the fake test backend — with keyless ServiceAccount auth, operator ambient identity, or declared static credentials.

## What Gets Created

- **ClusterSecretStore** — named after the resource; the name ExternalSecrets reference with `kind: ClusterSecretStore`
- **Credential Secret** — static keys and tokens declared in the spec are materialized as a `<name>-credentials` Secret in the secrets namespace (never hand-created)

## Prerequisites

- External Secrets Operator on the cluster (**KubernetesExternalSecretsOperator**)
- For keyless backend reads: a ServiceAccount with a workload-identity binding (**KubernetesServiceAccount**), or workload identity configured on the operator itself

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesClusterSecretStore
metadata:
  name: aws-secrets-manager
spec:
  secretsNamespace:
    value: external-secrets
  config:
    aws:
      region: us-east-1
      serviceAccountName:
        value: eso-secrets-reader
      serviceAccountNamespace: external-secrets
```

## Stack Outputs

| Output | Description |
|---|---|
| `store_name` | The store handle ExternalSecrets reference |
| `secrets_namespace` | Where credential Secrets were materialized |

## Next Steps

Create **KubernetesExternalSecret** resources referencing this store's `store_name` output with `kind: ClusterSecretStore`.
