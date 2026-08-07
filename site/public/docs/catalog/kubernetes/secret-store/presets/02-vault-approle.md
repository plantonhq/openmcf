---
title: "Team-Scoped Vault KV with AppRole"
description: "This preset creates a namespaced store connected to a HashiCorp Vault (or OpenBao — same API) KV v2 engine, authenticating with AppRole machine identity: a role-id plus a secret-id. The secret-id is..."
type: "preset"
rank: "02"
presetSlug: "02-vault-approle"
componentSlug: "secret-store"
componentTitle: "Secret Store"
provider: "kubernetes"
icon: "package"
order: 2
---

# Team-Scoped Vault KV with AppRole

This preset creates a namespaced store connected to a HashiCorp Vault (or OpenBao — same API) KV v2 engine, authenticating with AppRole machine identity: a role-id plus a secret-id. The secret-id is a declared credential — the module materializes it as a `<name>-credentials` Kubernetes Secret in the team's namespace and the store references it; it never appears in the CR. Use this when the Vault server has no Kubernetes auth trust with this cluster (Vault outside the platform, cluster outside Vault's trusted issuers).

## When to Use

- A team's secrets live in Vault/OpenBao and only that team's namespace should sync them
- Vault's Kubernetes auth method is not configured to trust this cluster — AppRole is the standard machine-identity fallback
- You need per-team Vault identities: each team's store carries its own AppRole

## Key Configuration Choices

- **AppRole auth** (`appRole` block) -- role-id (public half) + secret-id (secret half); prefer Vault Kubernetes auth when the trust exists, since AppRole secret-ids must be issued and rotated
- **Declared credential** (`secretId`) -- materialized by the module as a Kubernetes Secret in the team's namespace; readable only there (the namespaced blast radius)
- **KV v2** (`version: v2`) -- the versioned, modern KV engine (the default); switch to `v1` only for legacy mounts
- **Mount paths** -- `path: secret` and `path: approle` are the upstream defaults; adjust if your engines are mounted elsewhere

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<team-namespace>` | The team's namespace (store, credential Secret, and ExternalSecrets all live here) | `kubectl get namespaces` or `KubernetesNamespace` outputs |
| `https://vault.example.internal:8200` | Your Vault/OpenBao server URL (not a placeholder — replace the whole value) | Your Vault/OpenBao deployment |
| `<vault-approle-role-id>` | The AppRole role ID | `vault read auth/approle/role/<name>/role-id` |
| `<vault-approle-secret-id>` | The AppRole secret ID | `vault write -f auth/approle/role/<name>/secret-id` |

## Related Presets

- **01-team-gcp-secret-manager** -- Use when the team's secrets live in GCP Secret Manager
- **03-fake-sandbox** -- Use for pipelines and sandboxes, without any external account
