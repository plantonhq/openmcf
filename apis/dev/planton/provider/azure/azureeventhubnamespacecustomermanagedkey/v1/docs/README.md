# AzureEventHubNamespaceCustomerManagedKey -- Design Research

## The Resource

Event Hubs customer-managed-key encryption is a PROPERTY of a namespace
(`Microsoft.EventHub/namespaces` encryption settings), not an ARM
object of its own. The component maps onto
`azurerm_eventhub_namespace_customer_managed_key` (azurerm v4.x,
`internal/services/eventhub/eventhub_namespace_customer_managed_key_resource.go`),
parity-verified against pulumi-azure v6
(`eventhub.NamespaceCustomerManagedKey`).

## Why a Separate Kind (split, not fold)

Azure models CMK as a configuration applied ONTO an existing namespace,
and the split mirrors that grain for a causal reason: encrypting with
the namespace's own system-assigned identity is only POSSIBLE as a
second step. The identity does not exist until the namespace does, and
it needs wrap/unwrap access on the vault before the encryption patch is
accepted. The working order is:

1. Create the namespace with its identity block (system-assigned or
   user-assigned).
2. Grant the identity vault access -- an AzureRoleAssignment for
   "Key Vault Crypto Service Encryption User" on the vault, scoped to
   the namespace's `identity_principal_id` output (system-assigned) or
   the user-assigned identity's principal.
3. Apply this kind.

A create-time block folded into the namespace spec could never express
that sequence for the system-assigned case.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `eventhub_namespace_id` | same | FK → AzureEventHubNamespace; ForceNew -- the configuration is bound to its namespace for life |
| `key_vault_key_ids` | same | 1-10 data-plane key IDs; FK defaults to AzureKeyVaultKey's `versionless_id` so vault rotation propagates automatically |
| `infrastructure_encryption_enabled` | same | Optional bool; ForceNew -- fixed the moment CMK is first configured; unset leaves Azure's default (false) |
| `user_assigned_identity_id` | same | Optional FK → AzureUserAssignedIdentity; must ALREADY ride the namespace's identity block with vault access; unset = system-assigned identity |

Output: `customer_managed_key_id` -- the provider's identity for the
configuration IS the namespace's ARM id (no ARM object of its own; no
tags for the same reason).

## The Add-Only Contract

Once CMK is enabled it can never be removed -- Azure has no
decrypt-back path. The provider's Delete is deliberately a NO-OP:
destroying this resource changes nothing on the namespace, and
returning to Microsoft-managed keys requires replacing the namespace
itself. Key ROTATION is routine, though: versionless references follow
the vault's latest version automatically.

## Azure's Platform Contracts (apply-time)

- **Single-tenant capacity required**: the namespace must sit on a
  dedicated cluster (`dedicated_cluster_id`) or be PREMIUM.
  Multi-tenant BASIC/STANDARD namespaces share hardware and cannot take
  tenant keys; Azure rejects the encryption patch. The contract depends
  on the referenced namespace's live placement, which validation cannot
  see -- documented, enforced verbatim by ARM.
- **The identity contract**: a user-assigned identity named here must
  already be attached via the namespace's identity block, with
  wrap/unwrap access granted -- Azure rejects the patch otherwise.
- **Purge protection**: the keys' vault must have purge protection
  enabled.

## E2E Exclusion

Live E2E is excluded for this kind (see `e2e/profile.yaml`): CMK needs
single-tenant capacity (a dedicated cluster -- itself live-blocked by
the 4-hour deletion moratorium -- or a PREMIUM namespace) plus a
purge-protected vault chain whose recycle-bin hold outlives the run;
and because the kind is add-only, verify-absent after destroy is
meaningless (the delete is a no-op) and the encrypted fixture namespace
itself must be destroyed to clean up. The component stands on its
offline gate -- full parity audit, a complete plan rendering the key
array, identity, and infrastructure-encryption seams, spec tests
covering every validation rule, and the outputs conformance case.

## Composition

- `eventhub_namespace_id` → `AzureEventHubNamespace.status.outputs.namespace_id`
- `key_vault_key_ids[]` → `AzureKeyVaultKey.status.outputs.versionless_id`
- `user_assigned_identity_id` → `AzureUserAssignedIdentity.status.outputs.identity_id`
- The vault grant ← AzureRoleAssignment ("Key Vault Crypto Service
  Encryption User") on the unwrapping identity's principal
