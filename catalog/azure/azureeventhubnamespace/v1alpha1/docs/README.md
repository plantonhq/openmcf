# AzureEventHubNamespace -- Design Research

## The Resource

An Event Hubs namespace (`Microsoft.EventHub/namespaces`) is the
container and billing boundary for high-throughput event streaming. The
component maps onto `azurerm_eventhub_namespace` (azurerm v4.x,
`internal/services/eventhub/eventhub_namespace_resource.go`),
parity-verified against pulumi-azure v6 (`eventhub.EventHubNamespace`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `namespace_name` | Required, ForceNew, 6-50 chars, start letter, end alphanumeric; the global DNS identity (`{name}.servicebus.windows.net` -- Event Hubs shares the Service Bus DNS zone) and the Kafka bootstrap host on port 9093 |
| `sku` | enum | Basic/Standard/Premium; STANDARD default; the Premium boundary is ForceNew (provider CustomizeDiff -- Azure cannot convert across the reserved/multi-tenant boundary in place); BASIC <-> STANDARD updates in place |
| `capacity` | same | Tier-dependent semantics: throughput units (1-40) on BASIC/STANDARD, processing units {1,2,4,8,16} on PREMIUM -- both value sets front-loaded as spec CELs |
| `auto_inflate_enabled` | same | STANDARD's elastic scaling; grows TUs up to the ceiling but never shrinks them back |
| `maximum_throughput_units` | same | 0-40; the ceiling/enable pairing is validated by ARM at apply time (see below) |
| `dedicated_cluster_id` | same | ForceNew; FK → `AzureEventHubCluster.cluster_id`; placement buys single-tenant capacity, 1024-partition hubs, 90-day retention, and CMK eligibility |
| `identity` | `identity` block | All three models (SystemAssigned, UserAssigned, combined); required for identity-based capture auth and CMK unwrapping |
| `local_authentication_enabled` | same | Default true; false = keyless posture (every SAS key, including the root rule's, stops working) |
| `public_network_access_enabled` | same | Default true; must agree with the rule set's block-level dial when that block is declared |
| `network_rulesets` | `network_rule_sets` block | Inline fold (see below); not on BASIC (CEL); explicit default_action required; DENY-requires-admitted-sources front-loads Azure's check; subnet rules FK `AzureSubnet` |
| `minimum_tls_version` | -- | Recorded skip below |
| `tags` | `tags` | User tags merged over metadata-derived identity tags |

Outputs: `namespace_id` (the parent seam for the whole family and the
private-endpoint target, subresource "namespace"), `namespace_name`,
`identity_principal_id`, and the root SAS rule's
(`RootManageSharedAccessKey`) six credential faces -- four addressing the
namespace hostname plus the alias pair addressing the geo-DR alias
hostname, populated only when an AzureEventHubDisasterRecoveryConfig
pairing exists. All six are sensitive.

## Decomposition Decisions

- **Event hubs are a first-class kind, not a bundle.** They are
  many-per-namespace with independent lifecycles, owned by stream teams
  rather than the namespace's owner, individually FK-referenced
  (consumer groups, hub-scoped SAS rules, hub-level data-plane RBAC),
  and azurerm itself models them as standalone resources addressed by
  `namespace_id`.
- **Consumer groups, SAS rules, schema groups, the geo-DR pairing, the
  dedicated cluster, and CMK are all first-class kinds** --
  AzureEventHubConsumerGroup, AzureEventHubAuthorizationRule,
  AzureEventHubSchemaGroup, AzureEventHubDisasterRecoveryConfig,
  AzureEventHubCluster, and AzureEventHubNamespaceCustomerManagedKey --
  each referencing `namespace_id`. Nothing is bundled into the
  namespace's spec.
- **The network rule set folds into the namespace.** The provider models
  it as an inline block on the namespace resource, it has no independent
  lifecycle, and nothing references it -- a singleton per-namespace
  document with no independent consumer.

## Recorded Skips (with reasons)

- **`minimum_tls_version`** -- a one-value constant, not a knob: Azure
  retired TLS 1.0/1.1 for Azure services and azurerm v5 restricts the
  field to "1.2". Both modules deliberately leave the provider default
  ("1.2") in place.
- **`zone_redundant`** -- not part of azurerm v4's surface: zone
  redundancy is service-managed (namespaces in regions with availability
  zones get it automatically), so there is nothing to configure.
- **Auto-inflate schema pairing** -- the provider has no schema-level
  pairing between `auto_inflate_enabled` and `maximum_throughput_units`;
  ARM validates the pairing at apply time (a non-zero ceiling without
  auto-inflate enabled is rejected server-side). The provider's only
  guard is zeroing the ceiling itself on a downgrade to Basic. The spec
  mirrors the provider: both fields are independent optionals.

## Apply-Time Contracts

Contracts Azure or the provider enforces at plan/apply rather than in
the schema:

- **The Premium boundary is ForceNew** -- the provider's CustomizeDiff
  forces replacement when the sku crosses into or out of Premium;
  BASIC <-> STANDARD updates in place.
- **TLS cannot be unset** -- the provider rejects clearing
  `minimum_tls_version` once the namespace exists (moot here: both
  modules never send it).
- **The auto-inflate ceiling/enable pairing** -- ARM rejects a non-zero
  ceiling on a namespace without auto-inflate enabled.
- **The block-level and namespace-level public-access dials must agree**
  -- Azure validates the pair server-side; the spec front-loads it as a
  message CEL.

## Operational Behavior Worth Knowing

- **Auto-inflate only scales UP.** Azure grows TUs under load up to
  `maximum_throughput_units` but never shrinks them back -- scale-down
  after a traffic spike is a manual `capacity` edit.
- **The namespace name is a global DNS identity** -- it becomes both the
  AMQP endpoint and the Kafka bootstrap host
  `{name}.servicebus.windows.net:9093`.
- **The firewall rides a separate ARM operation** after namespace
  create/update (the provider sequences it); each `ip_rules` entry is an
  allow rule -- Azure's per-rule action accepts exactly one value.
- **Keyless posture is namespace-wide**: `local_authentication_enabled`
  false disables every SAS rule's keys, including the root rule surfaced
  in this kind's outputs. Pair it with AzureRoleAssignment grants of the
  data-plane roles (Azure Event Hubs Data Owner/Sender/Receiver).
- **ARM tags are Azure's governance surface** -- Azure Policy enforces
  them and Microsoft Cost Management groups by them; user tags merge
  over the platform's identity tags and win on key conflicts.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `dedicated_cluster_id` → `AzureEventHubCluster.status.outputs.cluster_id`
- `identity.user_assigned_identity_ids` → `AzureUserAssignedIdentity.status.outputs.identity_id`
- `network_rule_sets.virtual_network_rules[].subnet_id` → `AzureSubnet.status.outputs.subnet_id`
- `namespace_id` output ← AzureEventHub, AzureEventHubAuthorizationRule,
  AzureEventHubSchemaGroup, AzureEventHubDisasterRecoveryConfig,
  AzureEventHubNamespaceCustomerManagedKey, AzurePrivateEndpoint,
  namespace-wide data-plane role assignments
