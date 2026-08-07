# Geo-DR Pairing

This preset pairs two Event Hubs namespaces across regions under a
failover-stable alias: hub/consumer-group/rule metadata continuously
replicates to the partner, and clients connected through the alias
survive a regional failover without reconfiguration.

## When to Use

- Streaming estates with a regional-outage recovery objective
- Note what replicates: METADATA only -- event data does not replicate;
  after failover, consumers resume on the partner's (empty) hubs

## Key Configuration Choices

- **The alias is the client-facing identity** -- point producers and
  consumers at `{aliasName}.servicebus.windows.net` (the
  `*_connection_string_alias` outputs on the namespace and
  authorization-rule kinds), never at either namespace directly
- **The partner must be empty and in a different region** -- Azure
  validates both at pairing time
- **Failover is an operational action** performed from the SECONDARY
  side during an incident, never a manifest change; changing
  `partnerNamespaceId` here breaks and re-pairs instead

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `myorg-telemetry-alias` | Globally unique alias (shares the namespace name scope) | Your naming convention |
| `my-telemetry-hubs` | The primary namespace | Your streaming composition |
| `my-telemetry-hubs-dr` | The empty partner namespace in the recovery region | Same composition |
