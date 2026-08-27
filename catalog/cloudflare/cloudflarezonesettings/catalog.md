# Cloudflare Zone Settings

Manages a Cloudflare zone's behavior settings -- the serving posture the dashboard shows under the Speed, Security, and Network tabs -- as one typed field per Cloudflare setting id, plus managed header transforms, URL normalization, origin cloud-region hints, and the waiting-room crawler bypass. Only the fields you set are managed: unset fields are never sent, and because Cloudflare has no delete for zone settings, destroying the resource abandons the last-applied values rather than reverting them.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Zone Settings** -- one `cloudflare_zone_setting` per managed field, sent to the zone settings API in Cloudflare's own vocabulary (63 settings available: on/off toggles, enums, numerics, and object values like HSTS)
- **Managed Transforms** -- created only when `managedRequestHeaders` or `managedResponseHeaders` has entries; one zone-wide object toggling Cloudflare-defined header transforms
- **URL Normalization Settings** -- created only when `urlNormalization` is set; controls how Cloudflare rewrites URLs before rules match and before they reach the origin
- **Origin Cloud Regions** -- created only when `originCloudRegions` has entries; one region hint per origin IP so Cloudflare can optimize origin-facing routing
- **Waiting Room Settings** -- created only when `waitingRoomCrawlerBypass` is set; the zone-wide toggle letting verified search crawlers bypass waiting rooms

Zone settings and the waiting-room toggle have no delete at Cloudflare: destroy abandons their live values. Managed transforms and URL normalization do have a real delete and reset on destroy.

## Prerequisites

- **A Cloudflare zone** -- typically a CloudflareDnsZone resource whose `zone_id` output this resource references, or a literal zone ID from the dashboard
- **A Cloudflare API token** with Zone Settings Edit permission on the target zone (managed transforms, URL normalization, and origin cloud regions ride the same zone-scoped token)
- **The right zone plan for gated settings** -- `advancedDdos`, `orangeToOrange`, `prefetchPreload`, `responseBuffering`, `sortQueryStringForCache`, `trueClientIpHeader`, and `proxyReadTimeout` need Enterprise; `polish`, `mirage`, and `imageResizing` need Pro or above. The apply fails with the API's editable=false error on a plan that lacks a setting; nothing is billed or upgraded

## Quick Start

The minimal manifest manages exactly two settings and touches nothing else:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZoneSettings
metadata:
  name: baseline-settings
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  alwaysUseHttps: true
  minTlsVersion: "1.2"
