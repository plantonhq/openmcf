---
title: "Vault KV with Kubernetes Auth (Keyless)"
description: "This preset creates a cluster-wide store connected to a HashiCorp Vault KV v2 engine, authenticating through Vault's Kubernetes auth method: the referenced ServiceAccount's token is exchanged for a..."
type: "preset"
rank: "04"
presetSlug: "04-vault-kubernetes-auth"
componentSlug: "kubernetesclustersecretstore"
componentTitle: "KubernetesClusterSecretStore"
provider: "kubernetes"
icon: "package"
order: 4
---

# Vault KV with Kubernetes Auth (Keyless)

This preset creates a cluster-wide store connected to a HashiCorp Vault KV v2 engine, authenticating through Vault's Kubernetes auth method: the referenced ServiceAccount's token is exchanged for a Vault token, so no Vault credential is stored on the cluster. **It works unchanged for OpenBao** — OpenBao speaks the Vault API; point `server` at the OpenBao endpoint. This is the production posture for in-cluster or company-internal Vault/OpenBao.

## When to Use

- Your secrets live in HashiCorp Vault or OpenBao (KV engine)
- The Vault server has the Kubernetes auth method configured to trust this cluster
- You want keyless authentication — no static Vault tokens (which expire) and no AppRole secret-ids on the cluster

## Key Configuration Choices

- **Kubernetes auth** (`kubernetes` block) -- the cluster's ServiceAccount token is exchanged for a Vault token; the `role` must be bound to the referenced ServiceAccount on the Vault side
- **KV v2** (`version: v2`) -- the versioned, modern KV engine (the default); switch to `v1` only for legacy mounts
- **Mount paths** -- `path: secret` and `mountPath: kubernetes` are the upstream defaults; adjust if your engines are mounted elsewhere
- **ServiceAccount reference** (`serviceAccountName`) -- the identity presented to Vault; leave it out to present the operator's own ServiceAccount instead

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `https://vault.example.internal:8200` | Your Vault/OpenBao server URL (not a placeholder — replace the whole value) | Your Vault/OpenBao deployment |
| `<vault-role-name>` | Vault Kubernetes-auth role bound to the ServiceAccount | `vault read auth/kubernetes/role/<name>` (or the OpenBao equivalent) |
| `<eso-reader-service-account>` | ServiceAccount whose token is presented to Vault | `KubernetesServiceAccount` outputs or `kubectl get sa -n external-secrets` |

## Related Presets

- **01-aws-secrets-manager-irsa** -- Use when secrets live in AWS Secrets Manager
- **02-gcp-secret-manager-workload-identity** -- Use when secrets live in GCP Secret Manager
- **03-azure-key-vault-workload-identity** -- Use when secrets live in Azure Key Vault
