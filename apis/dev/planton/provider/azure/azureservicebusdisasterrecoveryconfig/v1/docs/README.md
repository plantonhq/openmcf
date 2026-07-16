# AzureServiceBusDisasterRecoveryConfig -- Design Research

## The Resource

A Service Bus geo-DR config
(`Microsoft.ServiceBus/namespaces/disasterRecoveryConfigs`) pairs two
Premium namespaces under a failover-stable alias. The component maps
onto `azurerm_servicebus_namespace_disaster_recovery_config` (azurerm
v4.x, `internal/services/servicebus/servicebus_namespace_disaster_recovery_config_resource.go`),
parity-verified against pulumi-azure v6
(`servicebus.NamespaceDisasterRecoveryConfig`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `alias_name` | Required, ForceNew; becomes the alias DNS identity (namespace-name uniqueness scope); the provider carries no schema validator -- min 1/max 50 kept, DNS reality documented |
| `primary_namespace_id` | same | FK → AzureServiceBusNamespace; ForceNew |
| `partner_namespace_id` | same | FK → AzureServiceBusNamespace; change = break-pairing + re-pair (provider-managed) |
| `alias_authorization_rule_id` | same | Optional FK → AzureServiceBusAuthorizationRule.authorization_rule_id (zero-translation); unset = the root rule |

Outputs: `disaster_recovery_config_id`, `alias_name`, and four
sensitive credential faces (both alias connection strings + the paired
rule's two keys).

## Decomposition Decision

A first-class kind: the pairing has an independent lifecycle over two
namespaces (the Redis linked-server class), is created/broken/re-paired
without touching either namespace, and its alias credentials are a
distinct credential surface.

## Azure's Pairing Contracts (apply-time -- both live namespaces)

- Both namespaces PREMIUM; different regions; the partner EMPTY at
  pairing time. These involve two referenced resources' live state,
  which validation cannot see -- documented, enforced verbatim by ARM.

## Operational Behavior Worth Knowing

- **Failover is triggered on the SECONDARY** (portal/CLI/SDK), promotes
  it to primary, and BREAKS the pairing -- re-pair to a new partner
  afterwards. It is deliberately not a field on this spec.
- **The provider choreographs the lifecycle**: create waits for
  Succeeded; partner change breaks and re-pairs; destroy breaks the
  pairing, deletes, and polls until the alias NAME frees so immediate
  re-creation does not collide.
- **Metadata replicates; messages do not.** In-flight messages in the
  failed region stay there.

## E2E Profile Note

Live proof requires two Premium namespaces (the slow tier) in two
regions plus pair/break-pair settling -- recorded as an offline-gated
kind with ready-to-run scenarios; see the E2E profile for the standing
deferral and its unblock path.

## Composition

- `primary_namespace_id` / `partner_namespace_id` →
  `AzureServiceBusNamespace.status.outputs.namespace_id`
- `alias_authorization_rule_id` →
  `AzureServiceBusAuthorizationRule.status.outputs.authorization_rule_id`
- `primary_connection_string_alias` output ← DR-aware application
  configuration