```

```shell
planton apply -f zone-settings.yaml
```

This redirects all plain-http requests to HTTPS and refuses TLS below 1.2. Every other zone setting stays exactly as it was -- unset means unmanaged.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone whose settings are managed. Can reference a CloudflareDnsZone resource via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. At least one setting must also be configured -- a resource that manages nothing is rejected. |

### Optional Fields

All settings are optional and default to **not managed** -- the module never sends a field you did not set. The most-used settings:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `alwaysUseHttps` | bool | not managed | Answer every plain-http request with a 301 to the https equivalent |
| `ssl` | string | not managed | Encryption mode to the origin: `off`, `flexible`, `full`, `strict`. Production zones should use `strict` |
| `minTlsVersion` | string | not managed | Minimum TLS version accepted: `1.0`, `1.1`, `1.2`, `1.3`. Quote the value in YAML |
| `securityLevel` | string | not managed | Challenge aggressiveness: `off`, `essentially_off`, `low`, `medium`, `high`, `under_attack` |
| `browserCacheTtl` | int64 | not managed | Seconds browsers may cache resources. Cloudflare accepts a fixed value set (0 and durations from 30s to one year); 0 means respect existing headers |
| `cacheLevel` | string | not managed | `basic`, `simplified`, or `aggressive` (cache all query-string variants) |
| `securityHeader` | object | not managed | HSTS configuration: `enabled`, `includeSubdomains`, `maxAge`, `nosniff`, `preload` |
| `http3` | bool | not managed | HTTP/3 (QUIC) for client connections |
| `ipv6` | bool | not managed | IPv6 termination on all proxied hostnames |

The remaining settings group by shape; see the API reference for the full field list:

| Group | Fields | Description |
|-------|--------|-------------|
| On/off toggles | 39 booleans: `automaticHttpsRewrites`, `brotli`, `earlyHints`, `http2`, `websockets`, `zeroRtt` (setting id `0rtt`), `browserCheck`, `emailObfuscation`, `hotlinkProtection`, `ipGeolocation`, `rocketLoader`, `mirage`, `webp`, `developmentMode`, `sslRecommender`, `tlsClientAuth`, `waf`, `longLivedGrpc`, and the plan-gated toggles, among others | Sent to the API as `"on"`/`"off"` |
| Enum settings | `polish`, `tls13`, `h2Prioritization`, `imageResizing`, `pseudoIpv4`, `cnameFlattening`, `originMaxHttpVersion`, `transformations`, plus free-string `transformationsAllowedOrigins` | Validated against Cloudflare's accepted values at manifest validation time |
| Numeric settings | `challengeTtl`, `edgeCacheTtl`, `maxUpload`, `originH2MaxStreams`, `proxyReadTimeout` | Cloudflare validates each against a fixed accepted value set |
| Object settings | `nel`, `aegis`, `automaticPlatformOptimization`, `ciphers` (list) | APO requires every field of its object on writes; `aegis` needs an Aegis entitlement; `ciphers` writes need the zone's Advanced Certificate Manager subscription (API code 1023 without it) |
| Companions | `managedRequestHeaders`, `managedResponseHeaders`, `urlNormalization`, `originCloudRegions`, `waitingRoomCrawlerBypass` | Zone-wide serving configuration living outside the settings API |

## Examples

### Production Security Posture

HTTPS everywhere, strict origin encryption, modern TLS, and a one-year HSTS pin.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZoneSettings
metadata:
  name: prod-security
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  alwaysUseHttps: true
  ssl: strict
  minTlsVersion: "1.2"
  securityLevel: medium
  securityHeader:
    enabled: true
    includeSubdomains: true
    maxAge: 31536000
    nosniff: true
    preload: false
```

### Performance Tuning

Compression, early hints, HTTP/3, aggressive caching, and lossless image optimization (Polish needs Pro or above).

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZoneSettings
metadata:
  name: perf-tuning
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  brotli: true
  earlyHints: true
  http3: true
  cacheLevel: aggressive
  browserCacheTtl: 14400
  polish: lossless
```

### Referencing the Zone and Managing Companions

Wire `zoneId` from a CloudflareDnsZone in the same environment, toggle managed transforms, and pin URL normalization.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZoneSettings
metadata:
  name: zone-serving-config
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_id
  alwaysUseHttps: true
  managedRequestHeaders:
    - id: add_true_client_ip_headers
      enabled: true
  managedResponseHeaders:
    - id: remove_x-powered-by_header
      enabled: true
  urlNormalization:
    scope: incoming
    type: cloudflare
  originCloudRegions:
    - originIp: "203.0.113.10"
      region: us-east-1
      vendor: aws
```

### Full Zone Posture

Security, performance, TLS ciphers, and the waiting-room crawler bypass in one manifest. The `ciphers` list needs the zone's Advanced Certificate Manager subscription -- drop it on zones without ACM.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZoneSettings
metadata:
  name: full-zone-posture
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_id
  alwaysUseHttps: true
  ssl: strict
  minTlsVersion: "1.2"
  securityLevel: medium
  browserCacheTtl: 14400
  http3: true
  ipv6: true
  ciphers:
    - "ECDHE-ECDSA-AES128-GCM-SHA256"
    - "ECDHE-RSA-AES128-GCM-SHA256"
  securityHeader:
    enabled: true
    includeSubdomains: true
    maxAge: 31536000
    nosniff: true
    preload: false
  waitingRoomCrawlerBypass: false
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `zone_id` | string | The zone ID the settings belong to. Zone settings are a zone-scoped singleton with no resource id of their own -- the zone is the identity, passed through for downstream references. |

## Related Components

- **CloudflareDnsZone** -- creates the zone this resource configures and owns the DNS-level settings (SOA, NS TTL, DNSSEC) that are out of scope here; its `zone_id` output is this resource's foreign key
- **CloudflareCacheSettings** -- cache rules, tiered caching, and cache reserve for the same zone
- **CloudflareZoneTlsSettings** -- advanced per-zone TLS configuration beyond this resource's TLS toggles
- **CloudflareRuleset** -- per-request overrides of these zone-wide settings via the `http_config_settings` phase
