# OciBastion

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1alpha1`

OciBastionSpec defines the specification for an OCI Bastion --
a managed SSH gateway that provides secure, time-limited access
to resources in a private subnet without requiring a public IP
on the target resource.

The bastion creates a private endpoint in the target subnet and
allows authorized clients (by CIDR) to establish SSH sessions
(managed SSH, port forwarding, dynamic port forwarding) through it.

bastion_type is hardcoded to STANDARD in IaC modules (the only
publicly documented type). Fields specific to non-standard types
(phone_book_entry, static_jump_host_ip_addresses) are excluded.

Sessions (oci_bastion_session) are NOT bundled -- they are
ephemeral operational artifacts with TTLs, not infrastructure.
The bastion's max_session_ttl_in_seconds governs the upper bound.

Excluded from v1:
  - phone_book_entry -- not applicable to STANDARD bastions
  - static_jump_host_ip_addresses -- not applicable to STANDARD bastions
  - security_attributes -- Oracle ZPR, very low adoption
  - defined_tags, system_tags -- managed by platform
  - freeform_tags -- auto-populated from metadata labels

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.targetSubnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.clientCidrBlockAllowList` | `[]string` |  |  |  |
| `spec.maxSessionTtlInSeconds` | `int32` |  |  |  |
| `spec.isDnsProxyEnabled` | `bool` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the bastion will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.targetSubnetId

`string | valueFrom` · required

OCID of the subnet that the bastion connects to. The bastion
creates a private endpoint in this subnet. Immutable after creation.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.displayName

`string`

Display name for the bastion. When omitted, the metadata name is used.
Immutable after creation.

### spec.clientCidrBlockAllowList

`[]string`

CIDR ranges allowed to connect to sessions hosted by this bastion.
Example: ["10.0.0.0/16", "192.168.1.0/24"]. Updatable after creation.

### spec.maxSessionTtlInSeconds

`int32` · optional (explicit presence)

Maximum TTL in seconds for any session on this bastion.
When omitted, the OCI default applies (10800 = 3 hours).
Updatable after creation.

### spec.isDnsProxyEnabled

`bool` · optional (explicit presence)

Enable FQDN and SOCKS5 proxy support on the bastion.
When true, sessions can use DNS names to reach targets.
Immutable after creation. When omitted, defaults to disabled.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciBastion, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.bastion_id` | `string` | OCID of the bastion. |
| `status.outputs.private_endpoint_ip_address` | `string` | Private IP address of the bastion's endpoint in the target subnet. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.targetSubnetId` | OciSubnet | `status.outputs.subnet_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
