---
title: "DNS Zone"
description: "DNS Zone deployment documentation"
icon: "package"
order: 100
componentName: "openstackdnszone"
---

# OpenStack DNS Zone

Deploys a Designate DNS zone on OpenStack -- the authoritative container for DNS records under a domain name. The zone supports primary and secondary modes, optional inline records, and configurable TTL and SOA email. For independently managed records with DAG visibility in InfraCharts, use the separate OpenStackDnsRecord component.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Designate DNS Zone** -- an authoritative zone for the specified domain with configurable type (PRIMARY or SECONDARY), TTL, SOA email, and optional master nameservers for zone transfers
- **Inline DNS Recordsets** -- created only when `records` entries are specified; one recordset per entry, each provisioned as a separate resource keyed by record type + record name for stable IaC state

## Before You Deploy

### OpenStack Account

- **Designate service** -- the Designate DNS service must be enabled in your OpenStack deployment. Run `openstack zone list` to verify availability.
- **Domain ownership** -- for production zones, ensure the domain's registrar delegates to the Designate nameservers. Without proper delegation, records are resolvable only within the OpenStack deployment.
- **Secondary zone masters** -- if creating a SECONDARY zone, have the master nameserver addresses ready. Designate will perform zone transfers (AXFR/IXFR) from these servers.

## Deploy

### Console

Open the deployment store, find **OpenStack DNS Zone**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Primary Zone** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackDnsZone
metadata:
  name: example-zone
  org: acme-corp
  env: prod
spec:
  domainName: "example.com."
  email: "admin@example.com"
  ttl: 3600
```

```shell
planton apply -f dns-zone.yaml
```

This creates a primary DNS zone for `example.com.` with a 1-hour default TTL and an SOA admin email. No inline records, secondary masters, or region override are configured.

## Key Configuration

These are the most important decisions when configuring a DNS zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Zone type** -- PRIMARY zones are authoritative (Designate manages all records). SECONDARY zones replicate from upstream master nameservers via zone transfer. If omitted, Designate defaults to PRIMARY. The type is immutable after creation.

**Default TTL** -- The `ttl` field sets the zone-wide default TTL applied to records that do not specify their own. Use 3600 (1 hour) for stable production zones and lower values (60-300) for zones with frequently changing records.

**Inline vs standalone records** -- The `records` field creates recordsets alongside the zone in one manifest. For InfraChart deployments where individual records need DAG-level dependency tracking (e.g., an A record pointing to a floating IP resolved via ValueFromRef), use the standalone OpenStackDnsRecord component instead.

**SOA email** -- The `email` field sets the administrative contact in the zone's SOA record. While optional, production zones should include a valid contact email for operational transparency.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zone_id` | UUID of the DNS zone | DNS record creation via `zoneId` |
| `zone_name` | DNS zone name (the authoritative domain) | Monitoring labels, resource identification |
| `region` | OpenStack region where the zone was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Primary zone** -- Creates an authoritative DNS zone with a 1-hour default TTL and SOA email. The starting point for hosting DNS in Designate -- add records inline or as standalone OpenStackDnsRecord resources. Start from the **Primary Zone** preset.

## Works With

This component operates independently and does not reference other deployment components.