# HetznerCloudVolume

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1`

HetznerCloudVolumeSpec defines the specification for a Hetzner Cloud block
storage volume.

A volume is a network-attached block storage device that can be attached to
exactly one server at a time. Volumes persist independently of any server --
detaching or deleting the server does not affect the volume's data. This
makes them suitable for databases, application state, and any data that must
survive server replacement.

Volumes are always created in a specific Hetzner Cloud location and can only
be attached to servers in the same location. Size can be increased after
creation but cannot be decreased (a provider-enforced constraint).

An optional server attachment creates an `hcloud_volume_attachment` resource
that connects the volume to a server. The attachment is managed separately
from the volume itself, allowing the volume to be detached and reattached
to different servers without recreating the volume.

Bundled provider resources:
  - hcloud_volume:             The block storage volume itself.
  - hcloud_volume_attachment:  Optional attachment to a server. Only created
                               when server_id is set.

Fields not exposed in this spec (hardcoded or derived in IaC modules):
  - name:   Derived from metadata.name.
  - labels: Derived from metadata (CG01 pattern). Standard labels take
            precedence over user-specified metadata.labels.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudVolume
metadata:
  name: hetznercloudvolume-demo
spec:
  size: 100
  location: fsn1
  format: ext4
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.size` | `int32` |  |  |  |
| `spec.location` | `string` | yes |  |  |
| `spec.format` | `enum` |  |  |  |
| `spec.serverId` | `string \| valueFrom` |  |  | HetznerCloudServer (`status.outputs.server_id`) |
| `spec.automount` | `bool` |  |  |  |
| `spec.deleteProtection` | `bool` |  |  |  |

## Field Details

### spec.size

`int32`

Volume size in GB. Minimum is 10 GB, maximum is 10240 GB (10 TB).

Size can be increased after creation (the provider resizes the underlying
block device), but it can never be decreased -- the Hetzner Cloud API
rejects size reductions.

- rule: {"int32":{"lte":10240,"gte":10}}

### spec.location

`string` · required

Hetzner Cloud location where the volume is stored (e.g., "fsn1", "nbg1",
"hel1", "ash", "hil", "sin"). Determines the physical datacenter for
the block storage device.

The volume can only be attached to servers in the same location.

Changing this value forces replacement of the volume (data loss).

- rule: {"string":{"minLen":"1"}}

### spec.format

`enum`

Filesystem format applied when the volume is first created.

If unspecified (format_unspecified / 0), the volume is created raw without
a filesystem. A raw volume must be formatted manually from the server
before it can be mounted.

This is a create-time-only setting: the provider does not read it back
after creation. Changing it in the spec after the initial apply has no
effect on the existing volume.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `format_unspecified` -- Raw volume with no filesystem. Must be formatted manually from the server before mounting.
- `ext4` -- ext4 filesystem. The most common Linux filesystem; recommended for general-purpose workloads.
- `xfs` -- XFS filesystem. High-performance filesystem suited for large files and high-throughput workloads.

### spec.serverId

`string | valueFrom`

Server to attach this volume to. Optional.

Accepts a literal Hetzner Cloud server ID (as a string) or a reference
to a HetznerCloudServer resource's output via valueFrom. The referenced
server must be in the same location as the volume.

If omitted, the volume is created unattached (available). Attachment can
be added later by updating the spec.

When set, an hcloud_volume_attachment resource is created to manage the
attachment lifecycle separately from the volume itself.

Example (literal):
  serverId:
    value: "12345678"

Example (reference):
  serverId:
    valueFrom:
      kind: HetznerCloudServer
      name: my-server
      fieldPath: status.outputs.server_id

- references: HetznerCloudServer (`status.outputs.server_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: HetznerCloudServer, name: <that resource's name>, fieldPath: status.outputs.server_id}} -- a bare string does not parse

### spec.automount

`bool`

Automatically mount the volume on the server after attaching.

Only meaningful when server_id is set. When true, Hetzner Cloud
automatically mounts the volume at a system-assigned mount point
after attachment.

This is a create-time-only setting: the provider does not read it back
after creation. Changing it in the spec after the initial apply has no
effect.

### spec.deleteProtection

`bool`

Prevent accidental deletion of the volume via the Hetzner Cloud API.
When enabled, the volume cannot be deleted until protection is removed.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudVolume, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.volume_id` | `string` | The Hetzner Cloud numeric ID of the created volume (as a string). Can be referenced by other components via StringValueOrRef. |
| `status.outputs.linux_device` | `string` | The Linux device path for the volume on the attached server (e.g., "/dev/disk/by-id/scsi-0HC_Volume_12345678"). Empty if the volume is not attached to a server. Use this path to mount the volume from within the server's OS. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.serverId` | HetznerCloudServer | `status.outputs.server_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
