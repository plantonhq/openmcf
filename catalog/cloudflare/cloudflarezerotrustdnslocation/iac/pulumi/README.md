# CloudflareZeroTrustDnsLocation Pulumi Module

Pulumi (Go) IaC module for Gateway DNS locations.

## Architecture

```
main.go                   — Entrypoint loading the stack input
module/main.go            — Resources(): provider setup, resource, outputs
module/locals.go          — Locals initialization
module/dns_location.go    — cloudflare.ZeroTrustDnsLocation
module/outputs.go         — Stack output keys
```

## Behavior

Mirrors the Terraform module's contract exactly: a plain CRUD resource (real create/update/delete; only the account forces replacement). Update is a full replace at the API — `max_ttl` is sent whenever the spec declares it, and omitting it genuinely resets the TTL behavior to inherit. `dns_destination_ips_id` is only sent when set; unset lets Cloudflare auto-assign the shared IPv4 destination pair.

## Outputs

| Name | Description |
|------|-------------|
| `location_id` | The Cloudflare-assigned UUID of the location |
| `doh_subdomain` | The location's unique DNS-over-HTTPS subdomain |
| `ip` | The IPv4 destination assigned to the plain-DNS endpoint |

## Provider Version

Uses `pulumi-cloudflare` SDK v6 (Cloudflare Terraform provider v5 line).
