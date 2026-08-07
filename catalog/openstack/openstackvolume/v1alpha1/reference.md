# OpenStackVolume

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1alpha1`

OpenStackVolumeSpec defines the configuration for an OpenStack Cinder block
storage volume.

A volume is a persistent block storage device that can be attached to compute
instances. Volumes survive instance termination and can be detached and
reattached to different instances, making them the primary mechanism for
persistent storage in OpenStack.

Volumes can be created empty (blank), cloned from an existing volume, restored
from a snapshot, or initialized from a Glance image (for bootable volumes).
These source options are mutually exclusive.

The volume name is derived from metadata.name.

Terraform resource: openstack_blockstorage_volume_v3
Pulumi resource: openstack.blockstorage.VolumeV3

## Example

```yaml
apiVersion: openstack.planton.dev/v1alpha1
kind: OpenStackVolume
metadata:
  name: test-volume
spec:
  size: 10
  description: "Test volume for development"
  volume_type: "SSD"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.description` | `string` |  |  |  |
| `spec.size` | `int32` |  |  |  |
| `spec.volumeType` | `string` |  |  |  |
| `spec.availabilityZone` | `string` |  |  |  |
| `spec.snapshotId` | `string` |  |  |  |
| `spec.sourceVolId` | `string` |  |  |  |
| `spec.imageId` | `string \| valueFrom` |  |  | OpenStackImage (`status.outputs.image_id`) |
| `spec.metadata` | `map<string, string>` |  |  |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.description

`string`

description is a human-readable description of the volume.
This is stored on the OpenStack resource and visible in Horizon and API responses.

### spec.size

`int32`

size is the volume size in gigabytes (GB).
Required for blank volumes and image-based volumes.
For snapshot-based or clone-based volumes, if omitted or set to 0, the size
of the source is used. If provided, it must be >= the source size.
Must be greater than 0 for blank volumes.

- rule: {"int32":{"gt":0}}

### spec.volumeType

`string`

volume_type is the Cinder volume type (backend storage class).
Examples: "SSD", "HDD", "__DEFAULT__", "ceph-ssd", "lvm".
If omitted, the Cinder default volume type for the project is used.
Changing this on an existing volume triggers a retype operation.

### spec.availabilityZone

`string`

availability_zone is the availability zone where the volume is created.
ForceNew: changing the AZ requires recreating the volume.
If omitted, Cinder selects the default AZ.
Must match the instance's AZ for attachment to work (in most deployments).

### spec.snapshotId

`string`

snapshot_id is the UUID of a Cinder volume snapshot to restore from.
ForceNew: the source cannot be changed after creation.
Mutually exclusive with source_vol_id and image_id.

### spec.sourceVolId

`string`

source_vol_id is the UUID of an existing Cinder volume to clone.
ForceNew: the source cannot be changed after creation.
Mutually exclusive with snapshot_id and image_id.

### spec.imageId

`string | valueFrom`

image_id is the ID of a Glance image to initialize the volume from.
This creates a bootable volume that can be used as a root disk.
ForceNew: the source cannot be changed after creation.
Mutually exclusive with snapshot_id and source_vol_id.
Can reference an OpenStackImage resource's output or be a literal image UUID.

- references: OpenStackImage (`status.outputs.image_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OpenStackImage, name: <that resource's name>, fieldPath: status.outputs.image_id}} -- a bare string does not parse

### spec.metadata

`map<string, string>`

metadata is a map of key-value pairs stored on the volume.
Metadata is visible in the OpenStack API and Horizon dashboard.
Can be used for tagging, billing, or organizational purposes.

### spec.region

`string`

region overrides the region from the provider config for this volume.
If omitted, the region from the OpenStack provider config is used.
ForceNew: changing the region requires recreating the volume.
Example: "RegionOne"

## Validation Rules

- `volume_source.mutual_exclusion`: at most one of snapshot_id, source_vol_id, or image_id may be set -- they are mutually exclusive volume sources

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackVolume, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.volume_id` | `string` | volume_id is the unique identifier (UUID) of the volume in OpenStack. This is the primary output used as a foreign key by OpenStackVolumeAttach. |
| `status.outputs.name` | `string` | name is the name of the volume (derived from metadata.name). |
| `status.outputs.size` | `int32` | size is the actual provisioned size of the volume in gigabytes. |
| `status.outputs.volume_type` | `string` | volume_type is the Cinder volume type (backend storage class) of the volume. Computed by Cinder if not explicitly specified in the spec. |
| `status.outputs.availability_zone` | `string` | availability_zone is the availability zone where the volume was created. Computed by Cinder if not explicitly specified in the spec. |
| `status.outputs.region` | `string` | region is the OpenStack region where the volume was created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.imageId` | OpenStackImage | `status.outputs.image_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackVolumeAttach | `spec.volumeId` | `status.outputs.volume_id` |

## See Also

- [Overview](../README.md)
