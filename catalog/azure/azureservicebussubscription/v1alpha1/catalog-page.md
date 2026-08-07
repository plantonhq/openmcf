# Azure Service Bus Subscription

Creates a subscription under an Azure Service Bus topic -- an independent, optionally filtered view of the topic's stream with its own consumer semantics. Filter rules are part of the subscription document.

## What Gets Created

When you deploy an AzureServiceBusSubscription resource, Planton provisions:

- **Service Bus Subscription** -- an `azurerm_servicebus_subscription` on the referenced topic (via its ARM id), with your consumer dials
- **Filter Rules** -- one `azurerm_servicebus_subscription_rule` per declared rule (SQL or correlation)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureServiceBusTopic** to attach to (referenced through `topicId`)

## Quick Start

Create a file `subscription.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServiceBusSubscription
metadata:
  name: billing-consumer
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureServiceBusSubscription.billing-consumer
spec:
  topicId:
    valueFrom:
      kind: AzureServiceBusTopic
      name: order-events
      fieldPath: status.outputs.topic_id
  subscriptionName: billing
  maxDeliveryCount: 10
```

Deploy:

```shell
planton apply -f subscription.yaml
```

Without rules, the subscription receives EVERYTHING (Azure auto-creates a `$Default` catch-all). Declared rules ADD to that catch-all with OR semantics -- to make your filters restrictive, remove the catch-all once after creation (`az servicebus topic subscription rule delete --name '$Default' ...`); the service-created rule cannot be declared or adopted by the management plane.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `subscription_id` | The subscription's ARM identity |
| `subscription_name` | What consumers reference (with the topic name) when receiving |
| `topic_name` / `namespace_name` | The full receive triple, without extra references |

## Related Resources

- [Azure Service Bus Topic](/docs/catalog/azure/azureservicebustopic) -- the parent topic
- [Azure Service Bus Queue](/docs/catalog/azure/azureservicebusqueue) -- pair with `forwardTo` for fan-out-then-collect routing
