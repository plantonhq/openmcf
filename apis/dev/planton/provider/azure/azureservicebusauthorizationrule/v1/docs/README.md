# AzureServiceBusAuthorizationRule -- Design Research

## The Resource

Service Bus SAS authorization rules exist as THREE ARM types --
`namespaces/authorizationRules`, `namespaces/queues/authorizationRules`,
and `namespaces/topics/authorizationRules` -- with byte-identical
shapes: a name, the listen/send/manage rights trio, and six
key/connection-string faces. The component maps onto
`azurerm_servicebus_namespace_authorization_rule` /
`azurerm_servicebus_queue_authorization_rule` /
`azurerm_servicebus_topic_authorization_rule` (azurerm v4.x,
`internal/services/servicebus/`, one shared schema builder),
parity-verified against pulumi-azure v6.

## The One-Kind-Three-Types Verdict

One Planton kind with an exactly-one-of parent scope
(`namespace_id` XOR `queue_id` XOR `topic_id`) dispatching to the three
provider resources, rather than three near-duplicate kinds:

- The three schemas are byte-identical (azurerm builds them from ONE
  shared schema function and ONE shared CustomizeDiff); what varies is
  only the parent scope.
- The scope IS configuration -- "which entity does this credential
  cover" -- and an XOR of typed FKs expresses it more honestly than
  three kinds that differ in nothing else.
- The orchestration-mode dispatch precedent: one kind whose
  discriminator picks between provider resources, with each dispatch
  path live-proven separately.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `rule_name` | Required, ForceNew, provider regex (≤50); `RootManageSharedAccessKey` CEL-rejected (the built-in root rule -- its keys are namespace outputs) |
| `namespace_id` / `queue_id` / `topic_id` | same | Exactly-one-of message CEL; each a typed FK; ForceNew |
| `listen` / `send` / `manage` | same | The provider's shared CustomizeDiff front-loaded as CELs: at-least-one; manage ⇒ listen ∧ send |

Outputs: `authorization_rule_id`, `rule_name`, and the six sensitive
faces (`primary_key`, `secondary_key`, both connection strings, both
alias connection strings). Alias faces populate only when the namespace
carries a geo-DR pairing.

## Recorded Skips (with reasons)

- **None** -- the surface is complete.

## Operational Behavior Worth Knowing

- **Renaming a rule regenerates its keys** (ForceNew) -- every client
  holding the old credential breaks. Rotate via the secondary instead.
- **On a geo-DR-paired Premium namespace**, the provider waits for
  pairing replication to settle after create/delete -- rule operations
  can take a minute longer there.
- **Keyless namespaces** (`local_auth_enabled` false) make every SAS
  rule's keys unusable -- the rule still exists, but authentication
  fails.

## Composition

- `namespace_id` → `AzureServiceBusNamespace.status.outputs.namespace_id`
- `queue_id` → `AzureServiceBusQueue.status.outputs.queue_id`
- `topic_id` → `AzureServiceBusTopic.status.outputs.topic_id`
- `authorization_rule_id` output ←
  `AzureServiceBusDisasterRecoveryConfig.alias_authorization_rule_id`
  (zero-translation)
- `primary_connection_string` output ← application configuration,
  managed-secret wiring
