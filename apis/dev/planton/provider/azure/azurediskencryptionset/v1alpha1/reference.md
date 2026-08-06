# AzureDiskEncryptionSet

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureDiskEncryptionSetSpec** defines an Azure Disk Encryption Set: the
resource that binds a customer-managed key (CMK) to managed disks,
snapshots, and images so their data is encrypted at rest with a key you
control in Key Vault rather than a platform-managed key. Managed disks,
VMs, and VM scale sets reference a disk encryption set by ARM ID
(disk_encryption_set_id) to opt their OS and data disks into CMK
encryption.

The set carries a managed identity that unwraps the key at runtime: that
identity must be granted crypto access ("Key Vault Crypto Service
Encryption User", or get/wrapKey/unwrapKey) on the key before disks can
use the set. The referenced Key Vault must have purge protection enabled
-- Azure requires it for any vault backing disk encryption.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureDiskEncryptionSet
metadata:
  name: test-disk-encryption-set
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: test-des
  keyVaultKeyId:
    value: https://test-vault.vault.azure.net/keys/cmk-key
  autoKeyRotationEnabled: true
  identity:
    type: SYSTEM_ASSIGNED
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.keyVaultKeyId` | `string \| valueFrom` | yes |  | AzureKeyVaultKey (`status.outputs.versionless_id`) |
| `spec.autoKeyRotationEnabled` | `bool` |  |  |  |
| `spec.encryptionType` | `enum` |  |  |  |
| `spec.federatedClientId` | `string` |  |  |  |
| `spec.identity` | `AzureDiskEncryptionSetIdentity` | yes |  |  |
| `spec.identity.type` | `enum` | yes |  |  |
| `spec.identity.identityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region the set is created in. Must match the region of the
disks that reference it. Fixed at creation.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the set is created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the set, unique within the resource group. 1-80 characters
(letters, numbers, underscores, hyphens, periods; cannot start with an
underscore or end with a hyphen or period). Fixed at creation.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.keyVaultKeyId

`string | valueFrom` · required

The Key Vault key that encrypts disks bound to this set. Defaults to
referencing an AzureKeyVaultKey's versionless_id output, which pairs
with auto_key_rotation_enabled = true (the recommended posture): the
set follows the key's current version automatically as it rotates.

The key's versioning must match the rotation setting -- a VERSIONLESS
key id (no version suffix) when auto_key_rotation_enabled is true, a
VERSIONED key id when it is false. Azure enforces this at apply; both
engines pass the referenced id through and the provider validates it.
(A field-level CEL cannot express this because it cannot dereference a
StringValueOrRef's resolved value; keep the reference and the rotation
flag consistent -- point at versionless_id for rotation, key_id for a
pinned version.)

Keys backed by Azure Key Vault Managed HSM (a separate premium
service) are not modeled -- this field carries standard Key Vault key
URLs only.

- references: AzureKeyVaultKey (`status.outputs.versionless_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultKey, name: <that resource's name>, fieldPath: status.outputs.versionless_id}} -- a bare string does not parse

### spec.autoKeyRotationEnabled

`bool` · optional (explicit presence)

Whether the set automatically follows the key's latest version as it
rotates. True (the recommended posture) requires a VERSIONLESS key id;
false pins the set to one VERSIONED key id and requires manual rotation.
Azure's default is false; leave unset to keep that, or set true and
reference the key's versionless_id.

### spec.encryptionType

`enum`

What the set encrypts. Unspecified applies Azure's default
(ENCRYPTION_AT_REST_WITH_CUSTOMER_KEY -- disks encrypted with your CMK
only). Fixed at creation.

Allowed values (use exactly as shown):

- `azure_disk_encryption_set_encryption_type_unspecified` -- Not specified: Azure's default (EncryptionAtRestWithCustomerKey).
- `ENCRYPTION_AT_REST_WITH_CUSTOMER_KEY` -- Disks encrypted at rest with the customer-managed key only.
- `ENCRYPTION_AT_REST_WITH_PLATFORM_AND_CUSTOMER_KEYS` -- Double encryption: a platform-managed key plus the customer-managed key, for defense in depth.
- `CONFIDENTIAL_VM_ENCRYPTED_WITH_CUSTOMER_KEY` -- Confidential-VM guest-state encryption with the customer-managed key (for confidential VMs).

### spec.federatedClientId

`string`

The client ID of a multi-tenant application used to access the key in a
Key Vault in a DIFFERENT Entra tenant (cross-tenant CMK). Leave empty
for the same-tenant case (the norm). Must be a UUID.

- rule: federated_client_id must be a UUID

### spec.identity

`AzureDiskEncryptionSetIdentity` · required

The managed identity the set uses to unwrap the key. Required -- a set
cannot read its key without an identity. SYSTEM_ASSIGNED is the common
choice (Azure creates and rotates it); USER_ASSIGNED brings an identity
you manage and can grant vault access to BEFORE the set exists (avoiding
the system-assigned chicken-and-egg). Whichever you use, its principal
must be granted crypto access on the key.

- rule: {"required":true}
- rule: identity_ids is required for USER_ASSIGNED and SYSTEM_AND_USER_ASSIGNED and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

Identity flavor. SYSTEM_ASSIGNED is created and rotated by Azure with
the set; USER_ASSIGNED brings identities you manage (letting you grant
vault access before the set exists); SYSTEM_AND_USER_ASSIGNED carries
both. Changing from a user-assigned flavor to SYSTEM_ASSIGNED replaces
the set.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_disk_encryption_set_identity_type_unspecified` -- Not specified -- invalid; a type must be chosen.
- `SYSTEM_ASSIGNED` -- Azure creates and rotates a system-assigned identity with the set.
- `USER_ASSIGNED` -- One or more user-assigned identities you manage.
- `SYSTEM_AND_USER_ASSIGNED` -- Both a system-assigned and one or more user-assigned identities.

### spec.identity.identityIds

`[]string | valueFrom`

For USER_ASSIGNED / SYSTEM_AND_USER_ASSIGNED: the user-assigned
identities attached to the set, by ARM ID. Reference
AzureUserAssignedIdentity resources so the Key Vault crypto grant can be
composed before the set is created.

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

Free-form tags applied to the set, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag with
the same key wins. Updatable in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureDiskEncryptionSet, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.disk_encryption_set_id` | `string` | The Azure Resource Manager ID of the set. This is the composition seam: AzureManagedDisk, AzureVirtualMachine, and AzureVirtualMachineScaleSet reference it (disk_encryption_set_id) to encrypt their disks with the customer-managed key. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Compute/diskEncryptionSets/{name} |
| `status.outputs.disk_encryption_set_name` | `string` | The name of the set resource. |
| `status.outputs.identity_principal_id` | `string` | The principal (object) ID of the set's system-assigned identity, when one is used. This is the principal to grant Key Vault crypto access (get/wrapKey/unwrapKey, or the "Key Vault Crypto Service Encryption User" role) so the set can unwrap the key. |
| `status.outputs.identity_tenant_id` | `string` | The Entra tenant ID of the set's system-assigned identity. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.keyVaultKeyId` | AzureKeyVaultKey | `status.outputs.versionless_id` |
| `spec.identity.identityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureAksCluster | `spec.diskEncryptionSetId` | `status.outputs.disk_encryption_set_id` |
| AzureManagedDisk | `spec.diskEncryptionSetId` | `status.outputs.disk_encryption_set_id` |
| AzureManagedDisk | `spec.secureVmDiskEncryptionSetId` | `status.outputs.disk_encryption_set_id` |
| AzureVirtualMachine | `spec.osDisk.diskEncryptionSetId` | `status.outputs.disk_encryption_set_id` |
| AzureVirtualMachine | `spec.osDisk.secureVmDiskEncryptionSetId` | `status.outputs.disk_encryption_set_id` |
| AzureVirtualMachineScaleSet | `spec.osDisk.diskEncryptionSetId` | `status.outputs.disk_encryption_set_id` |
| AzureVirtualMachineScaleSet | `spec.osDisk.secureVmDiskEncryptionSetId` | `status.outputs.disk_encryption_set_id` |
| AzureVirtualMachineScaleSet | `spec.dataDisks[].diskEncryptionSetId` | `status.outputs.disk_encryption_set_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
