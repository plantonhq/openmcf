# Bulk Redirect List

A `redirect`-kind list holding source→target URL rules. A redirect ruleset
(`CloudflareRuleset`, http_request_redirect phase) resolves these with `from_list`,
enabling large-scale URL redirects managed as data.

## When to use

- Migrating URLs at scale (site relaunch, domain consolidation).
- Marketing vanity URLs that map to canonical destinations.

## Key choices

- `kind: redirect` — entries are redirect definitions (immutable).
- Add entries with `CloudflareListItem` using the `redirect` item shape
  (`sourceUrl`, `targetUrl`, optional status code and matching flags).
- Wire a `CloudflareRuleset` redirect rule's `from_list` to this list's name.

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
