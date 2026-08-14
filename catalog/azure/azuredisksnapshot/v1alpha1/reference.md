# AzureDiskSnapshot

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: an incremental
# Copy-mode snapshot of a managed disk with a private network posture
# and legacy ADE encryption settings carried from the source.
# References are literal values so the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureDiskSnapshot
metadata:
  name: test-disk-snapshot
  id: test-disk-snapshot
  org: test-org
  env: test
spec:
  resourceGroup:
    value: test-rg
  name: app-disk-snap
  region: eastus
  createOption: Copy
  sourceResourceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Compute/disks/app-disk
  incrementalEnabled: true
  networkAccessPolicy: AllowPrivate
  diskAccessId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Compute/diskAccesses/app-disk-access
  publicNetworkAccessEnabled: false
  encryptionSettings:
    diskEncryptionKey:
      secretUrl: https://app-vault.vault.azure.net/secrets/app-disk-dek/0000000000000000
      sourceVaultId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.KeyVault/vaults/app-vault
    keyEncryptionKey:
      keyUrl: https://app-vault.vault.azure.net/keys/app-disk-kek/0000000000000000
      sourceVaultId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.KeyVault/vaults/app-vault
  tags:
    backupChain: app-disk
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.createOption` | `string` | yes |  |  |
| `spec.sourceResourceId` | `string \| valueFrom` |  |  | AzureManagedDisk (`status.outputs.disk_id`) |
| `spec.sourceUri` | `string` |  |  |  |
| `spec.storageAccountId` | `string \| valueFrom` |  |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.incrementalEnabled` | `bool` |  |  |  |
| `spec.diskSizeGb` | `int32` |  |  |  |
| `spec.networkAccessPolicy` | `string` |  |  |  |
| `spec.diskAccessId` | `string \| valueFrom` |  |  |  |
| `spec.publicNetworkAccessEnabled` | `bool` |  |  |  |
| `spec.encryptionSettings` | `AzureDiskSnapshotEncryptionSettings` |  |  |  |
| `spec.encryptionSettings.diskEncryptionKey` | `AzureDiskSnapshotDiskEncryptionKey` | yes |  |  |
| `spec.encryptionSettings.diskEncryptionKey.secretUrl` | `string` | yes |  |  |
| `spec.encryptionSettings.diskEncryptionKey.sourceVaultId` | `string \| valueFrom` | yes |  | AzureKeyVault (`status.outputs.key_vault_id`) |
| `spec.encryptionSettings.keyEncryptionKey` | `AzureDiskSnapshotKeyEncryptionKey` |  |  |  |
| `spec.encryptionSettings.keyEncryptionKey.keyUrl` | `string` | yes |  |  |
| `spec.encryptionSettings.keyEncryptionKey.sourceVaultId` | `string \| valueFrom` | yes |  | AzureKeyVault (`status.outputs.key_vault_id`) |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Snapshot names allow up to 80 letters, numbers, dashes, and underscores
- rule: {"required":true}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.createOption

`string` · required

- rule: {"required":true,"string":{"in":["Copy","Import"]}}

### spec.sourceResourceId

`string | valueFrom`

- references: AzureManagedDisk (`status.outputs.disk_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureManagedDisk, name: <that resource's name>, fieldPath: status.outputs.disk_id}} -- a bare string does not parse

### spec.sourceUri

`string`

### spec.storageAccountId

`string | valueFrom`

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.incrementalEnabled

`bool`

### spec.diskSizeGb

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1}}

### spec.networkAccessPolicy

`string`

- rule: {"string":{"in":["","AllowAll","AllowPrivate","DenyAll"]}}

### spec.diskAccessId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.publicNetworkAccessEnabled

`bool` · optional (explicit presence)

### spec.encryptionSettings

`AzureDiskSnapshotEncryptionSettings`

### spec.encryptionSettings.diskEncryptionKey

`AzureDiskSnapshotDiskEncryptionKey` · required

- rule: {"required":true}

### spec.encryptionSettings.diskEncryptionKey.secretUrl

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.encryptionSettings.diskEncryptionKey.sourceVaultId

`string | valueFrom` · required

- references: AzureKeyVault (`status.outputs.key_vault_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVault, name: <that resource's name>, fieldPath: status.outputs.key_vault_id}} -- a bare string does not parse

### spec.encryptionSettings.keyEncryptionKey

`AzureDiskSnapshotKeyEncryptionKey`

### spec.encryptionSettings.keyEncryptionKey.keyUrl

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.encryptionSettings.keyEncryptionKey.sourceVaultId

`string | valueFrom` · required

- references: AzureKeyVault (`status.outputs.key_vault_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVault, name: <that resource's name>, fieldPath: status.outputs.key_vault_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureDiskSnapshot, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.snapshot_id` | `string` |  |
| `status.outputs.snapshot_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.sourceResourceId` | AzureManagedDisk | `status.outputs.disk_id` |
| `spec.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.encryptionSettings.diskEncryptionKey.sourceVaultId` | AzureKeyVault | `status.outputs.key_vault_id` |
| `spec.encryptionSettings.keyEncryptionKey.sourceVaultId` | AzureKeyVault | `status.outputs.key_vault_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureComputeGalleryImage | `spec.versions[].osDiskSnapshotId` | `status.outputs.snapshot_id` |

## See Also

- [Overview](../README.md)
