# KubernetesSecretStore

Creates one External Secrets Operator SecretStore — a namespaced backend connection that only ExternalSecrets in the SAME namespace can sync from. Full backend surface: AWS Secrets Manager, GCP Secret Manager, Azure Key Vault, Vault/OpenBao, another Kubernetes cluster, and the fake test backend — with keyless ServiceAccount auth, operator ambient identity, or declared static credentials. The per-team grain: credentials and blast radius end at the namespace boundary.

## What Gets Created

- **SecretStore** — named after the resource, in the spec's namespace; the name ExternalSecrets in that namespace reference (kind SecretStore, the upstream default)
- **Credential Secret** — static keys and tokens declared in the spec are materialized as a `<name>-credentials` Secret in the store's namespace (never hand-created)

## Prerequisites

- External Secrets Operator on the cluster (**KubernetesExternalSecretsOperator**)
- For keyless backend reads: a ServiceAccount with a workload-identity binding (**KubernetesServiceAccount**), or workload identity configured on the operator itself

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesSecretStore
metadata:
  name: team-gcp-secrets
spec:
  namespace:
    value: team-a
  config:
    gcpSecretManager:
      projectId:
        value: my-gcp-project
      serviceAccountName:
        value: team-a-secrets-reader
```

## Stack Outputs

| Output | Description |
|---|---|
| `store_name` | The store handle ExternalSecrets in the same namespace reference |
| `namespace` | Where the store and its credential Secrets live |

## Next Steps

Create **KubernetesExternalSecret** resources in the same namespace referencing this store's `store_name` output.
