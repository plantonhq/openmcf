# Image Variants

Enables Smart Tiered Cache and configures cache variants so Cloudflare may serve WebP and AVIF versions of cached JPEG and PNG assets to clients that accept them. Pair this with an image optimization pipeline (such as Polish with WebP) that actually produces the variants -- this setting tells the cache which alternate content types are acceptable per original extension, it does not convert images itself. Cache variants are one of the two settings with a real delete: destroying this resource resets variants to Cloudflare defaults.

## When to Use

- Image-heavy zones (media sites, storefronts, galleries) serving `.jpg`, `.jpeg`, and `.png` assets to modern browsers
- Zones running Polish or an equivalent pipeline that generates WebP/AVIF variants worth caching
- Cutting image bandwidth without touching application code -- the variant negotiation happens at the cache

## Key Configuration Choices

- **Per-extension MIME lists** (`cache_variants.jpg`, `.jpeg`, `.png`) -- each key is the original asset's extension; the list is the variant content types the cache may serve for it. Order lists both WebP and AVIF so either variant can be cached.
- **Only listed extensions are managed** -- the other eight supported extensions (`avif`, `bmp`, `gif`, `jp2`, `jpg2`, `tif`, `tiff`, `webp`) are omitted and stay at their current values.
- **Smart tiering included** (`smart_tiered_cache: true`) -- variants and tiered caching compound: fewer origin fetches, and the fetches that happen can serve multiple formats from one cached entry.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d` | The Cloudflare zone ID whose cache posture this manages | Cloudflare Dashboard -> zone Overview -> Zone ID (right sidebar), or a `CloudflareDnsZone`'s `status.outputs.zone_id` |

## Related Presets

- **01-tiered-caching** -- the base production posture without variants; start there if the zone is not image-heavy
