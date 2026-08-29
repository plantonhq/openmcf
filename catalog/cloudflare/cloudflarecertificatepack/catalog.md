# Certificate Pack on Cloudflare

Orders an advanced edge certificate for a zone: a publicly-trusted (browser-trusted) TLS certificate, provisioned and auto-renewed by Cloudflare, that covers the hostnames you list beyond the free Universal SSL certificate. Use it when you need a specific certificate authority, multiple or longer-lived certificates, or coverage for hostnames Universal SSL does not include. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Certificate Pack** -- an advanced, auto-renewed edge certificate for the zone
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has SSL and Certificates edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A zone** -- the certificate pack is ordered against an existing `CloudflareDnsZone`.

## Deploy

### Console

Open the deployment store, find **Certificate Pack on Cloudflare**, and click **Deploy**. The creation wizard captures the zone, the certificate authority, validation method, validity, and the covered hosts.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareCertificatePack
metadata:
  name: edge-cert-pack
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: acme-zone
      fieldPath: status.outputs.zone_id
  certificateAuthority: google
  validationMethod: txt
  validityDays: 90
  hosts:
    - example.com
    - "*.example.com"
```

```shell
planton apply -f cloudflare-certificate-pack.yaml
```

This orders a 90-day advanced certificate from Google Trust Services covering the apex and its subdomains, validated automatically over TXT. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a certificate pack. Explore the full field reference in the [API Explorer](#api-explorer) tab. Almost every attribute is immutable -- changing it re-orders (replaces) the pack.

**Zone (`zoneId`)** -- The zone the certificate is ordered for. Reference a `CloudflareDnsZone` to keep the dependency in the graph.

**Certificate Authority (`certificateAuthority`)** -- `google`, `lets_encrypt`, or `ssl_com`. CA-specific restrictions apply.

**Validation Method (`validationMethod`)** -- `txt`, `http`, or `email`. TXT validates automatically for zones on Cloudflare nameservers.

**Validity (`validityDays`)** -- 14, 30, 90, or 365 days. Cloudflare auto-renews the pack before expiry.

**Hosts (`hosts`)** -- The hostnames the certificate covers. Must include the zone apex and may not exceed 50 hosts.

## Outputs and Dependencies

### What This Component Consumes

The pack references a **CloudflareDnsZone** (via `zoneId`).

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `certificate_pack_id` | The Cloudflare-assigned pack identifier | Verification, dashboards |
| `primary_certificate` | Identifier of the primary certificate in the pack | Auditing |
| `zone_id` | The zone the pack was ordered in | Verification, imports, composition |

Issuance status is deliberately not an output — it transitions asynchronously (`initializing` → `pending_validation` → `active`); read it from the Cloudflare API or dashboard.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Apex + wildcard** -- cover `example.com` and `*.example.com` with a single advanced pack.

**Specific CA** -- order from a particular certificate authority when compliance requires it.

## Works With

- [**DNS Zone on Cloudflare**](/cloud-catalog/cloudflare-dns-zone) -- the zone the certificate is ordered for
