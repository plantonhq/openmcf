# Cloudflare Cache Settings

Manages a Cloudflare zone's caching and performance posture as one resource: Smart Tiered Cache, generic tiered caching, Cache Reserve, Regional Tiered Cache, Argo Smart Routing, and cache variants. A field left unset is not managed -- the module never sends it, so the zone keeps whatever value it already carries. Note the destroy semantics: only Smart Tiered Cache and cache variants have a real delete at Cloudflare; the other four settings keep their last-applied value after destroy, and Argo Smart Routing keeps billing until turned off explicitly.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions one zone-scoped API object per managed setting:

- **Tiered Cache (smart)** -- created only when `smartTieredCache` is set; the dashboard's Tiered Cache toggle. Real delete: destroy disables it.
- **Argo Tiered Caching (generic)** -- created only when `tieredCaching` is set. No delete at Cloudflare.
- **Zone Cache Reserve** -- created only when `cacheReserve` is set. Paid feature; no delete at Cloudflare.
- **Regional Tiered Cache** -- created only when `regionalTieredCache` is set. No delete at Cloudflare.
- **Argo Smart Routing** -- created only when `argoSmartRouting` is set. Paid feature; no delete at Cloudflare -- destroy keeps billing if left on.
- **Zone Cache Variants** -- created only when `cacheVariants` is set; only extensions with at least one MIME type are sent. Real delete: destroy resets variants to defaults.

## Prerequisites

- **Cloudflare Provider Connection** -- an active connection with an API token that has Cache Settings edit access on the target zone (Argo Smart Routing additionally requires the Argo edit permission).
- **An active zone** -- the zone must already exist; manage it with a `CloudflareDnsZone` and reference its `zone_id` output, or pass an externally-managed zone ID directly.
- **Plan entitlements for paid features** -- `cacheReserve` and `argoSmartRouting` are paid Cloudflare features; the account must have them available before the apply can enable them.

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareCacheSettings
metadata:
  name: prod-cache
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: prod-zone
      fieldPath: status.outputs.zone_id
  smartTieredCache: true
```

```shell
planton apply -f cache-settings.yaml
```

This enables Smart Tiered Cache on the zone and manages nothing else -- every other setting stays exactly as it is.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone whose cache settings are managed. Can reference a `CloudflareDnsZone` resource via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. At least one of the six settings below must also be set -- a resource that manages nothing is rejected at validation. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `smartTieredCache` | bool | not managed | Route cache misses through the single best upper-tier data center for the origin (the dashboard's Tiered Cache toggle). Real delete: destroy disables it. |
| `tieredCaching` | bool | not managed | Generic tiered cache topology where lower-tier data centers ask upper tiers before the origin. Use `smartTieredCache` unless a specific generic topology is required. No delete at Cloudflare. |
| `cacheReserve` | bool | not managed | Persist cached objects in Cloudflare's durable storage tier so they survive eviction. Paid, billed by storage and operations. No delete at Cloudflare -- turn it off explicitly before retiring. |
| `regionalTieredCache` | bool | not managed | Add a regional tier between lower tiers and the upper tier, cutting long-haul trips for cache misses in distant regions. No delete at Cloudflare. |
| `argoSmartRouting` | bool | not managed | Route origin-bound traffic over Cloudflare's measured fastest paths. Paid: monthly base fee plus per-GB usage. No delete at Cloudflare -- destroying with this true keeps the feature on and keeps billing. Apply `false` first when retiring. |
| `cacheVariants` | object | not managed | Per-extension lists of MIME types Cloudflare may serve as variants of cached assets. Keys are original-asset extensions (`avif`, `bmp`, `gif`, `jp2`, `jpeg`, `jpg`, `jpg2`, `png`, `tif`, `tiff`, `webp`); values are lists like `image/webp`. Only extensions with at least one value are sent. Real delete: destroy resets variants to defaults. |

## Examples

### Tiered Caching for Production

The free, recommended posture: Smart Tiered Cache picks the best upper tier, Regional Tiered Cache shortens long-haul misses.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareCacheSettings
metadata:
  name: prod-cache
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: prod-zone
      fieldPath: status.outputs.zone_id
  smartTieredCache: true
  regionalTieredCache: true
```

### Image Variants for a Media Zone

Serve WebP and AVIF variants of cached JPEG and PNG assets, with Smart Tiered Cache on.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareCacheSettings
metadata:
  name: media-cache
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  smartTieredCache: true
  cacheVariants:
    jpg:
      - image/webp
      - image/avif
    jpeg:
      - image/webp
    png:
      - image/webp
```

### Full Performance Posture

Everything on: tiered caching, Cache Reserve for long-tail assets, and Argo Smart Routing for origin-bound traffic. Both paid toggles are managed here, so retiring this zone means applying them as `false` before any destroy.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareCacheSettings
metadata:
  name: global-app-cache
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: global-app-zone
      fieldPath: status.outputs.zone_id
  smartTieredCache: true
  regionalTieredCache: true
  cacheReserve: true
  argoSmartRouting: true
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `zone_id` | string | The zone the cache settings belong to. Cache settings are a zone-scoped singleton with no resource ID of their own -- the zone is the identity. |

## Related Components

- **CloudflareDnsZone** -- the zone whose cache posture this resource manages, referenced via `zoneId`
- **CloudflareZoneSettings** -- general zone settings (SSL mode, minification, security level)
- **CloudflareZoneTlsSettings** -- the zone's TLS-specific posture
- **CloudflareRuleset** -- per-URL cache rules (TTLs, cache keys, bypass) via the cache settings phase
