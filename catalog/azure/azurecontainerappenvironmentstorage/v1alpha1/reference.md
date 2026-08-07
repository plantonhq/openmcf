# AzureContainerAppEnvironmentStorage

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureContainerAppEnvironmentStorageSpec** defines the configuration for
registering an Azure Files share as a storage resource on a Container App
Environment (Microsoft.App/managedEnvironments/storages).

Container Apps and Jobs cannot mount file shares directly: the share is
first registered on the environment as a named storage resource, and
workloads then declare AZURE_FILE / NFS_AZURE_FILE volumes that reference
the registration by `storage_name`. One registration can back volumes in
any number of apps and jobs in the environment.

**Two share protocols** (exactly one per registration):

SMB (`account_name` + `access_key`):
  - The common case -- a standard Azure Files share addressed by storage
    account name, authenticated with an account access key
  - Works in both external and VNet-injected environments

NFS (`nfs_server_url`):
  - An NFS Azure Files share (premium FileStorage accounts) addressed by
    the account's file endpoint
  - Requires a VNet-injected environment (NFS traffic never leaves the
    VNet); back it with NFS_AZURE_FILE volumes

**Update semantics**: only the SMB `access_key` is updatable in place
(key rotation); every other field forces the registration to be
destroyed and recreated. Recreating a registration briefly breaks
volume mounts that reference it -- plan rotations accordingly.

**Referenced by**: AzureContainerApp and AzureContainerAppJob volumes
(storage_name).

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerAppEnvironmentStorage
metadata:
  name: test-env-storage
spec:
  container_app_environment_id:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.App/managedEnvironments/test-env
  storage_name: app-data
  share_name:
    value: shared-files
  access_mode: READ_WRITE
  account_name:
    value: teststorageaccount
  access_key:
    value: dGVzdC1hY2Nlc3Mta2V5LXZhbHVl
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.containerAppEnvironmentId` | `string \| valueFrom` | yes |  | AzureContainerAppEnvironment (`status.outputs.environment_id`) |
| `spec.storageName` | `string` | yes |  |  |
| `spec.shareName` | `string \| valueFrom` | yes |  | AzureStorageShare (`status.outputs.share_name`) |
| `spec.accessMode` | `enum` |  |  |  |
| `spec.accountName` | `string \| valueFrom` |  |  | AzureStorageAccount (`status.outputs.storage_account_name`) |
| `spec.accessKey` | `string \| valueFrom` (sensitive) |  |  | AzureStorageAccount (`status.outputs.primary_access_key`) |
| `spec.nfsServerUrl` | `string` |  |  |  |

## Field Details

### spec.containerAppEnvironmentId

`string | valueFrom` · required

The Container App Environment to register the storage on.

**ForceNew**: Changing this destroys and recreates the registration.

- references: AzureContainerAppEnvironment (`status.outputs.environment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureContainerAppEnvironment, name: <that resource's name>, fieldPath: status.outputs.environment_id}} -- a bare string does not parse

### spec.storageName

`string` · required

The name of the storage registration -- the handle app and job
volumes reference in their `storage_name`.
Lowercase alphanumeric characters and hyphens; must start with a
letter and end with an alphanumeric character; no consecutive
hyphens; at most 32 characters.

**ForceNew**: Changing this destroys and recreates the registration.

- rule: storage name must be lowercase alphanumeric characters or hyphens, start with a letter, end with an alphanumeric character, and contain no consecutive hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"32"}}

### spec.shareName

`string | valueFrom` · required

The name of the Azure Files share to register. Can be a literal name
or a reference to an AzureStorageShare output.

**ForceNew**: Changing this destroys and recreates the registration.

- references: AzureStorageShare (`status.outputs.share_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageShare, name: <that resource's name>, fieldPath: status.outputs.share_name}} -- a bare string does not parse

### spec.accessMode

`enum`

How workloads may use the share: READ_ONLY for shared configuration
and reference data, READ_WRITE for working storage.

**ForceNew**: Changing this destroys and recreates the registration.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_container_app_environment_storage_access_mode_unspecified` -- Not specified -- invalid; choose READ_ONLY or READ_WRITE.
- `READ_ONLY` -- Volumes backed by this registration mount read-only.
- `READ_WRITE` -- Volumes backed by this registration mount read-write.

### spec.accountName

`string | valueFrom`

The storage account holding the share -- the SMB path. Can be a
literal account name or a reference to an AzureStorageAccount output.
Pair with `access_key`; mutually exclusive with `nfs_server_url`.

**ForceNew**: Changing this destroys and recreates the registration.

- references: AzureStorageAccount (`status.outputs.storage_account_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_name}} -- a bare string does not parse

### spec.accessKey

`string | valueFrom` · sensitive

The storage account access key authenticating the SMB mount. Pair
with `account_name`. This is the one field that updates in place --
rotate keys without recreating the registration.

- references: AzureStorageAccount (`status.outputs.primary_access_key`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.primary_access_key}} -- a bare string does not parse

### spec.nfsServerUrl

`string`

The NFS server URL for an NFS Azure Files share -- the NFS path.
Format: {account}.file.core.windows.net
Mutually exclusive with `account_name`/`access_key`; requires a
VNet-injected environment.

**ForceNew**: Changing this destroys and recreates the registration.

## Validation Rules

- `environment_storage_smb_xor_nfs`: register the share over exactly one protocol: SMB (account_name + access_key together) or NFS (nfs_server_url alone)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureContainerAppEnvironmentStorage, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.storage_id` | `string` | The Azure Resource Manager ID of the storage registration. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.App/managedEnvironments/{env}/storages/{name} |
| `status.outputs.storage_name` | `string` | The name of the storage registration -- what app and job volumes reference in their storage_name field. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.containerAppEnvironmentId` | AzureContainerAppEnvironment | `status.outputs.environment_id` |
| `spec.shareName` | AzureStorageShare | `status.outputs.share_name` |
| `spec.accountName` | AzureStorageAccount | `status.outputs.storage_account_name` |
| `spec.accessKey` | AzureStorageAccount | `status.outputs.primary_access_key` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureContainerApp | `spec.volumes[].storageName` | `status.outputs.storage_name` |
| AzureContainerAppJob | `spec.volumes[].storageName` | `status.outputs.storage_name` |

## See Also

- [Overview](../README.md)
