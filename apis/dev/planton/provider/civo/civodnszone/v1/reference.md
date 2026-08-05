# CivoDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1`

CivoDnsZoneSpec defines the specification required to create a DNS zone (domain) on Civo.
This allows you to manage DNS records for a given domain via Civo's DNS service, focusing on the essential parameters (80/20 principle).

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.domainName` | `string` | yes |  |  |
| `spec.records` | `[]CivoDnsZoneRecord` |  |  |  |
| `spec.records[].name` | `string` | yes |  |  |
| `spec.records[].values` | `[]string \| valueFrom` | yes |  |  |
| `spec.records[].ttlSeconds` | `uint32` |  | `3600` |  |
| `spec.records[].type` | `enum` | yes |  |  |

## Field Details

### spec.domainName

`string` · required

The domain name for the DNS zone.
Must be a valid fully-qualified domain name (e.g., "example.com").

- rule: {"required":true,"string":{"pattern":"^(?:[A-Za-z0-9-]+\\.)+[A-Za-z]{2,}$"}}

### spec.records

`[]CivoDnsZoneRecord`

A list of DNS records to create within the zone (optional).
Each record includes its type, name, value(s), and TTL.

### spec.records[].name

`string` · required

The host/name for the DNS record, relative to the zone.
For root (apex) records, use "@" to denote the zone itself.

- rule: {"required":true}

### spec.records[].values

`[]string | valueFrom` · required

The value or values for the DNS record.
- For A/AAAA: one or more IP addresses.
- For CNAME: the target domain name.
- For TXT: the text data (if multiple strings, they will be concatenated by DNS).
- For MX: one or more entries like "<priority> <mail-server-domain>".
Each value can be a literal or a reference to another resource’s output.

- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.records[].ttlSeconds

`uint32`

The time-to-live (TTL) for this DNS record, in seconds.
Determines how long resolvers cache the record. Defaults to 3600 seconds (1 hour) if not set.

- default: `3600`

### spec.records[].type

`enum` · required

The DNS record type.
This field is required and must be one of the supported record types (A, AAAA, CNAME, MX, TXT, etc.).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `unspecified`
- `A` -- host address
- `AAAA` -- ipv6 host address
- `ALIAS` -- auto resolved alias
- `CNAME` -- canonical name for an alias
- `MX` -- mail exchange
- `NS` -- name server
- `PTR` -- pointer
- `SOA` -- start of authority
- `SRV` -- location of service
- `TXT` -- descriptive text
- `CAA` -- certificate authority authorization

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_name` | `string` | The domain name of the DNS zone managed on Civo. |
| `status.outputs.zone_id` | `string` | The unique identifier of the DNS zone on Civo (UUID). |
| `status.outputs.name_servers` | `[]string` | The list of nameserver addresses for the DNS zone. These are the nameservers that should be set at the domain's registrar (e.g., ns0.civo.com, ns1.civo.com). |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| CivoDnsRecord | `spec.zoneId` | `status.outputs.zone_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
