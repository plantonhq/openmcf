# AzureServiceBusDisasterRecoveryConfig

A geo-disaster-recovery pairing between two PREMIUM Service Bus
namespaces: metadata (queues, topics, subscriptions, rules -- not
message data) continuously replicates from the primary to the partner,
and a failover-stable alias DNS name
(`{alias_name}.servicebus.windows.net`) fronts whichever namespace is
currently primary.

## When to Use

Use AzureServiceBusDisasterRecoveryConfig when you need:

- **Regional DR for messaging topology** -- a warm standby namespace
  whose entity structure always matches production
- **Failover without client reconfiguration** -- clients connect
  through the alias, so promoting the secondary changes nothing on
  their side

Know the boundary: geo-DR replicates METADATA, not messages. In-flight
messages in the failed region are not moved.

## Key Configuration

- `alias_name` -- becomes the global DNS identity; ForceNew
- `primary_namespace_id` -- the active namespace (PREMIUM; ForceNew)
- `partner_namespace_id` -- PREMIUM, different region, EMPTY at pairing
  time; changing it breaks and re-pairs
- `alias_authorization_rule_id` -- optional; a namespace-scoped
  AzureServiceBusAuthorizationRule for least-privilege alias
  credentials (unset = the root rule)

Failover is an operational action taken on the SECONDARY during an
incident (portal/CLI/SDK) -- never a config change here. Deleting this
resource breaks the pairing gracefully.

## Composition

```yaml
primaryNamespaceId:
  valueFrom:
    kind: AzureServiceBusNamespace
    name: orders-bus-eastus
    fieldPath: status.outputs.namespace_id
partnerNamespaceId:
  valueFrom:
    kind: AzureServiceBusNamespace
    name: orders-bus-westus
    fieldPath: status.outputs.namespace_id
```

## Documentation

  choreography
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
