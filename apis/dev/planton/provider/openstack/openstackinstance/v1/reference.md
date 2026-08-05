# OpenStackInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackInstanceSpec defines the configuration for an OpenStack Compute instance.

An instance is a virtual machine running on an OpenStack cloud. It is the most
important resource after networking -- every developer workload, application
server, and database runs as an instance.

The instance supports two networking modes:
  - Network-based: specify a network UUID and let OpenStack auto-create a port
  - Port-based: specify a pre-provisioned port for stable network identity

The instance supports two boot modes:
  - Image-based: boot from a Glance image (ephemeral root disk)
  - Volume-based: boot from a Cinder volume via block_device (persistent root disk)

The instance name is derived from metadata.name.

Terraform resource: openstack_compute_instance_v2
Pulumi resource: openstack.compute.Instance

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackInstance
metadata:
  name: test-instance
spec:
  flavor_name: m1.medium
  image_name: ubuntu-22.04
  networks:
    - uuid:
        value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  key_pair:
    value: "test-keypair"
  security_groups:
    - value: "default"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.flavorName` | `string` |  |  |  |
| `spec.flavorId` | `string` |  |  |  |
| `spec.imageName` | `string` |  |  |  |
| `spec.imageId` | `string` |  |  |  |
| `spec.keyPair` | `string \| valueFrom` |  |  | OpenStackKeypair (`status.outputs.name`) |
| `spec.networks` | `[]InstanceNetwork` |  |  |  |
| `spec.networks[].uuid` | `string \| valueFrom` |  |  | OpenStackNetwork (`status.outputs.network_id`) |
| `spec.networks[].port` | `string \| valueFrom` |  |  | OpenStackNetworkPort (`status.outputs.port_id`) |
| `spec.networks[].fixedIpV4` | `string` |  |  |  |
| `spec.networks[].accessNetwork` | `bool` |  |  |  |
| `spec.securityGroups` | `[]string \| valueFrom` |  |  | OpenStackSecurityGroup (`status.outputs.name`) |
| `spec.blockDevice` | `[]BlockDevice` |  |  |  |
| `spec.blockDevice[].sourceType` | `string` | yes |  |  |
| `spec.blockDevice[].uuid` | `string` |  |  |  |
| `spec.blockDevice[].destinationType` | `string` |  |  |  |
| `spec.blockDevice[].bootIndex` | `int32` |  |  |  |
| `spec.blockDevice[].volumeSize` | `int32` |  |  |  |
| `spec.blockDevice[].deleteOnTermination` | `bool` |  |  |  |
| `spec.blockDevice[].volumeType` | `string` |  |  |  |
| `spec.userData` | `string` |  |  |  |
| `spec.metadata` | `map<string, string>` |  |  |  |
| `spec.configDrive` | `bool` |  |  |  |
| `spec.serverGroupId` | `string \| valueFrom` |  |  | OpenStackServerGroup (`status.outputs.server_group_id`) |
| `spec.availabilityZone` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.flavorName

`string`

flavor_name is the human-readable name of the instance flavor (e.g., "m1.medium").
This is the most common way to specify instance size.
Mutually exclusive with flavor_id.

### spec.flavorId

`string`

flavor_id is the UUID of the instance flavor.
Use when you need to reference a specific flavor by ID rather than name.
Mutually exclusive with flavor_name.

### spec.imageName

`string`

image_name is the name of the Glance image to boot from (e.g., "ubuntu-22.04").
This is the most common way to specify the boot image.
Optional when using block_device with a boot volume.

### spec.imageId

`string`

image_id is the UUID of the Glance image to boot from.
Alternative to image_name for UUID-based image references.
Optional when using block_device with a boot volume.

### spec.keyPair

`string | valueFrom`

key_pair is the name of the SSH keypair to inject into the instance.
Can reference an OpenStackKeypair resource's output name (via value_from
for InfraChart DAG wiring) or be a literal keypair name.
Optional: instances without SSH keys rely on other access methods
(e.g., console, cloud-init password).

- references: OpenStackKeypair (`status.outputs.name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackKeypair, name: <that resource's name>, fieldPath: status.outputs.name}} -- a bare string does not parse

### spec.networks

`[]InstanceNetwork`

networks defines the network attachments for the instance.
Each entry connects the instance to a network (via UUID) or attaches
a pre-provisioned port (via port). At least one network is required.

- rule: exactly one of uuid or port must be set

### spec.networks[].uuid

`string | valueFrom`

uuid is the network UUID to attach the instance to.
OpenStack auto-creates a port on this network for the instance.
Can reference an OpenStackNetwork resource's output (via value_from).
Mutually exclusive with port.

