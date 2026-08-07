# DNS Record on Cloudflare

Deploys a DNS record within a Cloudflare zone, supporting A, AAAA, CNAME, MX, TXT, SRV, NS, and CAA record types with configurable proxy status, TTL, and priority. Integrates with Planton's Provider Connections for Cloudflare credential management and supports ValueFromRef wiring to resolve zone IDs from CloudflareDnsZone resources.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cloudflare DNS Record** -- a single DNS record in the specified zone with the configured type, value, TTL, and proxy setting
- **Proxy Configuration** -- created only when `proxied` is `true` on A, AAAA, or CNAME records; routes traffic through Cloudflare's CDN and WAF, hiding the origin IP
- **MX/SRV Priority** -- created only when the record type is MX or SRV; sets routing priority for mail exchange or service discovery

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **A Cloudflare DNS zone** -- the target zone must exist and be active. Provide the zone ID directly or reference a CloudflareDnsZone Cloud Resource via ValueFromRef.
- **Origin server or target** -- for A records, the origin server's IPv4 address; for AAAA, the IPv6 address; for CNAME, the target hostname; for MX, the mail server hostname.

## Deploy

### Console

Open the deployment store, find **DNS Record on Cloudflare**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Proxied A Record** preset in the [Presets](#presets) tab to pre-populate a proxied web-facing record.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareDnsRecord
metadata:
  name: www-record
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  name: www
  type: A
  value: "203.0.113.50"
  proxied: true
  ttl: 1
```

```shell
planton apply -f cloudflare-dns-record.yaml
```

This creates a proxied A record pointing `www` to the specified IP address with automatic TTL. Traffic flows through Cloudflare's CDN and WAF. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the DNS record to a zone deployed in the same InfraPipeline:

```yaml
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_id
```

The InfraPipeline resolves the dependency graph, deploys the DNS zone first, then provisions the DNS record with the resolved zone ID.

## Key Configuration

These are the most important decisions when configuring a DNS record. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Record type** -- Choose the DNS record type based on the use case: `A` for IPv4, `AAAA` for IPv6, `CNAME` for hostname aliases, `MX` for mail delivery, `TXT` for SPF/DKIM/verification, `SRV` for service discovery, `NS` for delegation, and `CAA` for certificate authority restrictions.

**Proxy vs. DNS-only** -- Set `proxied: true` (orange cloud) to route traffic through Cloudflare's CDN, WAF, and DDoS protection while hiding the origin IP. Set `proxied: false` (grey cloud) for DNS-only resolution that exposes the origin. Only A, AAAA, and CNAME records support proxying.

**TTL** -- Set `ttl: 1` for automatic TTL (recommended for proxied records). For DNS-only records, specify a value between 60 and 86400 seconds based on how frequently the target changes. Lower TTL enables faster failover but increases DNS query volume.

**Priority** -- Required for MX records, optional for SRV records. Lower values indicate higher priority. Use `priority: 10` for the primary mail server, with higher values (20, 30) for backups.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CloudflareDnsZone** | `zoneId` | `status.outputs.zone_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `record_id` | The unique identifier of the created DNS record in Cloudflare | Record management, update or delete operations |
| `hostname` | The fully qualified hostname of the DNS record | Application configuration, health check endpoints |
| `record_type` | The DNS record type that was created | Diagnostic and audit information |
| `proxied` | Whether the record is proxied through Cloudflare | Monitoring and configuration verification |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Proxied A record** -- An A record with Cloudflare proxy enabled for web servers and APIs. Traffic flows through Cloudflare's CDN and DDoS protection while hiding the origin IP. Start from the **Proxied A Record** preset.

**MX record for email** -- An MX record for mail delivery with priority routing. MX records cannot be proxied and are DNS-only. Configure for Google Workspace, Microsoft 365, or custom mail servers. Start from the **MX Record for Email** preset.

## Works With

- [**DNS Zone on Cloudflare**](/cloud-catalog/cloudflare-dns-zone) -- provides the zone ID where this DNS record is created