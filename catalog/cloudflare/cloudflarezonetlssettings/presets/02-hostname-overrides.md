# Per-Hostname TLS Overrides

Raises a single hostname's minimum TLS version to 1.3 and enables HTTP/2 for it, while the rest of the zone stays on the zone-wide settings. Each override becomes its own per-hostname object at Cloudflare, so this preset changes exactly one hostname and nothing else. These overrides have a real delete: destroying the resource removes them and the hostname falls back to the zone-wide settings.

## When to Use

- You want a strict TLS floor on a security-sensitive hostname (an API or auth endpoint) without forcing TLS 1.3 on the whole zone and breaking older clients elsewhere.
- You need per-hostname control that `CloudflareZoneSettings`' zone-wide minimum TLS version cannot give you.

## Key Configuration Choices

- `minTlsVersion: "1.3"` sets the floor for this hostname only. Valid values are `"1.0"` through `"1.3"`, always quoted so YAML does not read them as numbers.
- `http2: true` overrides the zone-wide HTTP/2 setting for this hostname.
- Add more rows to `hostnameSettings` for more hostnames; each row must set at least one of `minTlsVersion`, `http2`, or `ciphers`. Editing one row never disturbs another hostname's overrides.
- Renaming a hostname replaces its override objects, since the hostname is part of their API identity.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | The 32-character ID of the zone the hostname belongs to | Cloudflare dashboard, zone Overview page, API section; or reference a `CloudflareDnsZone` resource via `valueFrom` instead |
| `api.example.com` | The hostname the overrides apply to; must be a hostname in this zone | Your zone's DNS records |

## Related Presets

- [01-total-tls](01-total-tls.md) -- manage zone-wide certificate issuance with Total TLS instead of per-hostname behavior.
