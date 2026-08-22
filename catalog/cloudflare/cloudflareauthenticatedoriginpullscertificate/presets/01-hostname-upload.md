# Hostname-scoped client certificate

A hostname-scoped Authenticated Origin Pulls client-certificate upload. Nothing changes at the edge until a `CloudflareAuthenticatedOriginPulls` association pins a hostname to the resulting `certificate_id` -- the safe upload path. Self-signed pairs are the designed case (your origin validates them).

## When to Use

- Per-hostname client certificates (multi-tenant origins, different trust per hostname)
- Staging new AOP material before any hostname is pinned to it
- Prefer `scope: zone` only when deliberately replacing the zone-wide certificate for every origin pull

## Key Configuration Choices

- **scope: hostname** -- upload-only blast radius; associations opt hostnames in explicitly. `zone` replaces the shared certificate for the whole zone immediately.
- **private_key by reference** -- the API never returns it; your secret store is the system of record.
- **Rotate certificate and key TOGETHER** -- the zone-scoped surface silently ignores key-only changes at provider v5.23.0 (its Update is empty); one rotation habit for both scopes.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone | Cloudflare Dashboard -> zone Overview -> API section, or reference a CloudflareDnsZone via `value_from` |
| `certificate` | The client certificate PEM | Your PKI or `openssl req -x509 -new ...` (self-signed is normal) |
| `private_key.value` | The PEM private key | Your secret store -- never paste plaintext into a committed manifest |

## Related Presets

None yet -- a zone-scope variant belongs on a zone whose entire origin fleet shares one client certificate.
