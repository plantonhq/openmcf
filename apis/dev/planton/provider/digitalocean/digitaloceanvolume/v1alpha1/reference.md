# DigitalOceanVolume

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1alpha1`

DigitalOceanVolumeSpec defines the specification required to create a DigitalOcean block storage volume.
A block storage volume provides expandable storage that can be attached to Droplets.
This specification focuses on essential parameters for volume creation, adhering to the 80/20 principle.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.volumeName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
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

### spec.description

`string`

An optional description for the volume.
Constraints: Maximum 100 characters.

- rule: {"string":{"maxLen":"100"}}

### spec.region

`enum` · required

The DigitalOcean region where the volume will be created.
Must match the region of any Droplet that will attach to this volume.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_region_unspecified` -- 0: default / unspecified region
- `nyc3` -- new york 3
- `sfo3` -- san francisco 3
- `fra1` -- frankfurt 1
- `sgp1` -- singapore 1
- `lon1` -- london 1
- `tor1` -- toronto 1
- `blr1` -- bangalore 1
- `ams3` -- amsterdam 3

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

An optional snapshot ID or reference to a volume snapshot to create this volume from.
If provided, the new volume will be created from the given snapshot (inheriting its region and minimum size).

### spec.tags

`[]string`

A list of tags to apply to the volume.
Tags must be unique and consist of letters, numbers, colons, dashes, or underscores.

- rule: {"repeated":{"unique":true,"items":{"string":{"maxLen":"64","pattern":"^[A-Za-z0-9:_-]+$"}}}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanVolume, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.volume_id` | `string` | The unique identifier (UUID) of the created DigitalOcean volume. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| DigitalOceanDroplet | `spec.volumeIds` | `status.outputs.volume_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
