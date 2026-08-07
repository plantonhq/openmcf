---
title: "DNS Record"
description: "DNS Record deployment documentation"
icon: "package"
order: 100
componentName: "digitaloceandnsrecord"
---

# DNS Record on DigitalOcean

Creates a single DNS record within an existing DigitalOcean DNS zone. Supports A, AAAA, CNAME, MX, TXT, SRV, NS, and CAA record types, with type-specific fields for priority, weight, port, flags, and tag applied conditionally. Integrates with Planton's Provider Connections for DigitalOcean credential management and ValueFromRef for DNS zone and target dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS Record** -- a `digitalocean_record` resource in the specified domain with the configured type, name, value, and TTL
- **Type-Specific Attributes** -- `priority` is set for MX and SRV records; `weight` and `port` are set for SRV records; `flags` and `tag` are set for CAA records; all are omitted for inapplicable record types

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **An existing DigitalOcean DNS zone (domain)** managed by DigitalOcean's DNS service. Provide the domain name directly or reference a DigitalOceanDnsZone Cloud Resource via ValueFromRef.
- **A valid record value** matching the record type: an IPv4 address for A records, a hostname for CNAME records, a mail server for MX records, etc.

## Deploy

### Console

Open the deployment store, find **DNS Record on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **A Record** preset in the [Presets](#presets) tab to point a hostname to an IP address.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1
kind: DigitalOceanDnsRecord
metadata:
  name: www-a-record
  org: acme-corp
  env: prod
spec:
  domain:
    value: "example.com"
  name: www
  type: A
  value:
    value: "192.0.2.1"
  ttlSeconds: 3600
```

```shell
planton apply -f dns-record.yaml
```

This creates an A record pointing `www.example.com` to `192.0.2.1` with a one-hour TTL. No MX, SRV, or CAA-specific fields are configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the record to a DNS zone deployed in the same InfraPipeline:

```yaml
spec:
  domain:
    valueFrom:
      kind: DigitalOceanDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_name
```

The InfraPipeline resolves the dependency graph, deploys the DNS zone first, then provisions the DNS record within it.

## Key Configuration

These are the most important decisions when configuring a DNS record. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Record type** -- The `type` field determines the DNS record type. A and AAAA resolve to IP addresses. CNAME creates an alias to another hostname (cannot be used on the root domain `@` on DigitalOcean). MX routes mail to a specified server. TXT stores arbitrary text (SPF, DKIM, domain verification). SRV and CAA have additional required fields.

**TTL** -- The `ttlSeconds` field controls how long DNS resolvers cache this record, defaulting to 1800 seconds (30 minutes). Use lower values (60-300) during migrations or when records change frequently, and higher values (3600-86400) for stable production records.

**Type-specific fields** -- MX records require `priority` (lower values = higher priority). SRV records require `priority`, `weight`, and `port`. CAA records require `flags` and `tag` (`issue`, `issuewild`, or `iodef`). The protobuf schema enforces these cross-field constraints at validation time.

**Value references** -- The `value` field supports ValueFromRef, allowing you to reference outputs from other Cloud Resources (e.g., a Droplet's IP address or a Load Balancer's hostname) instead of hardcoding values.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanDnsZone** | `domain` | `status.outputs.zone_name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `record_id` | Unique identifier of the DNS record in DigitalOcean | API operations, record management |
| `hostname` | Fully qualified hostname (e.g., `www.example.com`) | Application configuration, health check targets |
| `record_type` | DNS record type that was created (A, CNAME, etc.) | Audit logs, inventory tracking |
| `domain` | Domain name (DNS zone) where the record was created | Cross-referencing with zone management |
| `ttl_seconds` | TTL in seconds applied to the record | Cache behavior verification |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**A record** -- Points a hostname to an IPv4 address with a one-hour TTL. Use for root domain or subdomain records targeting Droplets, Load Balancers, or external IP addresses. Start from the **A Record** preset.

**CNAME record** -- Aliases a subdomain to another hostname. Common for `www` pointing to the root domain, subdomains pointing to CDN origins, or third-party service integrations. Start from the **CNAME Record** preset.

## Works With

- [**DNS Zone on DigitalOcean**](/cloud-catalog/digital-ocean-dns-zone) -- provides the domain (DNS zone) in which records are created