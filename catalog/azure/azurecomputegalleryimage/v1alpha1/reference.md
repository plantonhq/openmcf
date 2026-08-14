# AzureComputeGalleryImage

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a Gen2
# trusted-launch Linux image with recommended sizing, disk-type
# exclusions, and one snapshot-sourced version replicated to two
# regions (one with a customer-managed-key disk encryption set).
# References are literal values so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureComputeGalleryImage
metadata:
  name: test-gallery-image
  id: test-gallery-image
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  galleryName:
    value: platform.images
  name: ubuntu-base
  region: eastus
  identifier:
    publisher: acme
    offer: ubuntu
    sku: 22-04-lts-gen2
  osType: Linux
  hyperVGeneration: V2
  trustedLaunchEnabled: true
  acceleratedNetworkSupportEnabled: true
  diskTypesNotAllowed:
    - Standard_LRS
  minRecommendedVcpuCount: 2
  maxRecommendedVcpuCount: 16
  minRecommendedMemoryInGb: 4
  maxRecommendedMemoryInGb: 64
  endOfLifeDate: "2030-01-01T00:00:00Z"
  releaseNoteUri: https://acme.example/images/ubuntu-base/releases
  description: Hardened Ubuntu 22.04 Gen2 base image
  versions:
    - name: 1.0.0
      osDiskSnapshotId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Compute/snapshots/ubuntu-base-1-0-0
      replicationMode: Full
      targetRegions:
        - name: eastus
          regionalReplicaCount: 2
          storageAccountType: Standard_ZRS
        - name: westeurope
          regionalReplicaCount: 1
          diskEncryptionSetId:
            value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Compute/diskEncryptionSets/images-des
      tags:
        release: stable
  tags:
    costCenter: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.galleryName` | `string \| valueFrom` | yes |  | AzureComputeGallery (`status.outputs.gallery_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.identifier` | `AzureComputeGalleryImageIdentifier` | yes |  |  |
| `spec.identifier.publisher` | `string` | yes |  |  |
| `spec.identifier.offer` | `string` | yes |  |  |
| `spec.identifier.sku` | `string` | yes |  |  |
| `spec.osType` | `string` | yes |  |  |
| `spec.specialized` | `bool` |  |  |  |
| `spec.architecture` | `string` |  |  |  |
| `spec.hyperVGeneration` | `string` |  |  |  |
| `spec.trustedLaunchSupported` | `bool` |  |  |  |
| `spec.trustedLaunchEnabled` | `bool` |  |  |  |
| `spec.confidentialVmSupported` | `bool` |  |  |  |
| `spec.confidentialVmEnabled` | `bool` |  |  |  |
| `spec.acceleratedNetworkSupportEnabled` | `bool` |  |  |  |
| `spec.hibernationEnabled` | `bool` |  |  |  |
| `spec.diskControllerTypeNvmeEnabled` | `bool` |  |  |  |
| `spec.diskTypesNotAllowed` | `[]string` |  |  |  |
| `spec.endOfLifeDate` | `string` |  |  |  |
| `spec.eula` | `string` |  |  |  |
| `spec.privacyStatementUri` | `string` |  |  |  |
| `spec.releaseNoteUri` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.purchasePlan` | `AzureComputeGalleryImagePurchasePlan` |  |  |  |
| `spec.purchasePlan.name` | `string` | yes |  |  |
| `spec.purchasePlan.publisher` | `string` |  |  |  |
| `spec.purchasePlan.product` | `string` |  |  |  |
| `spec.minRecommendedVcpuCount` | `int32` |  |  |  |
| `spec.maxRecommendedVcpuCount` | `int32` |  |  |  |
| `spec.minRecommendedMemoryInGb` | `int32` |  |  |  |
| `spec.maxRecommendedMemoryInGb` | `int32` |  |  |  |
| `spec.versions` | `[]AzureComputeGalleryImageVersion` |  |  |  |
| `spec.versions[].name` | `string` | yes |  |  |
| `spec.versions[].targetRegions` | `[]AzureComputeGalleryImageVersionTargetRegion` | yes |  |  |
| `spec.versions[].targetRegions[].name` | `string` | yes |  |  |
| `spec.versions[].targetRegions[].regionalReplicaCount` | `int32` | yes |  |  |
| `spec.versions[].targetRegions[].diskEncryptionSetId` | `string \| valueFrom` |  |  | AzureDiskEncryptionSet (`status.outputs.disk_encryption_set_id`) |
| `spec.versions[].targetRegions[].excludeFromLatestEnabled` | `bool` |  |  |  |
| `spec.versions[].targetRegions[].storageAccountType` | `string` |  |  |  |
| `spec.versions[].blobUri` | `string` |  |  |  |
| `spec.versions[].storageAccountId` | `string \| valueFrom` |  |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.versions[].osDiskSnapshotId` | `string \| valueFrom` |  |  | AzureDiskSnapshot (`status.outputs.snapshot_id`) |
| `spec.versions[].managedImageId` | `string \| valueFrom` |  |  |  |
| `spec.versions[].replicationMode` | `string` |  |  |  |
| `spec.versions[].excludeFromLatest` | `bool` |  |  |  |
| `spec.versions[].deletionOfReplicatedLocationsEnabled` | `bool` |  |  |  |
| `spec.versions[].endOfLifeDate` | `string` |  |  |  |
| `spec.versions[].tags` | `map<string, string>` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.galleryName

`string | valueFrom` · required

- references: AzureComputeGallery (`status.outputs.gallery_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureComputeGallery, name: <that resource's name>, fieldPath: status.outputs.gallery_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Image names allow up to 80 letters, numbers, dots, dashes, and underscores
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.identifier

