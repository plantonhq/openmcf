# DNS Record on Civo

Creates individual DNS records within a Civo DNS zone, supporting A, AAAA, CNAME, MX, TXT, SRV, and NS record types with configurable TTL and priority. Integrates with Planton's Provider Connections for Civo credential management and ValueFromRef for zone dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo DNS Domain Record** -- a single DNS record of the specified type within the referenced DNS zone, with the configured name, value, TTL, and priority

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **An existing Civo DNS zone** for the target domain. Provide the zone ID directly or reference a CivoDnsZone Cloud Resource via ValueFromRef.
- **The target value** for the record -- an IPv4 address for A records, a hostname for CNAME records, or the appropriate format for other record types.

## Deploy

### Console

Open the deployment store, find **DNS Record on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **A Record** preset in the [Presets](#presets) tab for a standard IPv4 address record.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoDnsRecord
metadata:
  name: api-record
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "abc12345-6789-def0-1234-567890abcdef"
  name: api
  type: A
  value: "192.0.2.1"
  ttl: 3600
```

```shell
planton apply -f civo-dns-record.yaml
```

This creates an A record pointing `api.{domain}` to the specified IPv4 address with a 1-hour TTL. No MX priority or additional record types are configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the record to a DNS zone deployed in the same InfraPipeline:

```yaml
spec:
  zoneId:
    valueFrom:
      kind: CivoDnsZone
      name: app-zone
      fieldPath: status.outputs.zone_id
```

The InfraPipeline resolves the dependency graph, deploys the DNS zone first, then provisions the DNS record within it.

## Key Configuration

These are the most important decisions when configuring a Civo DNS record. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Record type** -- The `type` field determines the DNS record type. Use `A` for IPv4 addresses, `AAAA` for IPv6, `CNAME` for hostname aliases, `MX` for mail routing, `TXT` for verification and SPF records, `SRV` for service discovery, and `NS` for delegating subdomains.

**Record name** -- The `name` field specifies the hostname relative to the zone. Use `"@"` for apex (root domain) records, or a subdomain like `"api"` or `"www"`. The resulting FQDN is `{name}.{zone-domain}`.

**TTL** -- The `ttl` field controls how long DNS resolvers cache the record, in seconds (60--86400). Defaults to 3600 (1 hour) when not specified. Use shorter TTLs (300s) for records that may change during migrations, and longer TTLs (86400s) for stable records to reduce DNS query load.

**Priority** -- Required for MX records, optional for SRV records. Lower values indicate higher priority. Set `priority` to configure mail server preference for MX records (e.g., 10 for primary, 20 for backup).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CivoDnsZone** | `zoneId` | `status.outputs.zone_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `record_id` | Unique identifier of the DNS record in Civo | Civo API operations, lifecycle management |
| `hostname` | Fully qualified hostname of the DNS record | Application configuration, health checks |
| `record_type` | DNS record type that was created | Monitoring, audit logs |
| `account_id` | Civo account ID where the record was created | Multi-account resource tracking |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**A record for a subdomain** -- points a subdomain (e.g., `api`) to an IPv4 address with a standard 1-hour TTL. The most common DNS record type for mapping hostnames to servers. Start from the **A Record** preset.

**CNAME alias** -- aliases a subdomain to another hostname, typically a load balancer or CDN endpoint. The target's IP is resolved at query time, so changes to the target propagate automatically. Start from the **CNAME Record** preset.

## Works With

- [**Civo DNS Zone**](/cloud-catalog/civo-dns-zone) -- provides the DNS zone in which this record is created