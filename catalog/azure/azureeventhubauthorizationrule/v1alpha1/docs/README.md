# AzureEventHubAuthorizationRule -- Design Research

## The Resource

Event Hubs SAS authorization rules exist as TWO ARM types --
`namespaces/authorizationRules` and
`namespaces/eventhubs/authorizationRules` -- with byte-identical shapes:
a name, the listen/send/manage rights trio, and six
key/connection-string faces. The component maps onto
`azurerm_eventhub_namespace_authorization_rule` /
`azurerm_eventhub_authorization_rule` (azurerm v4.x,
`internal/services/eventhub/`, one shared schema builder),
parity-verified against pulumi-azure v6.

## The One-Kind-Two-Types Verdict

One Planton kind with an exactly-one-of parent scope
(`namespace_id` XOR `event_hub_id`) dispatching to the two provider
resources, rather than two near-duplicate kinds:

- The two schemas are byte-identical (azurerm builds them from ONE
  shared schema function and ONE shared CustomizeDiff); what varies is
  only the parent scope.
- The scope IS configuration -- "which surface does this credential
  cover" -- and an XOR of typed FKs expresses it more honestly than two
  kinds that differ in nothing else.
- The orchestration-mode dispatch precedent: one kind whose
  discriminator picks between provider resources, with each dispatch
  path live-proven separately.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `rule_name` | Required, ForceNew; 1-60 chars, start/end alphanumeric, letters/numbers/periods/hyphens/underscores between; `RootManageSharedAccessKey` CEL-rejected (the built-in root rule -- its keys are AzureEventHubNamespace outputs) |
| `namespace_id` / `event_hub_id` | same | Exactly-one-of message CEL; each a typed FK; ForceNew |
| `listen` / `send` / `manage` | same | The provider's shared CustomizeDiff front-loaded as CELs: at-least-one; manage ⇒ listen ∧ send |

Outputs: `authorization_rule_id`, `rule_name` (the SharedAccessKeyName
clients present), and the six sensitive faces (`primary_key`,
`secondary_key`, both connection strings, both alias connection
strings). Alias faces populate only when the namespace carries an
AzureEventHubDisasterRecoveryConfig pairing. Hub-scoped connection
strings append `;EntityPath={hub}` -- ready-to-use as-is.

## Name-Addressed Provider Resources

Both azurerm resources still address the rule by discrete names
(`namespace_name` + `resource_group_name` [+ `eventhub_name`]), not by
parent id. The spec stays on the catalog's ARM-id grain -- one typed FK
per scope, no redundant contradictable name fields -- and both modules
parse the parent names from the resolved ARM id with anchored regexes
that fail loudly on a malformed id.

## Recorded Skips (with reasons)

- **None** -- the surface is complete. (ARM does not support tags on
  Event Hubs entities, so there is no tags field to carry.)

## Operational Behavior Worth Knowing

- **Renaming a rule regenerates its keys** (ForceNew) -- every client
  holding the old credential breaks. Rotate via the secondary instead:
  move clients to the secondary, regenerate the primary, move back.
- **Keyless namespaces** (`local_authentication_enabled` false) make
  every SAS rule's keys unusable -- the rule still exists, but
  authentication fails. The keyless alternative skips SAS rules
  entirely: Entra identities with data-plane roles (Azure Event Hubs
  Data Owner/Sender/Receiver) via AzureRoleAssignment.
- **The scope is fixed at creation** -- converting a namespace-wide rule
  to hub scope (or back) is a replace, with fresh keys.

## Composition

- `namespace_id` → `AzureEventHubNamespace.status.outputs.namespace_id`
- `event_hub_id` → `AzureEventHub.status.outputs.event_hub_id`
- `primary_connection_string` output ← application configuration,
  managed-secret wiring
