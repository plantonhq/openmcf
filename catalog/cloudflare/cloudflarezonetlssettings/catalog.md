# Cloudflare Zone TLS Settings

Manages a Cloudflare zone's edge TLS posture as one resource: Universal SSL issuance, Total TLS per-hostname certificates, automatic origin TLS key exchange, origin TLS compliance modes, per-hostname TLS overrides, and certificate-authority hostname associations. Any field left unset is not managed -- the zone keeps its current value. The one hazard to know upfront: `universalSslEnabled: false` stops certificate issuance and can make proxied hostnames unreachable over HTTPS unless other certificates cover them.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Universal SSL setting** -- created only when `universalSslEnabled` is set; controls Universal SSL certificate issuance for the zone (no delete at Cloudflare -- destroy abandons the last-applied value)
- **Total TLS configuration** -- created only when `totalTls` is set; issues individual certificates for every proxied hostname, including deep subdomains Universal SSL's wildcard does not cover (no delete at Cloudflare)
- **Automatic origin TLS key exchange setting** -- created only when `autoOriginTlsKex` is set; lets Cloudflare negotiate the strongest key exchange the origin supports (no delete at Cloudflare)
- **Origin TLS compliance modes** -- created only when `originTlsComplianceModes` is non-empty; requires the listed modes on Cloudflare-to-origin connections (real delete -- destroy clears the requirement)
- **Per-hostname TLS settings** -- one API object per set attribute per hostname row: a `minTlsVersion` override, an `http2` override, and a `ciphers` override each become their own object keyed by hostname (real delete -- destroy removes the overrides and hostnames fall back to zone-wide settings)
- **Certificate-authority hostname associations** -- one object per `caHostnameAssociations` row: the row without `mtlsCertificateId` manages the zone's managed-CA list, and each row with one manages that mTLS certificate's hostname list (no delete at Cloudflare)

## Prerequisites

- A Cloudflare API token with SSL and Certificates edit access on the target zone
- An existing zone, either as a literal zone ID or a `CloudflareDnsZone` resource to reference via `valueFrom`
- Total TLS may require the zone's Advanced Certificate Manager subscription
- An mTLS certificate (literal ID or a `CloudflareMtlsCertificate` resource) if managing per-certificate CA hostname associations

## Quick Start

At least one TLS setting must be configured -- a resource that manages nothing would deploy nothing. This minimal manifest enables Total TLS and touches nothing else:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZoneTlsSettings
metadata:
  name: zone-tls
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  totalTls:
    enabled: true
```

```shell
planton apply -f zone-tls.yaml
```

Every proxied hostname in the zone gets its own certificate. Universal SSL, origin settings, and per-hostname overrides remain untouched.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone whose TLS settings are managed. Can reference a `CloudflareDnsZone` resource via `valueFrom` (defaults to `status.outputs.zone_id`). | Required |

At the spec level, at least one of the six optional surfaces below must be configured.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `universalSslEnabled` | bool | unset (not managed) | Universal SSL certificate issuance. Setting false STOPS issuance and can make proxied hostnames unreachable over HTTPS unless other certificates cover them. No delete at Cloudflare. |
| `totalTls.enabled` | bool | unset (not managed) | Whether Total TLS issues per-hostname certificates for the zone. May require the zone's Advanced Certificate Manager subscription. No delete at Cloudflare. |
| `totalTls.certificateAuthority` | string | Cloudflare chooses | The CA issuing Total TLS certificates. One of `google`, `lets_encrypt`, `ssl_com`. Certificate validity is fixed by Cloudflare at 90 days and is not configurable. |
| `autoOriginTlsKex` | bool | unset (not managed) | Automatic origin TLS key exchange on origin-facing TLS. No delete at Cloudflare. |
| `originTlsComplianceModes` | list(string) | unset (not managed) | Compliance modes required on Cloudflare-to-origin connections. Cloudflare documents `fips` and `pqh` (post-quantum hybrid); unknown strings pass through unvalidated as an open vocabulary. Real delete. |
| `hostnameSettings[].hostname` | string | -- | The hostname the overrides apply to. Required per row. Changing it replaces the override objects (the hostname is part of their API identity). |
| `hostnameSettings[].minTlsVersion` | string | unset (not managed) | Minimum TLS version for this hostname, overriding the zone-wide setting. One of `"1.0"`, `"1.1"`, `"1.2"`, `"1.3"`. |
| `hostnameSettings[].http2` | bool | unset (not managed) | HTTP/2 support for this hostname, overriding the zone-wide setting. |
| `hostnameSettings[].ciphers` | list(string) | unset (not managed) | TLS cipher allowlist for this hostname in BoringSSL format. Each row must set at least one of the three overrides. |
| `caHostnameAssociations[].hostnames` | list(string) | -- | The hostnames the certificate authority may issue for. At least one per row. |
| `caHostnameAssociations[].mtlsCertificateId` | StringValueOrRef | unset (managed-CA list) | The mTLS certificate whose hostname associations this row manages. Leave unset to manage the zone's managed-CA list. Can reference a `CloudflareMtlsCertificate` resource via `valueFrom`. |

## Examples

### Total TLS with a Chosen CA

Enables Total TLS with Google Trust Services and automatic origin key exchange. Everything else in the zone's TLS posture stays as it is.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZoneTlsSettings
metadata:
  name: zone-tls
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  totalTls:
    enabled: true
    certificateAuthority: google
  autoOriginTlsKex: true
```

### Per-Hostname Overrides with a Zone Reference

Raises the API hostname to TLS 1.3 while a legacy hostname stays on TLS 1.0 with a pinned cipher. The zone is referenced from a `CloudflareDnsZone` resource so the dependency stays in the graph.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZoneTlsSettings
metadata:
  name: hostname-tls
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: acme-zone
      fieldPath: status.outputs.zone_id
  hostnameSettings:
    - hostname: api.example.com
      minTlsVersion: "1.3"
      http2: true
    - hostname: legacy.example.com
      minTlsVersion: "1.0"
      ciphers:
        - "ECDHE-RSA-AES128-GCM-SHA256"
```

### Full TLS Posture

Manages all six surfaces: Total TLS, origin key exchange, FIPS compliance to the origin, a hostname override, and CA hostname associations for an mTLS certificate.

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZoneTlsSettings
metadata:
  name: full-tls-posture
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  universalSslEnabled: true
  totalTls:
    enabled: true
    certificateAuthority: google
  autoOriginTlsKex: true
  originTlsComplianceModes:
    - fips
  hostnameSettings:
    - hostname: api.example.com
      minTlsVersion: "1.3"
      http2: true
  caHostnameAssociations:
    - hostnames:
        - mtls.example.com
      mtlsCertificateId:
        value: "8d773bb3-90b2-4bd3-9d33-9c3ba33d5b9a"
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `zone_id` | string | The zone ID the TLS settings belong to. TLS settings are a zone-scoped singleton with no resource ID of their own -- the zone is the identity. |

## Related Components

- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- the zone whose TLS posture this resource manages, referenced by `zoneId`
- [Cloudflare Certificate Pack](/docs/catalog/cloudflare/cloudflarecertificatepack) -- advanced edge certificates, often the coverage that makes disabling Universal SSL safe
- [Cloudflare Custom Hostname](/docs/catalog/cloudflare/cloudflarecustomhostname) -- TLS for SaaS vanity hostnames outside this zone
- [Cloudflare Zone Settings](/docs/catalog/cloudflare/cloudflarezonesettings) -- zone-wide settings like minimum TLS version and Always Use HTTPS
