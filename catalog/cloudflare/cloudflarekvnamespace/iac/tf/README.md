# Cloudflare KV Namespace - OpenTofu Module

Provisions a Cloudflare Workers KV namespace (`cloudflare_workers_kv_namespace`) and
exports its identifier for binding to Workers.

## Inputs

| Variable | Description |
|----------|-------------|
| `metadata` | Resource metadata (name, labels, ...). |
| `spec.namespace_name` | Human-readable title for the KV namespace (max 64 chars), unique within the account. |
| `spec.account_id` | Cloudflare account ID (32 hex characters) that owns the namespace. |

## Outputs

| Output | Description |
|--------|-------------|
| `namespace_id` | The unique identifier of the created KV namespace. |
| `supports_url_encoding` | Whether the namespace supports URL-encoded key access. |

## Provider

```hcl
terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5.23"
    }
  }
}

provider "cloudflare" {
  # Automatically uses CLOUDFLARE_API_TOKEN environment variable
}
```
