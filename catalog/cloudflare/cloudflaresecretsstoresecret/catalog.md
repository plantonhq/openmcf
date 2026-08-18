# Cloudflare Secrets Store Secret

One secret inside the account Secrets Store, readable by the surfaces its scopes allow (Workers, AI Gateway, DEX, Access). The value is write-only at Cloudflare -- never returned, never drift-detected -- and rotating it here rotates it for every consumer at once.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Store secret** -- one `cloudflare_secrets_store_secret` inside the referenced store

## Prerequisites

- **The account Secrets Store** (a `CloudflareSecretsStore` -- one per account)
- **A Cloudflare API token** with Account → Secrets Store → Edit

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareSecretsStoreSecret
metadata:
  name: openai-api-key
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  storeId:
    valueFrom:
      kind: CloudflareSecretsStore
      name: account-secrets
      fieldPath: status.outputs.store_id
  name: openai-api-key
  value:
    value: ${secrets-group/prod-ai/openai-key}
  scopes:
    - ai_gateway
    - workers
```

```shell
planton apply -f secret.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex. |
| `storeId` | StringValueOrRef | The store holding the secret. | Required; references a `CloudflareSecretsStore`. |
| `name` | string | The secret's name -- what consumers reference. | Required; replaces on change. |
| `value` | StringValueOrRef | The secret value (<= 64 KiB). | Required, sensitive, WRITE-ONLY -- provide a managed-secret reference. |
| `scopes` | list | Which surfaces may read it. | Required; from `access`, `ai_gateway`, `dex`, `workers`, listed alphabetically without duplicates (the API returns them sorted -- any other order would drift forever). |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `comment` | string | unset | Free-form note shown in the dashboard. |

## Examples

### Worker-only secret

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareSecretsStoreSecret
metadata:
  name: webhook-signing-key
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  storeId:
    value: "7b0a3d5c1e9f42c68d1a2b3c4d5e6f70"
  name: webhook-signing-key
  value:
    value: ${secrets-group/prod-web/webhook-key}
  scopes:
    - workers
  comment: HMAC key the webhook Worker verifies signatures with.
```

## Destroy Semantics

Destroy is a real delete. Consumers referencing the name fail reads at RUNTIME, not at apply time -- re-point or retire them first.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `secret_id` | string | The secret's ID within its store |
| `store_id` | string | The store holding the secret |

## Related Components

- [Cloudflare Secrets Store](/docs/catalog/cloudflare/cloudflaresecretsstore) -- the vault this secret lives in
- [Cloudflare Worker](/docs/catalog/cloudflare/cloudflareworker) -- secrets-store bindings read it at runtime
- [Cloudflare AI Gateway](/docs/catalog/cloudflare/cloudflareaigateway) -- BYO-keys authentication reads from the same store
