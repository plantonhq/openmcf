# OpenStack Image

Deploys a Glance image on OpenStack by importing a virtual disk from a URL or registering image metadata. Images serve as boot templates for compute instances and bootable volumes, and are referenced by Magnum cluster templates for container cluster node provisioning. Supports configurable container and disk formats, minimum hardware requirements, visibility controls, and tagging.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Glance Image** -- a virtual disk image registered in the OpenStack Image service with the specified container format, disk format, and optional source URL for data import
- **Image Data** -- created only when `imageSourceUrl` is provided; Glance downloads the image data from the URL and stores it in its configured backend (Ceph, Swift, or filesystem)
- **OpenStack Tags** -- user-defined tags applied to the image for filtering in Glance API queries and the Horizon dashboard

## Before You Deploy

### OpenStack Account

- **An HTTP/HTTPS URL** (optional) -- if importing an image from a URL, the Glance service must be able to reach the download endpoint. If omitted, the image is created as a metadata entry only and the data must be uploaded separately.
- **Sufficient Glance storage** -- the backend storage (Ceph, Swift, or filesystem) must have capacity for the image file. Cloud images range from 200 MB to several GB depending on the distribution and pre-installed software.

## Deploy

### Console

Open the deployment store, find **OpenStack Image**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Cloud Image from URL** preset in the [Presets](#presets) tab to import a standard Linux cloud image.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackImage
metadata:
  name: ubuntu-22-04
  org: acme-corp
  env: prod
spec:
  containerFormat: bare
  diskFormat: qcow2
  imageSourceUrl: "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
```

```shell
planton apply -f image.yaml
```

This imports an Ubuntu 22.04 cloud image in qcow2 format with private visibility. No minimum disk or RAM requirements are set, and the image is not protected from deletion.

## Key Configuration

These are the most important decisions when configuring a Glance image. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Container and disk format** -- Most Linux cloud images use `containerFormat: bare` (no envelope) and `diskFormat: qcow2` (QEMU copy-on-write). Use `raw` for maximum I/O performance at the cost of larger file size. Use `vmdk` for VMware-originated images. Both fields are immutable after creation.

**Image source** -- Set `imageSourceUrl` to an HTTP/HTTPS URL for pipeline-based imports. Glance downloads the image asynchronously. If omitted, the image is created as a metadata-only entry -- useful when uploading via the Glance CLI or when the image data already exists in the backend.

**Minimum requirements** -- Set `minDiskGb` and `minRamMb` to prevent users from booting the image on undersized flavors. Zero means no minimum. Useful for images with known hardware requirements (e.g., databases, GPU workloads).

**Visibility and protection** -- Default `visibility` is `private` (only the owning project). Change to `shared` to grant access to specific projects, or `public` for all projects (requires admin). Set `protected: true` to prevent accidental deletion of production images.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `image_id` | UUID of the image in Glance | OpenStackInstance `imageName`, container cluster templates, bootable volume creation |
| `name` | Image name | Labels, audit trails |
| `checksum` | MD5 checksum of the image data | Image integrity verification |
| `size_bytes` | Size of the image data in bytes | Capacity planning, download estimation |
| `status` | Lifecycle state (active, queued, saving) | Health checks, provisioning validation |
| `file` | URL path to image data in Glance | Direct Glance API access |
| `region` | OpenStack region of the image | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Cloud image from URL** -- Imports a standard Linux cloud image (Ubuntu, CentOS, Fedora CoreOS) from an HTTP URL in bare/qcow2 format with private visibility. Covers the most common image import workflow. Start from the **Cloud Image from URL** preset.

## Works With

This component operates independently and does not reference other deployment components.