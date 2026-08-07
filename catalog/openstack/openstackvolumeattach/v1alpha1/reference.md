# OpenStackVolumeAttach

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackVolumeAttachSpec defines the configuration for attaching an
OpenStack Cinder volume to a compute instance.

A volume attachment is a "join" resource -- it connects a volume (persistent
block storage) to an instance (compute). The attachment makes the volume
appear as a block device inside the instance (e.g., /dev/vdb).

All fields on the underlying resource are ForceNew: any change recreates
the attachment. This is expected -- you cannot move a volume between instances
without detaching and reattaching.

The resource name in Planton is derived from metadata.name. OpenStack volume
attachments do not have a user-visible "name" attribute -- the resource is
identified by the combination of instance and volume.

Terraform resource: openstack_compute_volume_attach_v2
Pulumi resource: openstack.compute.VolumeAttach

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackVolumeAttach
metadata:
  name: test-va
spec:
  instance_id:
    value: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
  volume_id:
    value: "b2c3d4e5-f6a7-8901-bcde-f12345678901"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.instanceId` | `string \| valueFrom` | yes |  | OpenStackInstance (`status.outputs.instance_id`) |
| `spec.volumeId` | `string \| valueFrom` | yes |  | OpenStackVolume (`status.outputs.volume_id`) |
| `spec.device` | `string` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.instanceId

`string | valueFrom` · required

instance_id is the ID of the compute instance to attach the volume to.
This is a required foreign key -- every volume attachment must reference
exactly one instance.
Can reference an OpenStackInstance resource's output or be a literal instance UUID.

- references: OpenStackInstance (`status.outputs.instance_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackInstance, name: <that resource's name>, fieldPath: status.outputs.instance_id}} -- a bare string does not parse

### spec.volumeId

`string | valueFrom` · required

volume_id is the ID of the Cinder volume to attach.
This is a required foreign key -- every volume attachment must reference
exactly one volume.
Can reference an OpenStackVolume resource's output or be a literal volume UUID.

- references: OpenStackVolume (`status.outputs.volume_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackVolume, name: <that resource's name>, fieldPath: status.outputs.volume_id}} -- a bare string does not parse

### spec.device

`string`

device is the device path where the volume appears inside the instance.
Example: "/dev/vdb", "/dev/vdc".
If omitted, OpenStack (Nova) automatically selects the next available device.
Computed by OpenStack if not specified.

### spec.region

`string`

region overrides the region from the provider config for this attachment.
If omitted, the region from the OpenStack provider config is used.
ForceNew: changing the region recreates the attachment.
Example: "RegionOne"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackVolumeAttach, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.id` | `string` | id is the Terraform resource identifier for the attachment. Format: "{instance_id}/{volume_id}" |
| `status.outputs.instance_id` | `string` | instance_id is the UUID of the compute instance the volume is attached to. |
| `status.outputs.volume_id` | `string` | volume_id is the UUID of the Cinder volume that was attached. |
| `status.outputs.device` | `string` | device is the device path where the volume appears inside the instance. Example: "/dev/vdb" Computed by OpenStack if not explicitly specified in the spec. |
| `status.outputs.region` | `string` | region is the OpenStack region where the attachment was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.instanceId` | OpenStackInstance | `status.outputs.instance_id` |
| `spec.volumeId` | OpenStackVolume | `status.outputs.volume_id` |

## See Also

- [Overview](../README.md)
