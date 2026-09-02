---
display_name: SNI Custom
---

# SNI custom certificate

An SNI-class custom certificate upload with the compatibility-maximizing bundle method. Business/Enterprise zone feature; the certificate must be issued by a publicly trusted CA covering the zone's hostnames.

## When to Use

- Compliance demands a specific issuer or an EV/OV certificate on the zone
- An external PKI team manages the certificate lifecycle
- Prefer `CloudflareCertificatePack` when Cloudflare-managed renewal is acceptable -- custom certificates put renewal on YOUR calendar

## Key Configuration Choices

- **type: sni_custom** -- modern clients, multiple uploads allowed. `legacy_custom` works on every TLS client but occupies the zone's single legacy slot.
- **bundle_method: ubiquitous** -- the default, maximum chain compatibility. `force` sends the chain exactly as uploaded.
- **private_key by reference** -- provide a managed-secret reference; the API never returns the key, and rotation replaces the upload (the certificate id changes).

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone | Cloudflare Dashboard -> zone Overview -> API section, or reference a CloudflareDnsZone via `value_from` |
| `certificate` | The PEM certificate | Your CA's issuance bundle |
| `private_key.value` | The PEM private key | Your secret store -- never paste plaintext into a committed manifest |

## Related Presets

None yet -- a geo-restricted variant belongs on a zone with a measured jurisdiction requirement.
