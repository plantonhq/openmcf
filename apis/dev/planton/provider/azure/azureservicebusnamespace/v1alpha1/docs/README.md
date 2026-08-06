# AzureServiceBusNamespace -- Design Research

## The Resource

A Service Bus namespace (`Microsoft.ServiceBus/namespaces`, API
2024-01-01) is the container and billing boundary for enterprise
messaging. The component maps onto `azurerm_servicebus_namespace`
(azurerm v4.x, `internal/services/servicebus/servicebus_namespace_resource.go`),
parity-verified against pulumi-azure v6 (`servicebus.Namespace`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `namespace_name` | Required, ForceNew, provider regex (6-50, start letter, end alphanumeric) + the reserved `-sb`/`-mgmt` suffix bans |
| `sku` | enum | Basic/Standard/Premium; STANDARD default; the Premium-boundary migration is ForceNew (provider CustomizeDiff) |
| `capacity` | same | Premium messaging units {1,2,4,8,16}; the provider's apply-time capacity-vs-sku contract is front-loaded as spec CELs in BOTH directions |
| `premium_messaging_partitions` | same | {1,2,4}, ForceNew; Premium-required/forbidden pairing as CELs |
| `identity` | `identity` block | All three models (SystemAssigned, UserAssigned, combined) |
| `customer_managed_key` | `customer_managed_key` block | Premium-only (CEL); key FK → `AzureKeyVaultKey.versionless_id` (rotation-follows-latest); the unwrapping UAI must ride the identity block (CEL); removal-forces-recreate documented (provider CustomizeDiff) |
| `local_auth_enabled` | same | Default true; false = keyless posture |
| `public_network_access_enabled` | same | Default true |
| `network_rule_set` | `network_rule_set` block | Inline in v4 (no standalone resource in the service); Premium-only (CEL); DENY-requires-admitted-sources front-loads the provider's apply-time check; subnet rules FK `AzureSubnet` |
| `minimum_tls_version` | -- | Recorded skip below |
| `tags` | `tags` | User tags merged over metadata-derived identity tags |

Outputs: `namespace_id`, `namespace_name`, `endpoint`,
`identity_principal_id`, and the root SAS rule's four credential faces
(all sensitive; read from `RootManageSharedAccessKey`).

## Decomposition Decisions

- **Queues and topics are first-class kinds, not bundles.** They are
  many-per-namespace with independent lifecycles, owned by entity teams
  rather than the namespace's owner, individually FK-referenced (SAS
  rules, data-plane RBAC scopes, diagnostic targets), and azurerm itself
  models them as standalone resources addressed by `namespace_id`.
- **SAS authorization rules are a first-class kind**
  (AzureServiceBusAuthorizationRule) -- the credential surface
  applications actually consume, FK-referenced by the geo-DR pairing.
- **The geo-DR pairing is a first-class kind**
  (AzureServiceBusDisasterRecoveryConfig) -- independent lifecycle over
  two namespaces.
- **CMK folds into the namespace.** azurerm ships both an inline block
  and a standalone `azurerm_servicebus_namespace_customer_managed_key`
  resource; they manage the same singleton namespace property, and the
  standalone form exists for adopt-existing workflows (its delete is a
  documented no-op). The inline block is the honest declarative shape.
- **The network rule set folds** -- a singleton per-namespace document
  with no independent consumer.

## Recorded Skips (with reasons)

- **`minimum_tls_version`** -- a one-value constant, not a knob: Azure
  retired TLS 1.0/1.1 for Azure services and azurerm v5 restricts the
  field to "1.2". Both modules deliberately leave the provider default
  ("1.2") in place.
- **`zone_redundant`** -- not part of the azurerm v4 schema (Premium
  namespaces are zone redundant by default in regions with zones).

## Operational Behavior Worth Knowing

- **Standard namespaces provision in ~1-2 minutes**; Premium takes
  several minutes longer (dedicated capacity allocation).
- **The namespace name is a global DNS identity** with a post-delete
  hold -- a just-deleted name cannot be immediately recreated.
- **CMK is irreversible**: once customer-managed-key encryption is set,
  Azure cannot remove it; dropping the block replaces the namespace.
- **The firewall rides a separate ARM operation** after namespace
  create/update (the provider sequences it); Azure validates the
  block-level and namespace-level public-access dials agree.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `identity.user_assigned_identity_ids` / `customer_managed_key.user_assigned_identity_id`
  → `AzureUserAssignedIdentity.status.outputs.identity_id`
- `customer_managed_key.key_vault_key_id` → `AzureKeyVaultKey.status.outputs.versionless_id`
- `network_rule_set.network_rules[].subnet_id` → `AzureSubnet.status.outputs.subnet_id`
- `namespace_id` output ← AzureServiceBusQueue, AzureServiceBusTopic,
  AzureServiceBusAuthorizationRule, AzureServiceBusDisasterRecoveryConfig,
  AzurePrivateEndpoint, namespace-wide data-plane role assignments
