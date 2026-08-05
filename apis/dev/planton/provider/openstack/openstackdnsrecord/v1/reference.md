# OpenStackDnsRecord

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackDnsRecordSpec defines the configuration for a standalone DNS record
in an OpenStack Designate DNS zone.

This component creates a single DNS recordset (one name + one type + one or more values)
within an existing zone. It is the DAG-visible counterpart to the inline records
supported by OpenStackDnsZone -- use this component when records need to be
independently managed or wired as explicit dependencies in InfraCharts.

Terraform resource: openstack_dns_recordset_v2
Pulumi resource:    dns.RecordSet

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackDnsRecord
metadata:
  name: app-a-record
spec:
  zone_id:
    value: "zone-uuid-1234"
  record_name: "app.example.com."
  type: 1
  values:
    - "192.0.2.1"
  description: "Test DNS record for local development"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zoneId` | `string \| valueFrom` | yes |  | OpenStackDnsZone (`status.outputs.zone_id`) |
| `spec.recordName` | `string` | yes |  |  |
| `spec.type` | `enum` | yes |  |  |
| `spec.values` | `[]string` | yes |  |  |
| `spec.ttl` | `int32` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.zoneId

`string | valueFrom` · required

(Required) The ID of the Designate zone where this record will be created.
Can be a literal UUID or a reference to an OpenStackDnsZone resource's output.

FK: OpenStackDnsZone.status.outputs.zone_id
ForceNew: changing this requires recreating the record.

- references: OpenStackDnsZone (`status.outputs.zone_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.recordName

`string` · required

(Required) The fully qualified domain name for this record.
Must end with a trailing dot to indicate FQDN.
Example: "www.example.com." or "api.example.com."
ForceNew: changing this requires recreating the record.

- rule: record_name must be a valid DNS domain name ending with a trailing dot (e.g., www.example.com.)
- rule: {"required":true}

### spec.type

`enum` · required

(Required) The DNS record type.
Supported types: A, AAAA, CNAME, MX, TXT, SRV, NS, PTR, CAA, SOA, SPF, SSHFP, NAPTR.
ForceNew: changing this requires recreating the record.

- rule: type must be specified (cannot be record_type_unspecified)
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

### spec.values

`[]string` · required

(Required) The DNS record values.
Format depends on record type:
  - A record: IPv4 addresses (e.g., ["192.0.2.1", "192.0.2.2"])
  - AAAA record: IPv6 addresses (e.g., ["2001:db8::1"])
  - CNAME record: target hostname with trailing dot (e.g., ["target.example.com."])
  - MX record: priority and mail server (e.g., ["10 mail1.example.com.", "20 mail2.example.com."])
  - TXT record: text values (e.g., ["v=spf1 include:_spf.google.com ~all"])
Multiple values create a round-robin record set.

- rule: {"repeated":{"minItems":"1"}}

### spec.ttl

`int32` · optional (explicit presence)

(Optional) Time To Live (in seconds) for this DNS record.
Determines how long resolvers cache this record.
If omitted, the zone's default TTL is used.

### spec.description

`string`

(Optional) Human-readable description of the DNS record.

### spec.region

`string`

(Optional) Override the region from the provider config for this record.
If omitted, the region from the OpenStack provider config is used.
ForceNew: changing this requires recreating the record.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackDnsRecord, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.recordset_id` | `string` | recordset_id is the unique identifier (UUID) of the DNS recordset. |
| `status.outputs.fqdn` | `string` | fqdn is the fully qualified domain name of the created DNS record. Example: "www.example.com." or "api.example.com." |
| `status.outputs.record_type` | `string` | record_type is the DNS record type that was created. Example: "A", "AAAA", "CNAME", "MX", "TXT". |
| `status.outputs.values` | `[]string` | values is the list of DNS record values that were set. |
| `status.outputs.zone_id` | `string` | zone_id is the ID of the zone containing this record. |
| `status.outputs.region` | `string` | region is the OpenStack region where the record was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.zoneId` | OpenStackDnsZone | `status.outputs.zone_id` |
