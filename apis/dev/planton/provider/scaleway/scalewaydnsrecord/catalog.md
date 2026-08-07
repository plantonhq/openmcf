# Scaleway DNS Record

Deploys a standalone DNS record on Scaleway within an existing DNS zone. Designed for records whose values come from other infrastructure resources -- A records pointing to a Load Balancer IP, CNAMEs to a Kapsule cluster endpoint -- creating visible dependency edges in InfraChart DAGs. Supports all Scaleway DNS record types (A, AAAA, ALIAS, CAA, CNAME, DNAME, MX, NS, PTR, SOA, SRV, TXT, TLSA) with ValueFromRef for both zone reference and record data wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS Record** -- a single `scaleway_domain_record` resource in the referenced DNS zone with the configured name, type, data, TTL, and optional priority

Note: Scaleway DNS records do not support tags. Unlike most other Scaleway resources, the DNS API does not accept tags or labels.

## Before You Deploy

### Scaleway Account

- **A Scaleway account** with an active project and API access key pair.
- **An existing DNS zone** on Scaleway. Provide the zone name directly or reference a ScalewayDnsZone Cloud Resource via ValueFromRef.
- **The record data value** -- an IP address, hostname, or other record-type-specific data. For cross-resource wiring, use ValueFromRef to reference another resource's output (e.g., a Load Balancer's `lb_ip_address` or an Instance's `public_ip_address`).

## Deploy

### Console

Open the deployment store, find **Scaleway DNS Record**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **A Record** preset in the [Presets](#presets) tab to create a hostname-to-IP mapping.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayDnsRecord
metadata:
  name: app-dns
  org: acme-corp
  env: prod
spec:
  zoneName:
    value: "example.com"
  name: app
  type: A
  data:
    value: "51.159.26.100"
  ttl: 3600
```

```shell
planton apply -f scaleway-dns-record.yaml
```

This creates an A record `app.example.com` pointing to the specified IP address with a 1-hour TTL. No cross-resource wiring is configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the record to a DNS zone and a Load Balancer deployed in the same InfraPipeline:

```yaml
spec:
  zoneName:
    valueFrom:
      kind: ScalewayDnsZone
      name: my-zone
      fieldPath: status.outputs.zone_name
  data:
    valueFrom:
      kind: ScalewayLoadBalancer
      name: web-lb
      fieldPath: status.outputs.lb_ip_address
```

The InfraPipeline resolves the dependency graph, deploys the DNS zone and Load Balancer first, then creates the DNS record with the resolved values.

## Key Configuration

These are the most important decisions when configuring a DNS record. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Record type** -- The `type` field selects the DNS record type. Use `A` for IPv4 addresses, `CNAME` for hostname aliases (cannot be used at zone apex), `ALIAS` for Scaleway-native zone apex aliases, `MX` for mail routing, and `TXT` for SPF, DKIM, and domain verification. Record type cannot be changed after creation.

**Record data** -- The `data` field accepts a literal value or a ValueFromRef reference. Data format depends on the record type -- e.g., `"192.0.2.1"` for A records, `"target.example.com."` (trailing dot) for CNAME records. Use ValueFromRef to reference outputs from Load Balancers, Instances, or Kapsule clusters.

**TTL** -- The `ttl` field controls how long resolvers cache the record. Use 300 (5 minutes) during migrations, 3600 (1 hour, default) for most records, or 86400 (24 hours) for static records that rarely change.

**Priority** -- The `priority` field is meaningful only for MX and SRV records. Lower values indicate higher priority. Ignored for all other record types.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **ScalewayDnsZone** | `zoneName` | `status.outputs.zone_name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `record_id` | Unique identifier of the created DNS record | Scaleway API operations, Terraform import |
| `fqdn` | Fully qualified domain name of the record | Downstream resource configuration, user reference without manual string concatenation |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**A record** -- Maps a hostname to an IPv4 address. The most common record type for pointing subdomains at Instances, Load Balancers, or Public Gateways. Start from the **A Record** preset.

**CNAME record** -- Creates a hostname alias pointing to another hostname. Used for `www` aliases, CDN endpoints, and Kapsule cluster wildcard DNS targets. Cannot be used at the zone apex -- use ALIAS type instead. Start from the **CNAME Record** preset.

## Works With

- [**Scaleway DNS Zone**](/cloud-catalog/scaleway-dns-zone) -- provides the DNS zone that this record is created in