# CivoComputeInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1`

CivoComputeInstanceSpec defines the user configuration for a Civo compute instance.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.instanceName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.size` | `string` | yes |  |  |
| `spec.image` | `string` | yes |  |  |
| `spec.network` | `string \| valueFrom` | yes |  | CivoVpc (`status.outputs.network_id`) |
| `spec.sshKeyIds` | `[]string` |  |  |  |
| `spec.firewallIds` | `[]string \| valueFrom` |  |  | CivoFirewall (`status.outputs.firewall_id`) |
| `spec.volumeIds` | `[]string \| valueFrom` |  |  | CivoVolume (`status.outputs.volume_id`) |
| `spec.reservedIpId` | `string \| valueFrom` |  |  | CivoIpAddress (`status.outputs.reserved_ip_id`) |
| `spec.tags` | `[]string` |  |  |  |
| `spec.userData` | `string` |  |  |  |

## Field Details

### spec.instanceName

`string` · required

instance hostname (letters, numbers, dashes, dots; <=63 chars, no trailing dash/dot)

- rule: {"required":true,"string":{"maxLen":"63","pattern":"^[a-z0-9]([a-z0-9\\.\\-]*[a-z0-9])?$"}}

### spec.region

`enum` · required

region code for the instance location

- rule: {"required":true}

Allowed values (use exactly as shown):

- `civo_region_unspecified` -- 0: default / unspecified region
- `lon1` -- london 1
- `lon2` -- london 2
- `fra1` -- frankfurt 1
- `nyc1` -- new york 1
- `phx1` -- phoenix 1
- `mum1` -- mumbai 1

### spec.size

`string` · required

instance size (flavor) slug, e.g. "g3.small"

- rule: {"required":true,"string":{"pattern":"^[a-z0-9]([a-z0-9\\.\\-]*[a-z0-9])?$"}}

### spec.image

`string` · required

base OS image slug for the instance (e.g. "ubuntu-focal")

- rule: {"required":true,"string":{"pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"}}

### spec.network

`string | valueFrom` · required

target network for the instance (must exist in the same region)

- references: CivoVpc (`status.outputs.network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CivoVpc, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.sshKeyIds

`[]string`

SSH public key(s) to enable passwordless login (optional; if empty, a password will be set)

### spec.firewallIds

`[]string | valueFrom`

firewall(s) to attach to the instance (must belong to the same network)

- references: CivoFirewall (`status.outputs.firewall_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CivoFirewall, name: <that resource's name>, fieldPath: status.outputs.firewall_id}} -- a bare string does not parse

### spec.volumeIds

`[]string | valueFrom`

existing storage volumes to attach to this instance (must reside in same region)

- references: CivoVolume (`status.outputs.volume_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CivoVolume, name: <that resource's name>, fieldPath: status.outputs.volume_id}} -- a bare string does not parse

### spec.reservedIpId

`string | valueFrom`

reserved IP to assign to this instance (optional static public IP address)

- references: CivoIpAddress (`status.outputs.reserved_ip_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CivoIpAddress, name: <that resource's name>, fieldPath: status.outputs.reserved_ip_id}} -- a bare string does not parse

### spec.tags

`[]string`

tags for the instance (for organization, optional)

- rule: {"repeated":{"unique":true}}

### spec.userData

`string`

cloud-init user data script to run on instance boot (<=32 KiB, optional)

- rule: {"string":{"maxBytes":"32768"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoComputeInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | unique identifier of the instance (UUID) |
| `status.outputs.public_ipv4` | `string` | public IPv4 address assigned to the instance (if any) |
| `status.outputs.private_ipv4` | `string` | private IPv4 address of the instance within its network |
| `status.outputs.status` | `string` | current status of the instance (e.g., "ACTIVE", "BUILDING") |
| `status.outputs.created_at_rfc3339` | `string` | timestamp when the instance was created (RFC 3339 format) |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.network` | CivoVpc | `status.outputs.network_id` |
| `spec.firewallIds` | CivoFirewall | `status.outputs.firewall_id` |
| `spec.volumeIds` | CivoVolume | `status.outputs.volume_id` |
| `spec.reservedIpId` | CivoIpAddress | `status.outputs.reserved_ip_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
