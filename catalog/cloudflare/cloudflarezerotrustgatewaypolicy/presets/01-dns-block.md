---
display_name: DNS Block
---

# DNS block

A Gateway DNS policy that blocks a domain and shows Cloudflare's block page. `enabled` and `precedence` are set explicitly -- omit `enabled: true` and the rule deploys disabled.

## When to Use

- First Gateway policy on an account
- Blocking a small set of domains (or a Zero Trust list of domains) for every WARP-enrolled device
- A template to copy for category-style DNS blocks

## Key Configuration Choices

- **enabled: true** -- Cloudflare defaults this to false. A review that does not see this field should bounce the change.
- **precedence: 1000** -- lower runs earlier. Pick numbers with room to insert policies between them.
- **No add_headers / override_ips** -- those settings drift on first apply at provider v5.23.0. This preset stays on `block_page_enabled` / `block_reason`.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `account_id` | The Cloudflare account that owns this policy | Cloudflare Dashboard -> Overview -> API section |
| `traffic` | Wirefilter over DNS | Cloudflare Gateway expression docs, or `any(dns.domains[*] in $your-list-id)` |

## Related Presets

- **02-http-session-check** -- HTTP allow with a re-auth window
