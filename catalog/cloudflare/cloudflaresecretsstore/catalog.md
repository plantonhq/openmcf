# Cloudflare Secrets Store

The account-level Secrets Store: the vault Worker bindings and AI Gateway authentication consume secrets from. The store is a named container; the secrets inside it are separate `CloudflareSecretsStoreSecret` resources. Everything here is create-only, and Cloudflare allows one store per account.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Secrets Store** -- one `cloudflare_secrets_store` (the account's vault)

## Prerequisites

- **A Cloudflare account** with no existing store (one per account -- adopt an existing store by import instead)
- **A Cloudflare API token** with Account → Secrets Store → Edit

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareSecretsStore
metadata:
  name: account-secrets
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: account-secrets
```

```shell
planton apply -f secrets-store.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `name` | string | The store's name. | Required; replaces on change -- and takes every secret with it. |

### Optional Fields

None -- the store is a two-field container by design.

## Examples

### The account vault

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareSecretsStore
metadata:
  name: account-secrets
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: account-secrets
```

## Destroy Semantics

Destroy is a real delete and takes every secret in the store with it. Treat the store as permanent infrastructure and destroy it last, if ever.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `store_id` | string | The store's ID -- what secrets, Worker bindings, and AI Gateway authentication reference |

## Related Components

- [Cloudflare Secrets Store Secret](/docs/catalog/cloudflare/cloudflaresecretsstoresecret) -- the secrets inside this store
- [Cloudflare Worker](/docs/catalog/cloudflare/cloudflareworker) -- secrets-store bindings read from here
- [Cloudflare AI Gateway](/docs/catalog/cloudflare/cloudflareaigateway) -- BYO-keys authentication reads from here
