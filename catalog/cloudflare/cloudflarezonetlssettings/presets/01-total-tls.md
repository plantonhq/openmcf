# Total TLS with Google Trust Services

Enables Total TLS for the zone so every proxied hostname gets its own certificate, including deep subdomains that Universal SSL's wildcard does not cover. Certificates are issued by Google Trust Services and auto-renewed by Cloudflare on a fixed 90-day validity. Automatic origin TLS key exchange is also enabled, letting Cloudflare negotiate the strongest key exchange the origin supports. Everything else in the zone's TLS posture is left unmanaged and keeps its current values.

## When to Use

- You serve proxied hostnames two or more levels deep (like `app.eu.example.com`) that Universal SSL's `*.example.com` wildcard does not cover.
- You want per-hostname certificates issued and renewed automatically instead of ordering certificate packs by hand.
- Note: Total TLS requires the zone's Advanced Certificate Manager subscription -- without it the API rejects the write with 401 code 1450.

## Key Configuration Choices

- `totalTls.enabled: true` turns on per-hostname certificate issuance. This has no delete at Cloudflare: destroying the resource abandons the setting rather than turning it off.
- `totalTls.certificateAuthority: google` picks Google Trust Services. Other options are `lets_encrypt` and `ssl_com`; leave the field unset to let Cloudflare choose. Validity is fixed at 90 days and is not configurable.
- `autoOriginTlsKex: true` strengthens origin-facing TLS with no origin changes needed. Also no delete at Cloudflare.
- Universal SSL is deliberately not touched here. Total TLS supplements it; you rarely need to disable Universal SSL, and doing so can make hostnames unreachable over HTTPS.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | The 32-character ID of the zone whose TLS posture this manages | Cloudflare dashboard, zone Overview page, API section; or reference a `CloudflareDnsZone` resource via `valueFrom` instead |

## Related Presets

- [02-hostname-overrides](02-hostname-overrides.md) -- raise the TLS floor for individual hostnames instead of managing zone-wide issuance.
