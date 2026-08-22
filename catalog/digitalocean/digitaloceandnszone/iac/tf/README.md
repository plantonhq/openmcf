# DigitalOcean DNS Zone -- Terraform Module

Deploys a `digitalocean_domain` plus one `digitalocean_record` per managed record value from a `DigitalOceanDnsZone` spec: the zone itself, the create-only `ip_address` convenience, and the inline records with their per-type fields on presence semantics. Provider pin is `~> 2.99`.

`variables.tf` is generated (`planton tofu generate-variables DigitalOceanDnsZone`). Do not hand-edit it. The API token lives in `credentials.tf`.

## Prerequisites

- OpenTofu or Terraform 1.5+
- DigitalOcean API token (`digitalocean_token`)

## Usage

```hcl
module "dns_zone" {
  source = "./path/to/module"

  metadata = {
    name = "example-com"
  }

  spec = {
    domain_name = "example.com"
    records = [
      {
        name        = "@"
        type        = "A"
        values      = ["203.0.113.10"]
        ttl_seconds = 3600
      },
      {
        name     = "@"
        type     = "MX"
        values   = ["aspmx.l.google.com."]
        priority = 1
      }
    ]
  }

  digitalocean_token = var.digitalocean_token
}
```

## Behavior notes

- Record `values` arrive flattened as plain strings — the Planton orchestrator resolves `valueFrom` references before Terraform runs. Each value becomes its own `digitalocean_record` (multi-value entries fan out), keyed `name-index-valueIndex` for stable addressing.
- The spec's shared enum value names ARE the provider's record types, so `type` wires through directly — no lookup table.
- `priority`/`weight`/`port`/`flags` pass through as null when unset; `ttl_seconds: 0`/unset leaves the ttl attribute Computed (DigitalOcean applies its 1800-second default).
- `ip_address` is sent only when non-empty; the provider's Read never returns it (recorded config-only import tolerance).

## Outputs

Exactly the kind's stack-output contract, identical to the Pulumi module: `zone_name`, `zone_id`, `name_servers`, `urn`.
