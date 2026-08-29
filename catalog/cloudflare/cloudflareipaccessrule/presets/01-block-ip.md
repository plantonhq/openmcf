---
display_name: Block IP
---

# Block an IP

An account-wide IP Access rule that blocks a single IPv4 address on every zone. Only `mode` and `notes` update in place -- changing the address later requires delete-and-recreate.

## When to Use

- A known-bad address that should be refused everywhere on the account
- First IP Access rule on an account
- Prefer this over a zone-scoped rule when the block must follow every zone

## Key Configuration Choices

- **account_id, not zone_id** -- exactly one scope. Setting both is rejected; the provider would silently prefer the account.
- **mode: block** -- refuse the request outright. `managed_challenge` is the less-destructive alternative.
- **target: ip** -- a plain IPv4 address. Use `ip6` (fully-expanded long form) for IPv6, `ip_range` for CIDR.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `account_id` | The Cloudflare account that owns this rule | Cloudflare Dashboard -> Overview -> API section |
| `configuration.value` | The IPv4 address to block | Your threat or allow/block inventory |

## Related Presets

- **02-challenge-country** -- a zone-scoped managed challenge on a country code
