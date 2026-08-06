# AWS Secrets Manager with IRSA (Keyless)

This preset creates a cluster-wide store connected to AWS Secrets Manager, authenticating through a referenced ServiceAccount whose IRSA (IAM Roles for Service Accounts) binding authorizes the reads. No static AWS keys touch the cluster — the production posture on EKS, and the most common ESO store configuration overall.

## When to Use

- Your secrets live in AWS Secrets Manager and your cluster runs on EKS
- You want keyless authentication — no long-lived AWS credentials on the cluster
- Every team should sync from one platform-wide store (add `conditions` to fence namespaces if not)

## Key Configuration Choices

- **Secrets Manager service** (`service: SecretsManager`) -- the default; switch to `ParameterStore` for SSM parameters
- **Keyless ServiceAccount auth** (`serviceAccountName`) -- references a ServiceAccount with an IRSA role binding; create it with `KubernetesServiceAccount` (its workload-identity arm carries the IAM role annotation) and reference its `service_account_name` output
- **Explicit ServiceAccount namespace** -- required on a ClusterSecretStore; cluster scope has no home namespace to default to
- **Credential Secrets namespace** (`secretsNamespace: external-secrets`) -- the operator's install namespace by convention; this keyless preset materializes no credentials, but the namespace anchors any later ones

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<aws-region>` | AWS region the secrets live in (e.g. `us-east-1`) | AWS Secrets Manager console |
| `<eso-reader-service-account>` | ServiceAccount with an IRSA binding allowing `secretsmanager:GetSecretValue` | `KubernetesServiceAccount` outputs or `kubectl get sa -n external-secrets` |

## Related Presets

- **02-gcp-secret-manager-workload-identity** -- Use when secrets live in GCP Secret Manager
- **03-azure-key-vault-workload-identity** -- Use when secrets live in Azure Key Vault
- **04-vault-kubernetes-auth** -- Use when secrets live in Vault or OpenBao
