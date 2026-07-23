# AzureEventHubDisasterRecoveryConfig

A geo-disaster-recovery pairing between two Event Hubs namespaces:
metadata (hubs, consumer groups, authorization rules -- not event
data) continuously replicates from the primary to the partner, and a
failover-stable alias DNS name
(`{alias_name}.servicebus.windows.net`) fronts whichever namespace is
currently primary.

## When to Use

Use AzureEventHubDisasterRecoveryConfig when you need:

- **Regional DR for streaming topology** -- a warm standby namespace
  whose hub/consumer-group/rule structure always matches production
- **Failover without client reconfiguration** -- clients connect
  through the alias, so promoting the secondary changes nothing on
  their side

Know the boundary: geo-DR replicates METADATA, not events. Event data
in the failed region is not moved; after failover, consumers resume on
the partner's (empty) hubs.

## Key Configuration

- `alias_name` -- becomes the global DNS identity (it shares the
  namespace name uniqueness scope); ForceNew
- `primary_namespace_id` -- the active namespace; the pairing lives
  under it (ForceNew)
- `partner_namespace_id` -- a DIFFERENT region, the same tier
  (STANDARD or higher -- geo-DR is not available on BASIC), and EMPTY
  at pairing time; changing it breaks and re-pairs

Azure validates the pairing contracts at apply time, since they involve
both live namespaces. Failover is an operational action taken on the
SECONDARY during an incident (portal/CLI/SDK) -- never a config change
here. Deleting this resource breaks the pairing gracefully; both
namespaces keep running independently.

This kind exports the pairing's identity only. Alias-addressed
credentials are an authorization-rule surface in Azure, so DR-aware
clients take the `*_connection_string_alias` outputs from
AzureEventHubNamespace (the root rule) or
AzureEventHubAuthorizationRule (scoped rules).

## Composition

```yaml
primaryNamespaceId:
  valueFrom:
    kind: AzureEventHubNamespace
    name: telemetry-hubs-eastus
    fieldPath: status.outputs.namespace_id
partnerNamespaceId:
  valueFrom:
    kind: AzureEventHubNamespace
    name: telemetry-hubs-westus
    fieldPath: status.outputs.namespace_id
```

## Documentation

- [Design research](docs/README.md) -- pairing contracts, lifecycle
  choreography
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
