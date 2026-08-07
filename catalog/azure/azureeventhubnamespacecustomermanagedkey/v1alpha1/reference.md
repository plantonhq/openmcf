# AzureEventHubNamespaceCustomerManagedKey

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureEventHubNamespaceCustomerManagedKeySpec** configures
customer-managed-key (BYOK) encryption on an Event Hubs namespace:
event data at rest is encrypted with YOUR Key Vault keys instead of
Microsoft-managed keys.

Azure models CMK as a configuration applied ONTO an existing namespace
(not a create-time property), and this kind mirrors that grain -- for a
real reason: encrypting with the namespace's own system-assigned
identity is only possible as a second step (create the namespace with
its identity -> grant that identity wrap/unwrap on the vault, e.g. a
"Key Vault Crypto Service Encryption User" AzureRoleAssignment -> apply
this kind). A folded create-time block could never express that
sequence.

**Platform contract enforced by Azure at apply time**: CMK requires
single-tenant capacity -- a namespace placed on a dedicated cluster
(dedicated_cluster_id) or on the PREMIUM tier. Multi-tenant
BASIC/STANDARD namespaces share hardware and cannot take tenant keys;
Azure rejects the encryption patch on them.

**Identity contract**: with user_assigned_identity_id set, that
identity unwraps the keys and MUST already be attached via the
namespace's identity block (with vault access granted before this kind
applies). Without it, the namespace's system-assigned identity is used.

**Add-only lifecycle (Azure's own contract)**: once CMK is enabled it
can never be removed -- Azure has no decrypt-back path. Deleting this
resource intentionally changes NOTHING on the namespace (the delete is
a no-op); returning to Microsoft-managed keys requires replacing the
namespace itself. Key ROTATION is routine, though: with versionless key
references, vault rotation propagates automatically.

**ForceNew fields**: `eventhub_namespace_id`,
`infrastructure_encryption_enabled`.

## Example

```yaml
# Offline-plan manifest: a CMK configuration exercising the key list,
# the infrastructure-encryption layer, and the user-assigned unwrapping
# identity.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventHubNamespaceCustomerManagedKey
metadata:
  name: test-eventhub-cmk
spec:
  eventhubNamespaceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.EventHub/namespaces/my-ehns
  keyVaultKeyIds:
    - value: https://my-vault.vault.azure.net/keys/my-key
  infrastructureEncryptionEnabled: true
  userAssignedIdentityId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/my-uai
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.eventhubNamespaceId` | `string \| valueFrom` | yes |  | AzureEventHubNamespace (`status.outputs.namespace_id`) |
| `spec.keyVaultKeyIds` | `[]string \| valueFrom` | yes |  | AzureKeyVaultKey (`status.outputs.versionless_id`) |
| `spec.infrastructureEncryptionEnabled` | `bool` |  |  |  |
| `spec.userAssignedIdentityId` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |

## Field Details

### spec.eventhubNamespaceId

`string | valueFrom` · required

The namespace to encrypt, by ARM ID. References an
AzureEventHubNamespace's namespace_id output. The namespace must have
single-tenant capacity (a dedicated cluster or PREMIUM) and already
carry the unwrapping identity. Fixed at creation.

- references: AzureEventHubNamespace (`status.outputs.namespace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHubNamespace, name: <that resource's name>, fieldPath: status.outputs.namespace_id}} -- a bare string does not parse

### spec.keyVaultKeyIds

`[]string | valueFrom` · required

The Key Vault keys that encrypt the namespace's data, by data-plane
key ID (1-10 keys; Azure applies them across the namespace). Each
defaults to referencing an AzureKeyVaultKey's versionless_id output
so vault-side rotation propagates automatically; pin versioned IDs
only when a compliance regime demands immutable key versions. The
keys' vault must have purge protection enabled.

- references: AzureKeyVaultKey (`status.outputs.versionless_id`)
- rule: {"repeated":{"minItems":"1","maxItems":"10"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureKeyVaultKey, name: <that resource's name>, fieldPath: status.outputs.versionless_id}} -- a bare string does not parse

### spec.infrastructureEncryptionEnabled

`bool` · optional (explicit presence)

Whether Azure applies a second layer of encryption (infrastructure
encryption) beneath the customer-managed keys.

**ForceNew**: fixed the moment CMK is first configured.

### spec.userAssignedIdentityId

`string | valueFrom`

The user-assigned identity Azure uses to unwrap the keys, by ARM ID.
References an AzureUserAssignedIdentity's identity_id output. Must
already be attached via the namespace's identity block, with
wrap/unwrap access on the keys' vault. Unset uses the namespace's
system-assigned identity instead (grant IT the vault access, using
the namespace's identity_principal_id output).

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventHubNamespaceCustomerManagedKey, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.customer_managed_key_id` | `string` | The provider's identity for the CMK configuration -- the namespace's ARM ID (the configuration is a property of the namespace; it has no ARM object of its own). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.eventhubNamespaceId` | AzureEventHubNamespace | `status.outputs.namespace_id` |
| `spec.keyVaultKeyIds` | AzureKeyVaultKey | `status.outputs.versionless_id` |
| `spec.userAssignedIdentityId` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## See Also

- [Overview](../README.md)
