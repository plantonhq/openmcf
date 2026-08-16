# HTTP allow with a session check

An HTTP Gateway policy that allows a host and requires a fresh Access session every 24 hours. Uses `block_page` / `check_session` -- settings that do not carry the first-apply drift of `add_headers` / `override_ips`.

## When to Use

- Allowlisting an internal web app through Gateway
- Forcing periodic re-authentication on a sensitive host
- An HTTP policy you intend to prove live (idempotency-safe settings only)

## Key Configuration Choices

- **filter: http** -- the module sends `["http"]`. Do not try to combine filters; the spec is singular on purpose.
- **check_session.duration: 24h** -- the provider normalizes this to `24h0m0s` on read. If a plan shows that rewrite, adopt the stored form.
- **block_page.target_uri is required** whenever the block_page object is set.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `account_id` | The Cloudflare account that owns this policy | Cloudflare Dashboard -> Overview -> API section |
| `traffic` | Wirefilter over the HTTP request | Your intranet hostname |
| `rule_settings.block_page.target_uri` | Where blocked/expired sessions land | Your Access block or login page |

## Related Presets

- **01-dns-block** -- DNS-layer block with a block reason
