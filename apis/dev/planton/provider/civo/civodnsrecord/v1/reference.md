# CivoDnsRecord

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1`

CivoDnsRecordSpec defines the configuration for creating a DNS record in a Civo DNS zone.
This component supports creating individual DNS records with common record types.

## Example

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoDnsRecord
metadata:
  name: test-dns-record
spec:
  zoneId:
    value: test-zone-id-abc123
  name: "www"
  type: A
  value: "192.0.2.1"
  ttl: 3600
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zoneId` | `string \| valueFrom` | yes |  | CivoDnsZone (`status.outputs.zone_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.type` | `enum` | yes |  |  |
| `spec.value` | `string` | yes |  |  |
| `spec.ttl` | `int32` |  |  |  |
| `spec.priority` | `int32` |  |  |  |

## Field Details

### spec.zoneId

`string | valueFrom` · required

The ID of the Civo DNS Zone where this record will be created.
Can be provided as a literal value or referenced from a CivoDnsZone resource.

When using value_from, the default kind is CivoDnsZone and the default field path
is "status.outputs.zone_id", allowing you to wire this directly to a zone resource.

Example value_from:
  value_from:
    name: my-zone

- references: CivoDnsZone (`status.outputs.zone_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CivoDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.name

`string` · required

The name of the DNS record (e.g., "www", "api", "@" for root).
Use "@" to create a record at the zone apex (root domain).

- rule: {"required":true}

### spec.type

`enum` · required

The type of DNS record to create.

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

### spec.value

`string` · required

The value/target of the DNS record.
For A records: IPv4 address (e.g., "192.0.2.1")
For AAAA records: IPv6 address (e.g., "2001:db8::1")
For CNAME records: target hostname (e.g., "example.com")
For MX records: mail server hostname (e.g., "mail.example.com")
For TXT records: text value (e.g., "v=spf1 include:_spf.google.com ~all")
For SRV records: target in format "priority weight port target"

- rule: {"required":true}

### spec.ttl

`int32`

Time to live (TTL) for the DNS record in seconds.
Determines how long resolvers cache the record.
Valid range: 60-86400 seconds (1 minute to 24 hours).
Defaults to 3600 (1 hour) if not specified.

- rule: ttl must be 0 (default) or between 60 and 86400 seconds

### spec.priority

`int32`

Priority for MX and SRV records.
Lower values indicate higher priority.
Required for MX records, optional for SRV records.
Range: 0-65535.

- rule: priority must be between 0 and 65535

## Validation Rules

- `spec.priority_required_for_mx`: priority must be set (> 0) for MX records

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoDnsRecord, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.record_id` | `string` | The unique identifier of the created DNS record in Civo. |
| `status.outputs.hostname` | `string` | The fully qualified hostname of the DNS record. For example: "www.example.com" or "example.com" for apex records. |
| `status.outputs.record_type` | `string` | The DNS record type that was created. |
| `status.outputs.account_id` | `string` | The account ID in Civo where the record was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.zoneId` | CivoDnsZone | `status.outputs.zone_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
