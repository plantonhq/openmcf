# OciDnsRecord

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciDnsRecordSpec defines the specification for an OCI DNS Record Set
(RRSet) -- a set of DNS resource records sharing the same (domain,
rtype) tuple within an OCI DNS zone.

The resource manages records atomically: updates replace the entire
set, not individual records.

Key behaviors:
  - zone_name_or_id, domain, rtype, and view_id are ForceNew
    (changing them destroys and recreates the record set)
  - items are updatable (the full set is replaced on each update)
  - RDATA may be normalized by the OCI service (e.g., IPv6
    compression, trailing-dot appended to hostnames, TXT quote
    stripping) -- returned values may differ from input

Notable deviation from standard OCI component pattern:
  - No compartment_id: the oci_dns_rrset resource infers the
    compartment from the target zone. compartment_id is deprecated
    on this resource in both the Terraform and Pulumi providers.
  - No freeform_tags: DNS record sets do not support OCI tagging.

Items are simplified from the provider model: the provider requires
each item to carry its own domain and rtype (which must match the
top-level values). This spec removes that redundancy -- IaC modules
inject domain and rtype from the top-level fields.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zoneNameOrId` | `string \| valueFrom` | yes |  | OciDnsZone (`status.outputs.zone_id`) |
| `spec.domain` | `string` | yes |  |  |
| `spec.rtype` | `string` | yes |  |  |
| `spec.viewId` | `string \| valueFrom` |  |  |  |
| `spec.items` | `[]RecordItem` | yes |  |  |
| `spec.items[].rdata` | `string` | yes |  |  |
| `spec.items[].ttl` | `int32` |  |  |  |

## Field Details

### spec.zoneNameOrId

`string | valueFrom` · required

OCID or name of the target DNS zone. ForceNew.

- references: OciDnsZone (`status.outputs.zone_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.domain

`string` · required

Fully qualified domain name for the record set (e.g.,
"app.example.com"). ForceNew.

- rule: {"string":{"minLen":"1"}}

### spec.rtype

`string` · required

DNS record type (e.g., A, AAAA, CNAME, MX, TXT, SRV, CAA, NS,
PTR). Plain string because DNS types are IETF-standardized,
numerous, and extensible. ForceNew.

- rule: {"string":{"minLen":"1"}}

### spec.viewId

`string | valueFrom`

OCID of the private DNS view. Required when accessing a private
zone by name. Not needed when zone_name_or_id is an OCID.
No default_kind/default_kind_field_path: OCI private DNS views are not
modeled as an Planton kind, so this stays literal-only with no
cross-resource reference default.
ForceNew.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.items

`[]RecordItem` · required

DNS records in this record set. At least one record is required.

- rule: {"repeated":{"minItems":"1"}}

### spec.items[].rdata

`string` · required

Record data in type-specific presentation format.
Examples: "192.0.2.1" (A), "2001:db8::1" (AAAA),
"mail.example.com." (MX with priority as part of rdata: "10 mail.example.com."),
"\"v=spf1 include:example.com ~all\"" (TXT).

- rule: {"string":{"minLen":"1"}}

### spec.items[].ttl

`int32`

Time to live in seconds. Controls how long resolvers cache
this record. Values below 30 are not recommended by OCI.

- rule: {"int32":{"gte":1}}

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.zoneNameOrId` | OciDnsZone | `status.outputs.zone_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
