# AzureServiceBusNamespace

An Azure Service Bus namespace: the container and billing boundary for
enterprise messaging. The namespace sets the pricing tier, network
posture, encryption, and authentication mode; the messaging entities
compose onto it as first-class kinds -- AzureServiceBusQueue,
AzureServiceBusTopic (with AzureServiceBusSubscription),
AzureServiceBusAuthorizationRule, and
AzureServiceBusDisasterRecoveryConfig.

## When to Use

Use AzureServiceBusNamespace when you need:

- **Enterprise messaging** -- ordered delivery, sessions, duplicate
  detection, dead-lettering, transactions (Event Hubs is the streaming
  sibling; Storage Queues the minimal one)
- **A shared messaging boundary** -- one namespace per environment or
  domain, with entity teams owning their queues/topics independently
- **Premium isolation** -- dedicated messaging units, VNet integration,
  CMK encryption, geo-DR, and 100 MB messages on the PREMIUM tier

## Key Configuration

- `namespace_name` -- globally unique; becomes the endpoint
  `{name}.servicebus.windows.net` (ForceNew)
- `sku` -- BASIC (queues only), STANDARD (default; full-featured
  multi-tenant), PREMIUM (dedicated capacity; migrating in/out replaces
  the namespace)
- `capacity` + `premium_messaging_partitions` -- PREMIUM's messaging
  units and partition layout (partitions are fixed at creation)
- `identity` + `customer_managed_key` -- BYOK encryption (PREMIUM;
  irreversible once set)
- `local_auth_enabled` -- false for the keyless posture (Entra-only
  data-plane auth; every SAS key stops working)
- `network_rule_set` -- the PREMIUM firewall (DENY + admitted IPs and
  subnets is the production posture)

## Composition

```yaml
resourceGroup:
  valueFrom:
    kind: AzureResourceGroup
    name: messaging-rg
    fieldPath: status.outputs.resource_group_name
```

Children reference `status.outputs.namespace_id`; applications consume
the root SAS credential outputs (quick starts) or scoped
AzureServiceBusAuthorizationRule credentials (production).

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
