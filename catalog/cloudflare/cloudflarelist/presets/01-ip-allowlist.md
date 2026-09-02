---
display_name: IP Allowlist
---

# IP Allowlist

An `ip`-kind list to collect trusted IPs/CIDRs that WAF or custom rules reference
with `ip.src in $office_allowlist`.

## When to use

- Allowlisting office/VPN egress IPs, partner ranges, or monitoring probes.
- Any rule that should match a maintained set of addresses by name.

## Key choices

- `kind: ip` — accepts IPv4/IPv6 addresses and CIDRs (immutable).
- `name` — referenced in rule expressions; keep it short and lowercase.
- Add entries with `CloudflareListItem` (one per IP/CIDR).

## Placeholders

| Placeholder | Description |
|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | 32-character Cloudflare account ID |
