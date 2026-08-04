# OpenStackImage

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `openstack.planton.dev/v1`

OpenStackImageSpec defines the configuration for an OpenStack Glance image.

A Glance image is a virtual disk template used to boot compute instances or
initialize bootable Cinder volumes. Images contain the operating system,
pre-installed software, and any customizations needed for workloads.

Images can be uploaded from a URL (image_source_url) or created as metadata
entries for images already present in Glance. The most common 80/20 use case
is uploading a cloud image from an HTTP URL.

The image name is derived from metadata.name.

Terraform resource: openstack_images_image_v2
Pulumi resource: openstack.images.Image

## Example

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackImage
metadata:
  name: test-image
spec:
  container_format: bare
  disk_format: qcow2
  image_source_url: https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img
  min_disk_gb: 10
  min_ram_mb: 512
  tags:
    - test
    - hack
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.containerFormat` | `string` | yes |  |  |
| `spec.diskFormat` | `string` | yes |  |  |
| `spec.imageSourceUrl` | `string` |  |  |  |
| `spec.minDiskGb` | `int32` |  |  |  |
| `spec.minRamMb` | `int32` |  |  |  |
| `spec.protected` | `bool` |  | `false` |  |
| `spec.hidden` | `bool` |  | `false` |  |
| `spec.tags` | `[]string` |  |  |  |
| `spec.visibility` | `string` |  | `private` |  |
| `spec.region` | `string` |  |  |  |

## Field Details

### spec.containerFormat

`string` · required

container_format describes the container or envelope format of the image.
Required. ForceNew: cannot be changed after creation.

Common values:
  "bare"   - no container (most common for cloud images)
  "ovf"    - Open Virtualization Format
  "docker" - Docker container image

- rule: {"string":{"minLen":"1","in":["bare","ovf","aki","ari","ami","ova","docker","compressed"]}}

### spec.diskFormat

`string` · required

disk_format describes the format of the disk image data.
Required. ForceNew: cannot be changed after creation.

Common values:
  "qcow2" - QEMU Copy On Write 2 (most common for KVM/libvirt)
  "raw"   - unstructured disk image
  "vmdk"  - VMware disk format
  "iso"   - optical disc image

- rule: {"string":{"minLen":"1","in":["raw","vhd","vhdx","vmdk","vdi","iso","ploop","qcow2","aki","ari","ami"]}}

### spec.imageSourceUrl

`string`

image_source_url is the HTTP or HTTPS URL from which Glance downloads the
image data. This is the primary mechanism for uploading images in
pipeline-based IaC workflows.
ForceNew: the source cannot be changed after creation.

Example: "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"

If omitted, the image is created as a metadata entry only -- the actual
image data must be uploaded separately (e.g., via glance CLI).

### spec.minDiskGb

`int32`

min_disk_gb is the minimum disk size in GB required to boot this image.
Set this to prevent users from booting the image on a flavor with
insufficient disk. Zero means no minimum.

- rule: {"int32":{"gte":0}}

### spec.minRamMb

`int32`

min_ram_mb is the minimum RAM in MB required to boot this image.
Set this to prevent users from booting the image on a flavor with
insufficient memory. Zero means no minimum.

- rule: {"int32":{"gte":0}}

### spec.protected

`bool` · optional (explicit presence)

protected prevents the image from being deleted.
Set to true for production images that should not be accidentally removed.
Default: false

- default: `false`

### spec.hidden

`bool` · optional (explicit presence)

hidden hides the image from default listing queries.
Hidden images are still accessible by UUID but do not appear in
`openstack image list` output.
Default: false

- default: `false`

### spec.tags

`[]string`

tags is a set of tags to associate with the image.
Tags are simple strings used for filtering images in Glance.
Example: ["ubuntu", "22.04", "cloud-init"]

### spec.visibility

`string` · optional (explicit presence)

visibility controls who can see and use this image.
Default: "private"

Values:
  "private"   - only the image owner (project) can see it
  "shared"    - owner can share with specific projects via member API
  "community" - visible to all projects but not in default listings
  "public"    - visible and usable by all projects (requires admin)

- default: `private`
- rule: {"string":{"in":["","public","private","shared","community"]}}

### spec.region

`string`

region overrides the region from the provider config for this image.
If omitted, the region from the OpenStack provider config is used.
Example: "RegionOne"

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OpenStackImage, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.image_id` | `string` | image_id is the unique identifier (UUID) of the image in Glance. This is the primary output used as a foreign key by downstream resources. |
| `status.outputs.name` | `string` | name is the name of the image (derived from metadata.name). |
| `status.outputs.checksum` | `string` | checksum is the MD5 checksum of the image data. Computed by Glance after the image data is uploaded. |
| `status.outputs.size_bytes` | `int64` | size_bytes is the size of the image data in bytes. Computed by Glance after upload. |
| `status.outputs.status` | `string` | status is the current lifecycle state of the image. Common values: "active", "queued", "saving", "killed", "deleted". |
| `status.outputs.file` | `string` | file is the URL path to the image data in the Glance store. Example: "/v2/images/<uuid>/file" |
| `status.outputs.region` | `string` | region is the OpenStack region where the image was created. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| OpenStackVolume | `spec.imageId` | `status.outputs.image_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
