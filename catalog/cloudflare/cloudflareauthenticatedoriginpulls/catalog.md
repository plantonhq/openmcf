# Cloudflare Authenticated Origin Pulls

A zone's Authenticated Origin Pulls enablement: the zone-wide toggle making Cloudflare present a client certificate when pulling from your origin, plus per-hostname associations pinning hostnames to uploaded certificates. The origin completes the story by requiring and validating the certificate. Free-plan feature. Neither surface deletes on destroy: the toggle is abandoned, associations are reverted.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Zone toggle** -- one `cloudflare_authenticated_origin_pulls_settings` when `zoneEnabled` is set
- **Associations** -- one `cloudflare_authenticated_origin_pulls` per hostname row (the provider requires exactly one hostname per resource; the module fans rows out)

## Prerequisites

- **A Cloudflare zone** (free plans included)
- **An origin configured to require and validate Cloudflare's client certificate** -- without the origin-side check, enabling AOP protects nothing
- **Uploaded client certificates** (`CloudflareAuthenticatedOriginPullsCertificate`) when pinning per hostname
- **A Cloudflare API token** with Zone → SSL and Certificates → Edit

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareAuthenticatedOriginPulls
metadata:
  name: www-origin-pulls
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  zoneEnabled: true
```

```shell
planton apply -f origin-pulls.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone whose AOP surface is managed. Can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. |

At least one of `zoneEnabled` / `hostnameAssociations` must be set.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `zoneEnabled` | bool | unset (not managed) | The zone-wide toggle. Presence, not truthiness, decides whether the module manages it. |
| `hostnameAssociations` | list | `[]` | Rows of `{hostname, certificateId, enabled}`. `certificateId` references a CloudflareAuthenticatedOriginPullsCertificate; unset rides the zone-level material. `enabled` unset is sent as true (Cloudflare treats null as "void the association"). |

## Examples

### Zone toggle plus a pinned hostname

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareAuthenticatedOriginPulls
metadata:
  name: www-origin-pulls
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: acme-zone
      fieldPath: status.outputs.zone_id
  zoneEnabled: true
  hostnameAssociations:
    - hostname: app.example.com
      certificateId:
        valueFrom:
          kind: CloudflareAuthenticatedOriginPullsCertificate
          name: app-client-certificate
    - hostname: legacy.example.com
      enabled: false
```

## Destroy Semantics

Neither surface deletes. The zone toggle has NO delete at Cloudflare -- destroy abandons whatever value is live. An association is removed by a revert write (`enabled: null`) issued from state. Disable explicitly before destroying when the zone outlives this resource.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `zone_id` | string | The zone whose AOP surface is managed |

## Related Components

- [Cloudflare Authenticated Origin Pulls Certificate](/docs/catalog/cloudflare/cloudflareauthenticatedoriginpullscertificate) -- the uploaded client certificates rows reference
- [Cloudflare mTLS Certificate](/docs/catalog/cloudflare/cloudflaremtlscertificate) -- account-level CA material
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key
