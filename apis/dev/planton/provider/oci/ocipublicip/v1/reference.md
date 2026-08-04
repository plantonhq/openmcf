# OciPublicIp

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciPublicIpSpec defines the specification for an Oracle Cloud Infrastructure
public IP address.

A public IP is an Oracle-assigned IPv4 address that enables internet
connectivity for OCI resources. Two lifetime modes exist:

  - **RESERVED**: A persistent, region-scoped IP with a user-controlled
    lifecycle. It survives instance termination and can be reassigned to
    different private IPs over time. Use reserved IPs for production workloads
    that need a stable, well-known address (DNS records, firewall allowlists).

  - **EPHEMERAL**: An IP whose lifecycle is tied to the assigned entity.
    It is automatically released when the instance, VNIC, or NAT gateway it
    is attached to is terminated. Always scoped to an availability domain
    when assigned to a private IP.

This component manages a single public IP. Reserved IPs can optionally be
drawn from a Bring-Your-Own-IP (BYOIP) pool via public_ip_pool_id.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.lifetime` | `string` |  |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.privateIpId` | `string \| valueFrom` |  |  |  |
| `spec.publicIpPoolId` | `string \| valueFrom` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where this public IP will reside.
For ephemeral IPs, must match the compartment of the private IP.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.lifetime

`string`

Lifetime controls when the public IP is deleted and released back to the
OCI public IP pool. Cannot be changed after creation.
Valid values: "RESERVED", "EPHEMERAL".

- rule: {"string":{"in":["RESERVED","EPHEMERAL"]}}

### spec.displayName

`string`

User-friendly display name for the public IP.
Falls back to metadata.name if not provided.

### spec.privateIpId

`string | valueFrom`

OCID of the private IP to assign this public IP to.

Required for ephemeral IPs (must be a primary private IP on a VNIC).
Optional for reserved IPs -- when omitted the IP is created unassigned
and can be attached later via the OCI Console or API.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.publicIpPoolId

`string | valueFrom`

OCID of a public IP pool for BYOIP (Bring Your Own IP) scenarios.
When set, the public IP is allocated from the specified pool instead
of Oracle's pool. Cannot be changed after creation.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

## Validation Rules

- `ephemeral_requires_private_ip`: private_ip_id is required when lifetime is EPHEMERAL

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciPublicIp, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.public_ip_id` | `string` | OCID of the created public IP resource. |
| `status.outputs.ip_address` | `string` | The allocated IPv4 address (e.g. "203.0.113.2"). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
