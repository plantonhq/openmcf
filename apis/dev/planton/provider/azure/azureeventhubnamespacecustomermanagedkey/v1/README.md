# AzureEventHubNamespaceCustomerManagedKey

Customer-managed-key (BYOK) encryption applied ONTO an existing Event
Hubs namespace: event data at rest is encrypted with YOUR Key Vault
keys instead of Microsoft-managed keys. Azure models CMK as a
configuration applied after namespace creation -- not a create-time
property -- and this kind mirrors that grain for a real reason:
encrypting with the namespace's own system-assigned identity is only
possible as a second step. Create the namespace with its identity,
grant that identity wrap/unwrap on the vault (a "Key Vault Crypto
Service Encryption User" AzureRoleAssignment), then apply this kind. A
folded create-time block could never express that sequence.

## When to Use

Use AzureEventHubNamespaceCustomerManagedKey when you need:

- **Tenant-controlled encryption keys** -- compliance regimes requiring
  your keys, your vault, your revocation power over event data at rest
- **Rotation-follows-latest** -- versionless key references make
  vault-side rotation propagate automatically, with no manifest change
- **Double encryption** -- `infrastructure_encryption_enabled` adds a
  second layer beneath the customer keys

Know the contracts: CMK requires SINGLE-TENANT capacity -- a namespace
placed on a dedicated cluster or on the PREMIUM tier. Multi-tenant
BASIC/STANDARD namespaces share hardware and cannot take tenant keys;
Azure rejects the encryption patch on them. And the lifecycle is
ADD-ONLY: once CMK is enabled it can never be removed -- Azure has no
decrypt-back path. Deleting this resource intentionally changes NOTHING
on the namespace; returning to Microsoft-managed keys requires
replacing the namespace itself.

## Key Configuration

- `eventhub_namespace_id` -- the namespace to encrypt, referenced from
  an AzureEventHubNamespace output; must have single-tenant capacity
  and already carry the unwrapping identity (ForceNew)
- `key_vault_key_ids` -- 1-10 Key Vault keys by data-plane key ID;
  defaults to referencing an AzureKeyVaultKey's `versionless_id` output
  so rotation propagates automatically (pin versioned IDs only when
  compliance demands immutable versions). The keys' vault must have
  purge protection enabled
- `user_assigned_identity_id` -- the identity that unwraps the keys;
  MUST already ride the namespace's identity block with vault access
  granted. Unset uses the namespace's system-assigned identity (grant
  IT the vault access via its `identity_principal_id` output)
- `infrastructure_encryption_enabled` -- ForceNew: fixed the moment CMK
  is first configured

## Composition

```yaml
eventhubNamespaceId:
  valueFrom:
    kind: AzureEventHubNamespace
    name: premium-hubs
    fieldPath: status.outputs.namespace_id
keyVaultKeyIds:
  - valueFrom:
      kind: AzureKeyVaultKey
      name: streaming-key
      fieldPath: status.outputs.versionless_id
```

Sequence the grant first: an AzureRoleAssignment giving the namespace's
identity "Key Vault Crypto Service Encryption User" on the vault, then
this kind.

## Documentation

- [Design research](docs/README.md) -- the split-not-fold reasoning,
  lifecycle contracts
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
