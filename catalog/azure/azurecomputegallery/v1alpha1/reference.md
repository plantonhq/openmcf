# AzureComputeGallery

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a
# Community-shared gallery exercising the whole sharing tree.
# References are literal values so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureComputeGallery
metadata:
  name: test-compute-gallery
  id: test-compute-gallery
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: platform.images
  region: eastus
  description: The platform team's approved golden images
  sharing:
    permission: Community
    communityGallery:
      eula: https://example.com/image-eula
      prefix: acmeimages
      publisherEmail: images@acme.example
      publisherUri: https://acme.example
  tags:
    costCenter: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.sharing` | `AzureComputeGallerySharing` |  |  |  |
| `spec.sharing.permission` | `string` | yes |  |  |
| `spec.sharing.communityGallery` | `AzureComputeGalleryCommunitySharing` |  |  |  |
| `spec.sharing.communityGallery.eula` | `string` | yes |  |  |
| `spec.sharing.communityGallery.prefix` | `string` | yes |  |  |
| `spec.sharing.communityGallery.publisherEmail` | `string` | yes |  |  |
| `spec.sharing.communityGallery.publisherUri` | `string` | yes |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Gallery names allow up to 80 letters, numbers, dots, and underscores (no dashes)
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.description

`string`

### spec.sharing

`AzureComputeGallerySharing`

- rule: community_gallery must be set when permission is Community

### spec.sharing.permission

`string` · required

- rule: {"required":true,"string":{"in":["Community","Groups","Private"]}}

### spec.sharing.communityGallery

`AzureComputeGalleryCommunitySharing`

### spec.sharing.communityGallery.eula

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.sharing.communityGallery.prefix

`string` · required

- rule: The community prefix must be 5-16 letters and numbers
- rule: {"required":true}

### spec.sharing.communityGallery.publisherEmail

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.sharing.communityGallery.publisherUri

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.tags

`map<string, string>`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureComputeGallery, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.gallery_id` | `string` |  |
| `status.outputs.gallery_name` | `string` |  |
| `status.outputs.unique_name` | `string` |  |
| `status.outputs.community_gallery_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureComputeGalleryImage | `spec.galleryName` | `status.outputs.gallery_name` |

## See Also

- [Overview](../README.md)
