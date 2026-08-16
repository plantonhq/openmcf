# Domain list

A DOMAIN-type Zero Trust list two or more Gateway DNS policies can share. Type is immutable -- changing it later replaces the list and breaks every policy that referenced the old ID.

## When to Use

- A shared blocklist or allowlist of hostnames
- First Zero Trust list on an account
- Prefer this over URL-type lists (URL values drift at provider v5.23.0)

## Key Configuration Choices

- **type: DOMAIN** -- uppercase. Lowercase would round-trip as permanent plan drift.
- **items are a set** -- order is not preserved. Do not depend on position.
- **value is required** -- an entry with only a description is rejected.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `account_id` | The Cloudflare account that owns this list | Cloudflare Dashboard -> Overview -> API section |
| `items[].value` | Domain names to share across policies | Your block/allow inventory |

## Related Presets

- **02-ip** -- CIDRs and addresses instead of hostnames
