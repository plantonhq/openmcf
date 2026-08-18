# CloudflareAuthenticatedOriginPulls Terraform Module

Terraform IaC module for a zone's Authenticated Origin Pulls enablement: the zone-wide toggle plus per-hostname certificate associations.

## Architecture

```
provider.tf   — Cloudflare provider configuration (~> 5.23)
variables.tf  — Input variables mirroring CloudflareAuthenticatedOriginPullsSpec
locals.tf     — Association fan-out map (one resource per hostname row)
main.tf       — cloudflare_authenticated_origin_pulls_settings (count-gated)
                + cloudflare_authenticated_origin_pulls (for_each per row)
outputs.tf    — zone_id
```

## Behavior

The zone toggle is managed only when `zone_enabled` is set (presence, not truthiness). Association rows fan out one provider resource per hostname -- the provider hard-fails config lists that do not hold exactly one element. An omitted row `enabled` is sent as true (Cloudflare treats null as "void the association"). Destroy deletes nothing: the toggle is abandoned (no delete exists at Cloudflare) and associations revert by a null write issued from state.

## Outputs

| Name | Description |
|------|-------------|
| `zone_id` | The zone whose Authenticated Origin Pulls surface is managed |

## Provider Version

Uses `cloudflare/cloudflare ~> 5.23`.
