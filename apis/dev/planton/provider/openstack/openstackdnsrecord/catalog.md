# OpenStack DNS Record

Deploys a standalone Designate DNS recordset on OpenStack, mapping a fully qualified domain name to one or more values within an existing DNS zone. The record references an OpenStackDnsZone via ValueFromRef, making it a DAG-visible, independently managed alternative to inline records defined directly in the zone spec.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Designate DNS Recordset** -- a single recordset (one name + one type + one or more values) in the specified DNS zone, supporting A, AAAA, CNAME, MX, TXT, SRV, NS, PTR, CAA, SOA, SPF, SSHFP, and NAPTR record types
- **Round-Robin Values** -- created only when multiple values are provided; Designate returns all values to resolvers for client-side load distribution

## Before You Deploy

### OpenStack Account

- **DNS zone** -- an existing Designate zone where the record will be created. Provide the zone ID directly or reference an OpenStackDnsZone Cloud Resource via ValueFromRef.
- **Designate service** -- the Designate DNS service must be enabled in your OpenStack deployment. Run `openstack zone list` to verify availability.
- **FQDN format** -- record names must be fully qualified with a trailing dot (e.g., `app.example.com.`). Designate rejects names without the trailing dot.

## Deploy

### Console

Open the deployment store, find **OpenStack DNS Record**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **A Record** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackDnsRecord
metadata:
  name: app-a-record
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "<zone-id>"
  recordName: "app.example.com."
  type: A
  values:
    - "203.0.113.42"
```

```shell
planton apply -f dns-record.yaml
```

This creates an A record pointing `app.example.com.` to a single IPv4 address with the zone's default TTL. No custom TTL or description is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the record to a zone deployed in the same InfraPipeline:

```yaml
spec:
  zoneId:
    valueFrom:
      kind: OpenStackDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_id
```

The InfraPipeline resolves the dependency graph, deploys the DNS zone first, then provisions the record with the resolved zone ID.

## Key Configuration

These are the most important decisions when configuring a DNS record. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Record type** -- Determines the kind of DNS mapping. A and AAAA records map to IP addresses, CNAME records alias to another hostname, MX records direct mail delivery, and TXT records store arbitrary text (SPF, DKIM, verification tokens). The type is immutable after creation.

**TTL** -- Controls how long resolvers cache this record. If omitted, the zone's default TTL applies. Use lower values (60-300 seconds) for records that change frequently (blue-green deployments, failover) and higher values (3600+ seconds) for stable records.

**Multiple values** -- Providing multiple entries in the `values` field creates a round-robin recordset. For A records, this distributes traffic across multiple IP addresses at the DNS level.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackDnsZone** | `zoneId` | `status.outputs.zone_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `recordset_id` | UUID of the DNS recordset | Import, audit, and debugging |
| `fqdn` | Fully qualified domain name of the record | Downstream service endpoint references |
| `record_type` | DNS record type that was created | Monitoring labels, automation |
| `values` | List of DNS record values | Verification, downstream configuration |
| `zone_id` | ID of the zone containing this record | Zone-aware downstream resources |
| `region` | OpenStack region where the record was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**A record** -- Maps a hostname to one or more IPv4 addresses with a 5-minute TTL. The standard pattern for pointing hostnames at instance floating IPs or load balancer VIPs. Start from the **A Record** preset.

**CNAME record** -- Aliases one hostname to another. Use for vanity names, service aliases, or pointing to external service endpoints. Both the record name and target value must be FQDNs with trailing dots. Start from the **CNAME Record** preset.

## Works With

- [**OpenStack DNS Zone**](/cloud-catalog/openstack-dns-zone) -- provides the zone ID where the DNS record is created