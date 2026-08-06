# OciDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciDnsZoneSpec defines the specification for an OCI DNS Zone --
a managed authoritative DNS zone supporting both public (GLOBAL)
and private resolution scopes, PRIMARY and SECONDARY zone types,
zone transfers via external masters/downstreams, and DNSSEC signing.

Key behaviors:
  - zone_type, scope, view_id, and name are ForceNew (changing
    them destroys and recreates the zone)
  - compartment_id is updatable (compartment move)
  - external_masters and external_downstreams are updatable
  - is_dnssec_enabled is updatable (toggle after creation)

Constraints enforced via CEL:
  - zone_type must be explicitly set (no implicit default)
  - PRIVATE zones require view_id
  - SECONDARY zones cannot be PRIVATE (OCI limitation)
  - SECONDARY zones require at least one external master

Excluded from v1:
  - dnssec_config -- deeply nested computed DNSSEC key version
    data (KSK/ZSK details); operational concern for key rotation
  - zone_transfer_servers -- computed; not needed for composability
  - is_protected -- read-only system flag
  - DNSSEC key lifecycle actions (stage/promote) -- operational
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.zoneType` | `enum` |  |  |  |
| `spec.scope` | `enum` |  |  |  |
| `spec.viewId` | `string \| valueFrom` |  |  |  |
| `spec.isDnssecEnabled` | `bool` |  |  |  |
| `spec.externalMasters` | `[]ExternalServer` |  |  |  |
| `spec.externalMasters[].address` | `string` | yes |  |  |
| `spec.externalMasters[].port` | `int32` |  |  |  |
| `spec.externalMasters[].tsigKeyId` | `string` |  |  |  |
| `spec.externalDownstreams` | `[]ExternalServer` |  |  |  |
| `spec.externalDownstreams[].address` | `string` | yes |  |  |
| `spec.externalDownstreams[].port` | `int32` |  |  |  |
| `spec.externalDownstreams[].tsigKeyId` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where this zone will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.zoneType

`enum`

Zone type: PRIMARY or SECONDARY. Required. ForceNew.

Allowed values (use exactly as shown):

- `unspecified`
- `primary`
- `secondary`

### spec.scope

`enum`

Resolution scope. When omitted, defaults to GLOBAL (public DNS).
ForceNew.

Allowed values (use exactly as shown):

- `scope_unspecified` -- Treat as GLOBAL (public DNS). This is the default when scope is omitted, matching OCI API behavior.
- `global`
- `scope_private` -- Private DNS zone, resolvable only within VCNs via a DNS view.

### spec.viewId

`string | valueFrom`

OCID of the private DNS view. Required when scope is private.
Not applicable for global zones. ForceNew.
No default_kind/default_kind_field_path: OCI private DNS views are not
modeled as an Planton kind, so this stays literal-only with no
cross-resource reference default.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.isDnssecEnabled

`bool` · optional (explicit presence)

Enable DNSSEC signing for the zone. When true, OCI generates
KSK and ZSK key pairs and signs zone records. Only meaningful
for GLOBAL zones. Nil/omitted = OCI default (disabled).

### spec.externalMasters

`[]ExternalServer`

External master DNS servers that this SECONDARY zone replicates
from. Required when zone_type is secondary. Not applicable for
PRIMARY zones (ignored if set).

### spec.externalMasters[].address

`string` · required

IPv4 or IPv6 address of the external DNS server.

- rule: {"string":{"minLen":"1"}}

### spec.externalMasters[].port

`int32` · optional (explicit presence)

Port number. Must be 53 or omitted (OCI validates server-side).

### spec.externalMasters[].tsigKeyId

`string`

OCID of the TSIG key for authenticating zone transfers.
TSIG keys are not modeled as Planton components.

### spec.externalDownstreams

`[]ExternalServer`

External downstream DNS servers that receive zone transfers
from this PRIMARY zone. Only supported for PRIMARY zones with
GLOBAL scope.

### spec.externalDownstreams[].address

`string` · required

IPv4 or IPv6 address of the external DNS server.

- rule: {"string":{"minLen":"1"}}

### spec.externalDownstreams[].port

`int32` · optional (explicit presence)

Port number. Must be 53 or omitted (OCI validates server-side).

### spec.externalDownstreams[].tsigKeyId

`string`

OCID of the TSIG key for authenticating zone transfers.
TSIG keys are not modeled as Planton components.

## Validation Rules

- `zone_type_required`: zone_type must be explicitly set (primary or secondary)
- `private_requires_view_id`: private zones require view_id to be set
- `no_private_secondary`: SECONDARY zones cannot have private scope (OCI only supports PRIMARY for private zones)
- `secondary_requires_external_masters`: SECONDARY zones require at least one external master

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_id` | `string` | OCID of the DNS zone. |
| `status.outputs.nameservers` | `string` | Comma-separated list of OCI-assigned authoritative nameserver hostnames. Users configure these as NS records at their domain registrar. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OciDnsRecord | `spec.zoneNameOrId` | `status.outputs.zone_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
