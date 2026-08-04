# DigitalOceanDnsRecord

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1`

DigitalOceanDnsRecordSpec defines the configuration for creating a DNS record in a DigitalOcean DNS zone (domain).
This component supports creating individual DNS records with common record types including A, AAAA, CNAME, MX, TXT, SRV, NS, and CAA.
DNS records are created within an existing DigitalOcean DNS zone (domain) managed by DigitalOcean's DNS service.

## Example

```yaml
apiVersion: digital-ocean.planton.dev/v1
kind: DigitalOceanDnsRecord
metadata:
  name: test-dns-record
spec:
  domain:
    value: "example.com"
  name: "www"
  type: A
  value:
    value: "192.0.2.1"
  ttl_seconds: 3600
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.domain` | `string \| valueFrom` | yes |  | DigitalOceanDnsZone (`status.outputs.zone_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.type` | `enum` | yes |  |  |
| `spec.value` | `string \| valueFrom` | yes |  |  |
| `spec.ttlSeconds` | `int32` |  | `1800` |  |
| `spec.priority` | `int32` |  |  |  |
| `spec.weight` | `int32` |  |  |  |
| `spec.port` | `int32` |  |  |  |
| `spec.flags` | `int32` |  |  |  |
| `spec.tag` | `string` |  |  |  |

## Field Details

### spec.domain

`string | valueFrom` · required

The DigitalOcean domain name (DNS zone) where this DNS record will be created.
This can be a direct value or a reference to a DigitalOceanDnsZone resource's output.
Example: "example.com"

- references: DigitalOceanDnsZone (`status.outputs.zone_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the DNS record (hostname/subdomain).
Use "@" for root domain records, or specify the subdomain name.
Examples:
  - "@" for root domain (example.com)
  - "www" for subdomain (www.example.com)
  - "api.v1" for nested subdomain (api.v1.example.com)

- rule: {"required":true}

### spec.type

`enum` · required

The type of DNS record to create.
Supported types: A, AAAA, CNAME, MX, TXT, SRV, NS, CAA.

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
- `CAA` -- Certificate Authority Authorization record.

### spec.value

`string | valueFrom` · required

The value/target of the DNS record.
This can be a direct value or a reference to another resource's output.
Format depends on record type:
  - A record: IPv4 address (e.g., "192.0.2.1")
  - AAAA record: IPv6 address (e.g., "2001:db8::1")
  - CNAME record: Target hostname (e.g., "target.example.com")
  - MX record: Mail server hostname (e.g., "mail.example.com")
  - TXT record: Text value (e.g., "v=spf1 include:_spf.google.com ~all")
  - NS record: Nameserver hostname
  - SRV record: Target hostname
  - CAA record: CA domain (e.g., "letsencrypt.org")

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.ttlSeconds

`int32` · optional (explicit presence)

Time to live (TTL) for the DNS record in seconds.
Determines how long DNS resolvers should cache this record.
Common values: 60 (1 min), 300 (5 min), 1800 (30 min), 3600 (1 hour), 86400 (1 day).
Default: 1800 seconds (30 minutes).
Valid range: 30-86400 seconds.

- default: `1800`
- rule: {"int32":{"lte":86400,"gte":30}}

### spec.priority

`int32`

Priority for MX and SRV records.
Lower values indicate higher priority.
For MX records: Typically 10, 20, 30 for primary, secondary, tertiary mail servers.
For SRV records: Used with weight for load distribution.
Range: 0-65535. Defaults to 0 if not specified.

- rule: {"int32":{"lte":65535,"gte":0}}

### spec.weight

`int32`

Weight for SRV records.
Specifies the relative weight for records with the same priority.
Higher weights receive proportionally more traffic.
Range: 0-65535. Defaults to 0 if not specified.
Ignored for non-SRV record types.

- rule: {"int32":{"lte":65535,"gte":0}}

### spec.port

`int32`

Port for SRV records.
Specifies the TCP or UDP port on which the service is available.
Range: 0-65535.
Required for SRV records, ignored for other types.

- rule: {"int32":{"lte":65535,"gte":0}}

### spec.flags

`int32`

Flags for CAA records.
Values:
  - 0: Non-critical (default) - CA may ignore unknown tags
  - 128: Critical - CA must refuse if tag is not understood
Ignored for non-CAA record types.

- rule: {"int32":{"lte":255,"gte":0}}

### spec.tag

`string`

Tag for CAA records.
Specifies the property being authorized:
  - "issue": Authorizes a CA to issue certificates for the domain
  - "issuewild": Authorizes a CA to issue wildcard certificates
  - "iodef": URL to report policy violations
Required for CAA records, ignored for other types.

## Validation Rules

- `spec.port_required_for_srv`: port must be specified for SRV records
- `spec.tag_required_for_caa`: tag must be specified for CAA records (issue, issuewild, or iodef)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanDnsRecord, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.record_id` | `string` | The unique identifier of the created DNS record in DigitalOcean. This ID is assigned by DigitalOcean and can be used for API operations. |
| `status.outputs.hostname` | `string` | The fully qualified hostname of the DNS record. For example: "www.example.com" or "example.com" for root domain records. |
| `status.outputs.record_type` | `string` | The DNS record type that was created. One of: A, AAAA, CNAME, MX, TXT, SRV, NS, CAA. |
| `status.outputs.domain` | `string` | The domain name (DNS zone) where the record was created. |
| `status.outputs.ttl_seconds` | `int32` | The TTL (time to live) in seconds that was applied to the record. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.domain` | DigitalOceanDnsZone | `status.outputs.zone_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
