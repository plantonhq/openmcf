# AzureEventHubDisasterRecoveryConfig -- Design Research

## The Resource

An Event Hubs geo-DR config
(`Microsoft.EventHub/namespaces/disasterRecoveryConfigs`) pairs two
namespaces under a failover-stable alias. The component maps onto
`azurerm_eventhub_namespace_disaster_recovery_config` (azurerm v4.x,
`internal/services/eventhub/eventhub_namespace_disaster_recovery_config_resource.go`),
parity-verified against pulumi-azure v6
(`eventhub.EventhubNamespaceDisasterRecoveryConfig`).

Metadata (hubs, consumer groups, authorization rules) continuously
replicates primary → partner; EVENT DATA does not. The alias
(`{alias_name}.servicebus.windows.net`) is the client-facing DNS
identity that fronts whichever namespace is currently primary.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `alias_name` | Required, ForceNew; becomes the alias DNS identity, sharing the namespace-name uniqueness scope (globally unique); 1-60 chars, alphanumeric ends, `[-._]` interior |
| `namespace_name` / `resource_group_name` | `primary_namespace_id` | azurerm addresses the primary side by discrete names; both modules parse them from the spec's single ARM-id reference with an anchored regex that fails loudly on a malformed id; ForceNew |
| `partner_namespace_id` | same | FK → AzureEventHubNamespace; change = break-pairing + re-pair (provider-managed) |

Outputs: `disaster_recovery_config_id`, `alias_name`.

## The No-Credential-Outputs Design

This kind exports the pairing's IDENTITY only. Azure exposes
alias-addressed connection strings through authorization rules, not
through the DR resource itself -- so DR-aware clients take the
`*_connection_string_alias` outputs from AzureEventHubNamespace (the
root rule) or AzureEventHubAuthorizationRule (scoped rules), which
surface once a pairing exists. Duplicating that credential surface here
would invent outputs Azure does not provide.

## Azure's Pairing Contracts (apply-time -- both live namespaces)

- Different regions; SAME tier, STANDARD or higher (geo-DR is not
  available on BASIC); the partner EMPTY (no hubs) at pairing time.
  These involve two referenced resources' live state, which validation
  cannot see -- documented on the spec, enforced verbatim by ARM.

## Lifecycle Choreography (provider-managed)

- **Create** waits for the pairing to reach the Succeeded provisioning
  state (polling Accepted → Succeeded).
- **Partner change** breaks the existing pairing first, waits, then
  re-pairs to the new partner.
- **Destroy** breaks the pairing, deletes the config, then waits BOTH
  for the config to 404 AND for the alias NAME to be released by
  Azure's name-availability check -- the alias name stays reserved
  briefly after deletion, so destroys take minutes by the service's
  own design.

## Operational Behavior Worth Knowing

- **Failover is triggered on the SECONDARY** (portal/CLI/SDK), promotes
  it to primary, and breaks the pairing -- re-pair to a new partner
  afterwards. It is deliberately not a field on this spec.
- **Metadata replicates; events do not.** After failover, consumers
  resume on the partner's (empty) hubs -- plan retention and recovery
  objectives accordingly.
- **Deleting the pairing is graceful**: both namespaces keep running
  independently.

## Composition

- `primary_namespace_id` / `partner_namespace_id` →
  `AzureEventHubNamespace.status.outputs.namespace_id`
- Alias credentials ← `AzureEventHubNamespace` /
  `AzureEventHubAuthorizationRule` `*_connection_string_alias` outputs