- references: OpenStackNetwork (`status.outputs.network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackNetwork, name: <that resource's name>, fieldPath: status.outputs.network_id}} -- a bare string does not parse

### spec.networks[].port

`string | valueFrom`

port is the UUID of a pre-provisioned port to attach to the instance.
Use when you need stable MAC/IP addresses, specific security groups,
or other port-level configuration.
Can reference an OpenStackNetworkPort resource's output (via value_from).
Mutually exclusive with uuid.

- references: OpenStackNetworkPort (`status.outputs.port_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackNetworkPort, name: <that resource's name>, fieldPath: status.outputs.port_id}} -- a bare string does not parse

### spec.networks[].fixedIpV4

`string`

fixed_ip_v4 requests a specific IPv4 address on the network.
Only meaningful when uuid is used (not port -- port has its own fixed_ips).
If omitted, an available IP from the subnet's allocation pool is assigned.

### spec.networks[].accessNetwork

`bool`

access_network marks this network as the instance's access network.
The access network determines which IP appears in access_ip_v4 output.
Only one network should be marked as access_network.

### spec.securityGroups

`[]string | valueFrom`

security_groups is the list of security group names to apply to the instance.
Each entry can reference an OpenStackSecurityGroup resource's output name
(via value_from) or be a literal security group name.
If omitted, OpenStack applies the default security group.

Note: the Compute API uses security group NAMES (not UUIDs). The FK
resolves to status.outputs.name, not status.outputs.security_group_id.
The literal value must be a security group name, not a UUID.

- references: OpenStackSecurityGroup (`status.outputs.name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.name}} -- a bare string does not parse

### spec.blockDevice

`[]BlockDevice`

block_device defines block device mappings for the instance.
Use this for boot-from-volume (persistent root disk) or attaching
additional volumes at launch time.
When using block_device with boot_index=0, image_name/image_id are optional.

### spec.blockDevice[].sourceType

`string` · required

source_type is the type of the source for this block device.
Required.
  "image"    - boot from a Glance image (most common)
  "volume"   - attach an existing Cinder volume
  "snapshot" - create from a volume snapshot
  "blank"    - create an empty volume (for additional storage)

- rule: {"string":{"minLen":"1","in":["blank","image","snapshot","volume"]}}

### spec.blockDevice[].uuid

`string`

uuid is the UUID of the source (image, volume, or snapshot).
Required when source_type is not "blank".

### spec.blockDevice[].destinationType

`string`

destination_type controls where the block device is created.
  "local"  - ephemeral storage on the hypervisor (default for image)
  "volume" - persistent Cinder volume (recommended for production)

### spec.blockDevice[].bootIndex

`int32`

boot_index determines the boot order.
  0  - this is the boot device
 -1  - not bootable (data volume)
  N  - boot priority (lower = higher priority)

### spec.blockDevice[].volumeSize

`int32`

volume_size is the size of the block device in GB.
Required for image-to-volume and blank-to-local mappings.

### spec.blockDevice[].deleteOnTermination

`bool`

delete_on_termination controls whether the volume is deleted when
the instance is terminated. Default: false (volume persists).

### spec.blockDevice[].volumeType

`string`

volume_type specifies the Cinder volume type (e.g., "SSD", "HDD").
Only meaningful when destination_type is "volume".

### spec.userData

`string`

user_data is the cloud-init configuration or script to run at first boot.
Accepts cloud-config YAML, shell scripts, or any cloud-init format.
The value is base64-encoded before being passed to the instance.
ForceNew: changing user_data recreates the instance.

### spec.metadata

`map<string, string>`

metadata is a map of key-value pairs to attach to the instance.
Metadata is visible in the OpenStack API and Horizon dashboard.
Can be updated without recreating the instance.

### spec.configDrive

`bool` · optional (explicit presence)

config_drive enables the config drive for metadata delivery.
When true, a small read-only disk is attached containing instance
metadata, user_data, and network config. Useful when DHCP is not
available or for environments requiring metadata on a local disk.
ForceNew: changing config_drive recreates the instance.

### spec.serverGroupId

`string | valueFrom`

server_group_id is the UUID of the server group for placement control.
Can reference an OpenStackServerGroup resource's output (via value_from
for InfraChart DAG wiring) or be a literal server group UUID.
Maps to scheduler_hints.group in the Terraform/Pulumi resources.
ForceNew: changing the server group recreates the instance.

- references: OpenStackServerGroup (`status.outputs.server_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackServerGroup, name: <that resource's name>, fieldPath: status.outputs.server_group_id}} -- a bare string does not parse

### spec.availabilityZone

`string`

availability_zone specifies the AZ where the instance should be launched.
If omitted, Nova selects an AZ based on its scheduling algorithms.
ForceNew: changing the AZ recreates the instance.
Example: "nova", "az1", "az:host:node"

### spec.tags

`[]string`

tags are string tags to associate with the instance in OpenStack.
Tags are stored on the OpenStack resource and can be used for filtering
and organization in the OpenStack API and Horizon dashboard.

- rule: {"repeated":{"unique":true}}

### spec.region

`string`

region overrides the region from the provider config for this instance.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Validation Rules

- `flavor.exactly_one`: exactly one of flavor_name or flavor_id must be set
- `networks.required`: at least one network must be specified

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | instance_id is the unique identifier (UUID) of the instance in OpenStack. This is the primary output used as a foreign key by downstream components. |
| `status.outputs.name` | `string` | name is the name of the instance (derived from metadata.name). |
| `status.outputs.access_ip_v4` | `string` | access_ip_v4 is the best IPv4 address for accessing the instance. OpenStack computes this from the instance's network interfaces, prioritizing the access_network if one is marked. |
| `status.outputs.access_ip_v6` | `string` | access_ip_v6 is the best IPv6 address for accessing the instance. Empty if the instance has no IPv6 connectivity. |
| `status.outputs.region` | `string` | region is the OpenStack region where the instance was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.keyPair` | OpenStackKeypair | `status.outputs.name` |
| `spec.networks[].uuid` | OpenStackNetwork | `status.outputs.network_id` |
| `spec.networks[].port` | OpenStackNetworkPort | `status.outputs.port_id` |
| `spec.securityGroups` | OpenStackSecurityGroup | `status.outputs.name` |
| `spec.serverGroupId` | OpenStackServerGroup | `status.outputs.server_group_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackVolumeAttach | `spec.instanceId` | `status.outputs.instance_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
