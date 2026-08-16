# Tiered Caching

Enables Smart Tiered Cache and Regional Tiered Cache on a zone -- the free, recommended production posture. Smart Tiered Cache routes cache misses through the single best upper-tier data center for your origin (the dashboard's Tiered Cache toggle), and Regional Tiered Cache adds a regional tier that cuts long-haul trips for misses in distant regions. Every other cache setting is left unmanaged, so the zone keeps whatever it already carries.

## When to Use

- The default starting point for any production zone: fewer origin requests, better cache hit ratios, no cost
- Zones with a single origin region serving a global audience, where regional tiering shortens miss latency the most
- Adopting managed cache settings on a zone configured by hand -- only these two toggles are touched

## Key Configuration Choices

- **Smart over generic** (`smart_tiered_cache: true`) -- Cloudflare picks the upper tier automatically; the generic `tiered_caching` field exists for rare topology-specific needs and is left unmanaged here.
- **Regional tier on** (`regional_tiered_cache: true`) -- free, and complements smart tiering; note it has no delete at Cloudflare, so destroying the resource abandons the value rather than reverting it.
- **Nothing else managed** -- Cache Reserve, Argo Smart Routing, and cache variants stay untouched; add them deliberately when you need them.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | The Cloudflare zone ID whose cache posture this manages | Cloudflare Dashboard -> zone Overview -> Zone ID (right sidebar), or a `CloudflareDnsZone`'s `status.outputs.zone_id` |

## Related Presets

- **02-image-variants** -- adds WebP/AVIF cache variants for image-heavy zones on top of smart tiering