`AzureComputeGalleryImageIdentifier` · required

- rule: {"required":true}

### spec.identifier.publisher

`string` · required

- rule: publisher allows up to 128 letters, numbers, dots, dashes, and underscores, and must not end with a dot
- rule: {"required":true}

### spec.identifier.offer

`string` · required

- rule: offer allows up to 64 letters, numbers, dots, dashes, and underscores, and must not end with a dot
- rule: {"required":true}

### spec.identifier.sku

`string` · required

- rule: sku allows up to 64 letters, numbers, dots, dashes, and underscores, and must not end with a dot
- rule: {"required":true}

### spec.osType

`string` · required

- rule: {"required":true,"string":{"in":["Linux","Windows"]}}

### spec.specialized

`bool`

### spec.architecture

`string`

- rule: {"string":{"in":["","x64","Arm64"]}}

### spec.hyperVGeneration

`string`

- rule: {"string":{"in":["","V1","V2"]}}

### spec.trustedLaunchSupported

`bool`

### spec.trustedLaunchEnabled

`bool`

### spec.confidentialVmSupported

`bool`

### spec.confidentialVmEnabled

`bool`

### spec.acceleratedNetworkSupportEnabled

`bool`

### spec.hibernationEnabled

`bool`

### spec.diskControllerTypeNvmeEnabled

`bool`

### spec.diskTypesNotAllowed

`[]string`

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["Standard_LRS","Premium_LRS"]}}}}

### spec.endOfLifeDate

`string`

- rule: end_of_life_date must be an RFC3339 timestamp, e.g. "2030-01-01T00:00:00Z"

### spec.eula

`string`

### spec.privacyStatementUri

`string`

### spec.releaseNoteUri

`string`

### spec.description

`string`

### spec.purchasePlan

`AzureComputeGalleryImagePurchasePlan`

### spec.purchasePlan.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.purchasePlan.publisher

`string`

### spec.purchasePlan.product

`string`

### spec.minRecommendedVcpuCount

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":80,"gte":1}}

### spec.maxRecommendedVcpuCount

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":80,"gte":1}}

### spec.minRecommendedMemoryInGb

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":640,"gte":1}}

### spec.maxRecommendedMemoryInGb

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":640,"gte":1}}

### spec.versions

`[]AzureComputeGalleryImageVersion`

- rule: exactly one of blob_uri, os_disk_snapshot_id, and managed_image_id must be set -- the version has one source
- rule: blob_uri and storage_account_id are set together -- the blob's storage account carries the read grant
- rule: disk_encryption_set_id cannot be used in target_regions when replication_mode is Shallow

### spec.versions[].name

`string` · required

- rule: Version names are three dot-separated numeric segments (e.g. "1.2.0"), or the literal "latest" or "recent"
- rule: {"required":true}

### spec.versions[].targetRegions

`[]AzureComputeGalleryImageVersionTargetRegion` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.versions[].targetRegions[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.versions[].targetRegions[].regionalReplicaCount

`int32` · required

- rule: {"required":true,"int32":{"gte":1}}

### spec.versions[].targetRegions[].diskEncryptionSetId

`string | valueFrom`

- references: AzureDiskEncryptionSet (`status.outputs.disk_encryption_set_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDiskEncryptionSet, name: <that resource's name>, fieldPath: status.outputs.disk_encryption_set_id}} -- a bare string does not parse

### spec.versions[].targetRegions[].excludeFromLatestEnabled

`bool`

### spec.versions[].targetRegions[].storageAccountType

`string`

- rule: {"string":{"in":["","Premium_LRS","Standard_LRS","Standard_ZRS"]}}

### spec.versions[].blobUri

`string`

- rule: blob_uri must be an http or https URL

### spec.versions[].storageAccountId

`string | valueFrom`

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.versions[].osDiskSnapshotId

`string | valueFrom`

- references: AzureDiskSnapshot (`status.outputs.snapshot_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDiskSnapshot, name: <that resource's name>, fieldPath: status.outputs.snapshot_id}} -- a bare string does not parse

### spec.versions[].managedImageId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.versions[].replicationMode

`string`

- rule: {"string":{"in":["","Full","Shallow"]}}

### spec.versions[].excludeFromLatest

`bool`

### spec.versions[].deletionOfReplicatedLocationsEnabled

`bool`

### spec.versions[].endOfLifeDate

`string`

- rule: end_of_life_date must be an RFC3339 timestamp, e.g. "2030-01-01T00:00:00Z"

### spec.versions[].tags

`map<string, string>`

### spec.tags

`map<string, string>`

## Validation Rules

- `gallery_image_security_flags_exclusive`: at most one of trusted_launch_supported, trusted_launch_enabled, confidential_vm_supported, and confidential_vm_enabled can be set -- they are mutually exclusive security postures
- `gallery_image_vcpu_range_ordered`: max_recommended_vcpu_count must be greater than or equal to min_recommended_vcpu_count
- `gallery_image_memory_range_ordered`: max_recommended_memory_in_gb must be greater than or equal to min_recommended_memory_in_gb
- `gallery_image_version_names_unique`: versions names must be unique -- the name is the version's identity under the image

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureComputeGalleryImage, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.image_id` | `string` |  |
| `status.outputs.image_name` | `string` |  |
| `status.outputs.version_ids` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.galleryName` | AzureComputeGallery | `status.outputs.gallery_name` |
| `spec.versions[].targetRegions[].diskEncryptionSetId` | AzureDiskEncryptionSet | `status.outputs.disk_encryption_set_id` |
| `spec.versions[].storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.versions[].osDiskSnapshotId` | AzureDiskSnapshot | `status.outputs.snapshot_id` |

## See Also

- [Overview](../README.md)
