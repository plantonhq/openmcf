# Zone-wide toggle

Enable Authenticated Origin Pulls for the whole zone using Cloudflare's zone-level client certificate. The origin must REQUIRE and VALIDATE the certificate for the control to protect anything -- Cloudflare's side alone only presents it.

## When to Use

- First step of locking the origin down to Cloudflare-only traffic
- Zones where one client certificate for all hostnames is enough (per-hostname pinning is the associations list)

## Key Configuration Choices

- **zone_enabled: true** -- presence matters: unset means "leave the toggle alone", explicit false asserts OFF. Destroy abandons the live value (no delete exists at Cloudflare) -- assert false BEFORE destroying when AOP must end OFF.
- **No hostname_associations** -- add rows pinning hostnames to uploaded certificates (`CloudflareAuthenticatedOriginPullsCertificate`) when different hostnames need different client certificates.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone | Cloudflare Dashboard -> zone Overview -> API section, or reference a CloudflareDnsZone via `value_from` |

## Related Presets

None yet -- a per-hostname pinning variant belongs with an uploaded client certificate to reference.
