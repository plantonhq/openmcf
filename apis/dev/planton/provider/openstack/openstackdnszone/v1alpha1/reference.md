# OpenStackDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackDnsZoneSpec defines the configuration for an OpenStack Designate DNS zone
with optional inline DNS records.

A DNS zone represents an authoritative domain (e.g., "example.com") managed by
the OpenStack Designate service. The zone can be PRIMARY (Designate is authoritative)
or SECONDARY (replicated from upstream master nameservers).

This component supports two modes for managing DNS records:
  - **Inline records** (via the `records` field): Convenient for self-contained
    zones where all records are defined in one place.
  - **Standalone records** (via the separate OpenStackDnsRecord component):
    For DAG-visible, independently managed records in InfraCharts.

Terraform resource: openstack_dns_zone_v2 + openstack_dns_recordset_v2 (for inline records)
Pulumi resource:    dns.Zone + dns.RecordSet (for inline records)

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackDnsZone
metadata:
  name: test-zone
spec:
  domain_name: "example.com"
  email: "admin@example.com"
  description: "Test DNS zone for local development"
  records:
    - record_type: 1
      record_name: "www.example.com."
      values:
        - "192.0.2.1"
      ttl: 300
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.domainName` | `string` | yes |  |  |
| `spec.email` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.ttl` | `int32` |  |  |  |
| `spec.type` | `string` |  |  |  |
| `spec.masters` | `[]string` |  |  |  |
| `spec.records` | `[]OpenStackDnsRecord` |  |  |  |
| `spec.records[].recordType` | `enum` | yes |  |  |
| `spec.records[].recordName` | `string` | yes |  |  |
| `spec.records[].values` | `[]string` | yes |  |  |
| `spec.records[].ttl` | `int32` |  | `60` |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.domainName

`string` · required

(Required) The DNS domain name for the zone (e.g., "example.com").
This is the authoritative domain managed by Designate.
Must be a valid domain name (lowercase labels separated by dots, ending with a TLD).
Note: Designate may auto-append a trailing dot internally.
ForceNew: changing this requires recreating the zone.

- rule: domain_name must be a valid DNS domain (e.g., example.com)
- rule: {"required":true}

### spec.email

`string`

(Optional) Email address of the zone administrator.
Used in the SOA record for the zone.
Example: "admin@example.com"

### spec.description

`string`

(Optional) Human-readable description of the DNS zone.

### spec.ttl

`int32` · optional (explicit presence)

(Optional) Default Time To Live (in seconds) for records in this zone.
Determines how long resolvers cache records from this zone.
If omitted, Designate uses its deployment default.

### spec.type

`string`

(Optional) The zone type.
"PRIMARY" for zones where Designate is the authoritative source.
"SECONDARY" for zones replicated from upstream master nameservers.
If omitted, Designate defaults to PRIMARY.
ForceNew: changing this requires recreating the zone.

### spec.masters

`[]string`

(Optional) List of master nameserver addresses for SECONDARY zones.
Required when type is "SECONDARY". Ignored for PRIMARY zones.
Example: ["ns1.upstream.com", "ns2.upstream.com"]

### spec.records

`[]OpenStackDnsRecord`

(Optional) Inline DNS records to create alongside the zone.
Each record is provisioned as a separate openstack_dns_recordset_v2 resource,
keyed by record_type + record_name for stable IaC state management.
For DAG-visible, independently managed records in InfraCharts, use the
standalone OpenStackDnsRecord component instead.

### spec.records[].recordType

`enum` · required

(Required) The DNS record type.
Supported types: A, AAAA, CNAME, MX, TXT, SRV, NS, PTR, CAA, SOA, SPF, SSHFP, NAPTR.

- rule: record_type must be specified (cannot be record_type_unspecified)
- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `record_type_unspecified` -- Unspecified record type (invalid).
- `A` -- IPv4 address record.
- `AAAA` -- IPv6 address record.
- `CNAME` -- Canonical name (alias) record.
- `MX` -- Mail exchange record.
- `TXT` -- Text record (SPF, DKIM, verification, etc.).
- `SRV` -- Service locator record.
- `NS` -- Nameserver record.
- `PTR` -- Pointer record (reverse DNS).
- `CAA` -- Certificate Authority Authorization record.
- `SOA` -- Start of Authority record.
- `SPF` -- Sender Policy Framework record (deprecated in favor of TXT).
- `SSHFP` -- SSH Public Key Fingerprint record.
- `NAPTR` -- Naming Authority Pointer record.

### spec.records[].recordName

`string` · required

(Required) The fully qualified domain name for this record.
Must end with a trailing dot to indicate FQDN.
Example: "www.example.com." or "api.example.com."

- rule: record_name must be a valid DNS domain name ending with a trailing dot (e.g., www.example.com.)
- rule: {"required":true}

### spec.records[].values

`[]string` · required

(Required) The DNS record values.
For A records: IPv4 addresses (e.g., "192.0.2.1")
For AAAA records: IPv6 addresses (e.g., "2001:db8::1")
For CNAME records: target hostname with trailing dot (e.g., "target.example.com.")
For MX records: priority and mail server (e.g., "10 mail.example.com.")
For TXT records: text values (e.g., "v=spf1 include:_spf.google.com ~all")
Multiple values create a round-robin record set.

- rule: {"repeated":{"minItems":"1"}}

### spec.records[].ttl

`int32` · optional (explicit presence)

(Optional) Time To Live (in seconds) for this specific record.
Overrides the zone-level TTL for this record.
Default: 60 seconds.

- default: `60`

### spec.region

`string`

(Optional) Override the region from the provider config for this zone.
If omitted, the region from the OpenStack provider config is used.
ForceNew: changing this requires recreating the zone.

## Validation Rules

- `secondary.requires_masters`: SECONDARY zones require at least one master nameserver in the masters field
- `type.valid_values`: type must be 'PRIMARY' or 'SECONDARY' when specified

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_id` | `string` | zone_id is the unique identifier (UUID) of the DNS zone. This is the primary output used as a foreign key by DNS record components. |
| `status.outputs.zone_name` | `string` | zone_name is the DNS zone name (derived from domain_name in spec). This is the authoritative domain managed by this zone. |
| `status.outputs.region` | `string` | region is the OpenStack region where the DNS zone was created. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackDnsRecord | `spec.zoneId` | `status.outputs.zone_id` |
