# CloudflareDnsRecord Terraform Module

This Terraform module provisions a Cloudflare DNS record.

## Prerequisites

- Terraform 1.0+
- Cloudflare API token with DNS:Edit permissions

## Usage

### As Part of OpenMCF

This module is typically invoked through the OpenMCF CLI:

```bash
planton apply -f manifest.yaml --iac terraform
```

### Standalone Usage

```hcl
module "dns_record" {
  source = "./path/to/module"

  metadata = {
    name = "www-a-record"
  }

  spec = {
    zone_id = "your-zone-id"
    name    = "www"
    type    = "A"
    value   = "192.0.2.1"
    proxied = true
    ttl     = 1
    comment = "Primary web server"
  }
}
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `CLOUDFLARE_API_TOKEN` | Cloudflare API token | Yes |

## Inputs

### metadata

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | string | Yes | Resource name |
| `id` | string | No | Optional resource ID |
| `org` | string | No | Organization |
| `env` | string | No | Environment |
| `labels` | map(string) | No | Labels |
| `tags` | list(string) | No | Tags |

### spec

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `zone_id` | string | Yes | - | Cloudflare Zone ID |
| `name` | string | Yes | - | DNS record name |
| `type` | string | Yes | - | Record type (A, AAAA, CNAME, MX, TXT, SRV, NS, CAA) |
| `value` | string | Yes | - | Record value |
| `proxied` | bool | No | false | Proxy through Cloudflare |
| `ttl` | number | No | 1 | TTL in seconds (1 = auto) |
| `priority` | number | No | 0 | Priority for MX/SRV |
| `comment` | string | No | "" | Comment (max 100 chars) |

## Outputs

| Name | Description |
|------|-------------|
| `record_id` | Cloudflare DNS record ID |
| `hostname` | Fully qualified hostname |
| `record_type` | DNS record type |
| `proxied` | Whether record is proxied |

## Examples

### A Record

```hcl
spec = {
  zone_id = "abc123"
  name    = "www"
  type    = "A"
  value   = "192.0.2.1"
  proxied = true
}
```

### MX Record

```hcl
spec = {
  zone_id  = "abc123"
  name     = "@"
  type     = "MX"
  value    = "mail.example.com"
  priority = 10
}
```

### TXT Record

```hcl
spec = {
  zone_id = "abc123"
  name    = "@"
  type    = "TXT"
  value   = "v=spf1 include:_spf.google.com ~all"
}
```

## Troubleshooting

### "authentication failed"

Ensure `CLOUDFLARE_API_TOKEN` environment variable is set with a valid token.

### "zone not found"

Verify the `zone_id` matches an existing Cloudflare zone.

### "invalid record type"

Ensure `type` is one of: A, AAAA, CNAME, MX, TXT, SRV, NS, CAA.

## Validation

```bash
terraform init
terraform validate
```
