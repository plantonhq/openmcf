# CivoVolume

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `civo.planton.dev/v1alpha1`

CivoVolumeSpec defines the specification required to create a Civo block storage volume.
A block storage volume provides expandable storage that can be attached to Civo instances.
This specification focuses on essential parameters for volume creation, adhering to the 80/20 principle.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.volumeName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.sizeGib` | `uint32` | yes |  |  |
| `spec.filesystemType` | `enum` |  |  |  |
| `spec.snapshotId` | `string` |  |  |  |
| `spec.tags` | `[]string` |  |  |  |

## Field Details

### spec.volumeName

`string` · required

The name of the volume. Must be lowercase letters, numbers, and hyphens only,
starting with a letter and ending with a letter or number. Maximum 64 characters.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"64","pattern":"^[a-z]([a-z0-9-]*[a-z0-9])?$"}}

### spec.region

`enum` · required

The Civo region where the volume will be created.
Must match the region of any instance that will attach to this volume.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `civo_region_unspecified` -- 0: default / unspecified region
- `lon1` -- london 1
- `lon2` -- london 2
- `fra1` -- frankfurt 1
- `nyc1` -- new york 1
- `phx1` -- phoenix 1
- `mum1` -- mumbai 1

### spec.sizeGib

`uint32` · required

The size of the volume in GiB.
Constraints: between 1 and 16000 (inclusive).

- rule: {"required":true,"uint32":{"lte":16000,"gte":1}}

### spec.filesystemType

`enum`

The initial filesystem to format the volume with.
Allowed values: ext4, xfs, or none (no pre-formatting). Default is none.

Allowed values (use exactly as shown):

- `unformatted`
- `ext4`
- `xfs`

### spec.snapshotId

`string`

An optional snapshot ID or reference to create this volume from.
If provided, the new volume will be created from the given snapshot (inheriting its region and minimum size).

### spec.tags

`[]string`

A list of tags to apply to the volume.
Tags must be unique and consist of letters, numbers, colons, dashes, or underscores.

- rule: {"repeated":{"unique":true,"items":{"string":{"maxLen":"64","pattern":"^[A-Za-z0-9:_-]+$"}}}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CivoVolume, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.volume_id` | `string` | The unique identifier (ID) of the created Civo volume. |
| `status.outputs.attached_instance_id` | `string` | The ID of the Civo instance the volume is attached to (if any). |
| `status.outputs.device_path` | `string` | The device path of the volume on the attached instance (if any). |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| CivoComputeInstance | `spec.volumeIds` | `status.outputs.volume_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
