# CloudflareSecretsStoreSecret Terraform Module

Terraform IaC module for one secret inside the account Secrets Store.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareSecretsStoreSecretSpec (generated)
locals.tf     — Empty-string drop (comment)
main.tf       — cloudflare_secrets_store_secret
outputs.tf    — secret_id, store_id
```

## Behavior

The value is write-only at Cloudflare (never returned, never drift-detected) and the provider marks the attribute sensitive, so plans redact it. account_id, store_id, and name are create-only; value, scopes, and comment update in place (a merge-patch). The spec's CEL wall guarantees scopes arrive in Cloudflare's canonical alphabetical order -- the API returns them sorted, and an unsorted config would drift forever. Destroy is a real delete with an honest 404.

## Outputs

| Name | Description |
|------|-------------|
| `secret_id` | The secret's ID within its store |
| `store_id` | The store holding the secret |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
