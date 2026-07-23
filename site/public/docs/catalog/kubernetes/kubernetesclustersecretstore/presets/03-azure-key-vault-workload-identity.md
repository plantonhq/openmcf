---
title: "Azure Key Vault with Workload Identity (Keyless)"
description: "This preset creates a cluster-wide store connected to Azure Key Vault, authenticating through a referenced ServiceAccount whose AKS Workload Identity federation authorizes the reads. No client secret..."
type: "preset"
rank: "03"
presetSlug: "03-azure-key-vault-workload-identity"
componentSlug: "kubernetesclustersecretstore"
componentTitle: "KubernetesClusterSecretStore"
provider: "kubernetes"
icon: "package"
order: 3
---

# Azure Key Vault with Workload Identity (Keyless)

This preset creates a cluster-wide store connected to Azure Key Vault, authenticating through a referenced ServiceAccount whose AKS Workload Identity federation authorizes the reads. No client secret touches the cluster — the production posture on AKS.

## When to Use

- Your secrets live in Azure Key Vault and your cluster runs on AKS with Workload Identity enabled
- You want keyless authentication — no service-principal client secrets
- Every team should sync from one platform-wide store (add `conditions` to fence namespaces if not)

## Key Configuration Choices

- **Workload Identity auth** (`authType: WorkloadIdentity`) -- the default and the keyless posture; `ManagedIdentity` covers AKS without Workload Identity, `ServicePrincipal` is the static-credential fallback for clusters outside Azure
- **Keyless ServiceAccount auth** (`serviceAccountName`) -- references a ServiceAccount with an AKS Workload Identity binding; create it with `KubernetesServiceAccount` (its workload-identity arm carries the Entra client-id annotation) and reference its `service_account_name` output
- **Explicit ServiceAccount namespace** -- required on a ClusterSecretStore; cluster scope has no home namespace to default to
- **Tenant ID** (`tenantId`) -- the Entra (Azure AD) tenant the federated identity lives in

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `https://my-key-vault.vault.azure.net` | Your Key Vault URL (not a placeholder — replace the whole value) | Azure portal > Key Vault > Overview |
| `<entra-tenant-id>` | Entra (Azure AD) tenant ID | Azure portal > Microsoft Entra ID > Overview |
| `<eso-reader-service-account>` | ServiceAccount with a Workload Identity federation allowing Key Vault secret reads | `KubernetesServiceAccount` outputs or `kubectl get sa -n external-secrets` |

## Related Presets

- **01-aws-secrets-manager-irsa** -- Use when secrets live in AWS Secrets Manager
- **02-gcp-secret-manager-workload-identity** -- Use when secrets live in GCP Secret Manager
- **04-vault-kubernetes-auth** -- Use when secrets live in Vault or OpenBao
