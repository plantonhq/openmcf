# DNS Record on OCI

Deploys an Oracle Cloud Infrastructure DNS record set (RRSet) -- a set of DNS resource records sharing the same domain name and record type within an OCI DNS zone. Updates replace the entire record set atomically. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to DNS zones.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS Record Set** -- a `dns.Rrset` containing one or more DNS records of the specified type (A, AAAA, CNAME, MX, TXT, SRV, CAA, NS, PTR) for the specified domain within the target zone

This component does not create freeform tags -- OCI DNS record sets do not support tagging.

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- An OCI DNS zone (public or private) to place the records in. Provide the zone name or OCID directly, or reference an OciDnsZone Cloud Resource via ValueFromRef.
- For private DNS zones: the OCID of the private DNS view, provided via `viewId`.

## Deploy

### Console

Open the deployment store, find **DNS Record on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **A Record** preset in the [Presets](#presets) tab to pre-populate a simple IPv4 address record.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciDnsRecord
metadata:
  name: app-a-record
  org: acme-corp
  env: prod
spec:
  zoneNameOrId:
    value: "example.com"
  domain: app.example.com
  rtype: A
  items:
    - rdata: "203.0.113.10"
      ttl: 300
```

```shell
planton apply -f dns-record.yaml
```

This creates an A record pointing `app.example.com` to a single IPv4 address with a 5-minute TTL. Add multiple entries to `items` for round-robin DNS.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the record to a DNS zone deployed in the same InfraPipeline:

```yaml
spec:
  zoneNameOrId:
    valueFrom:
      kind: OciDnsZone
      name: main-zone
      fieldPath: status.outputs.zoneId
```

The InfraPipeline resolves the dependency graph, deploys the DNS zone first, then provisions the record set with the resolved zone identifier.

## Key Configuration

These are the most important decisions when configuring a DNS record. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Record type** -- Set `rtype` to the DNS record type (A, AAAA, CNAME, MX, TXT, SRV, CAA, NS, PTR). This is a ForceNew field -- changing the record type destroys and recreates the record set. CNAME records cannot coexist with other record types on the same domain.

**Domain** -- The `domain` field is the fully qualified domain name (e.g., `app.example.com`). This is a ForceNew field. The domain must be within the scope of the target zone.

**Record data format** -- Each item's `rdata` uses type-specific presentation format: IPv4 addresses for A records, FQDNs with trailing dot for CNAME and MX records, quoted strings for TXT records. OCI may normalize returned values (IPv6 compression, trailing-dot appended to hostnames, TXT quote stripping) -- the stored value may differ from input.

**Atomic replacement** -- Updates replace the entire record set, not individual records. To add a record to an existing set, include all existing records plus the new one in `items`.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciDnsZone** | `zoneNameOrId` | `status.outputs.zoneId` |

### What This Component Provides

This component does not produce stack outputs. The record set is identified by its (zone, domain, rtype) tuple, all of which are input fields.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**A record** -- A single IPv4 address record for directing traffic to a compute instance or load balancer. Start from the **A Record** preset.

**CNAME alias** -- A CNAME record aliasing one domain to another (e.g., `www.example.com` to `app.example.com`). Start from the **CNAME Alias** preset.

## Works With

- [**DNS Zone on OCI**](/cloud-catalog/oci-dns-zone) -- provides the DNS zone that contains this record set