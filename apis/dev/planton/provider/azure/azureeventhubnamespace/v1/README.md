# AzureEventHubNamespace

An Azure Event Hubs namespace: the container and billing boundary for
high-throughput event streaming. The namespace sets the pricing tier,
throughput capacity, network posture, and authentication mode; the
streaming entities compose onto it as first-class kinds --
AzureEventHub (streams, with capture-to-storage),
AzureEventHubConsumerGroup, AzureEventHubAuthorizationRule,
AzureEventHubSchemaGroup, AzureEventHubDisasterRecoveryConfig,
AzureEventHubCluster, and AzureEventHubNamespaceCustomerManagedKey.

## When to Use

Use AzureEventHubNamespace when you need:

- **High-throughput event streaming** -- telemetry, log aggregation, IoT
  ingestion, change-data capture at millions of events per second
  (Service Bus is the enterprise-messaging sibling; Storage Queues the
  minimal one)
- **A Kafka endpoint without a Kafka cluster** -- STANDARD and above
  expose `{name}.servicebus.windows.net:9093` to existing Kafka
  producers and consumers
- **A shared streaming boundary** -- one namespace per environment or
  domain, with stream teams owning their hubs, consumer groups, and
  credentials independently

## Key Configuration

- `namespace_name` -- globally unique; becomes the endpoint
  `{name}.servicebus.windows.net` and the Kafka bootstrap host on port
  9093 (ForceNew)
- `sku` -- BASIC (single consumer group, no Kafka), STANDARD (default;
  full-featured multi-tenant with auto-inflate), PREMIUM (reserved
  processing units; migrating in/out replaces the namespace)
- `capacity` -- throughput units (1-40) on BASIC/STANDARD; processing
  units (1/2/4/8/16) on PREMIUM
- `auto_inflate_enabled` + `maximum_throughput_units` -- STANDARD's
  elastic scaling; grows TUs under load but never shrinks them back
- `dedicated_cluster_id` -- place the namespace on an
  AzureEventHubCluster for single-tenant capacity, 1024-partition hubs,
  90-day retention, and CMK eligibility (ForceNew)
- `local_authentication_enabled` -- false for the keyless posture
  (Entra-only data-plane auth; every SAS key stops working)
- `network_rule_sets` -- the namespace firewall, not available on BASIC
  (DENY + admitted IPs and subnets is the production posture)

## Composition

```yaml
resourceGroup:
  valueFrom:
    kind: AzureResourceGroup
    name: streaming-rg
    fieldPath: status.outputs.resource_group_name
```

Children reference `status.outputs.namespace_id`; applications consume
the root SAS credential outputs (quick starts) or scoped
AzureEventHubAuthorizationRule credentials (production).

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
