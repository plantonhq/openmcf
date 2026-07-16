---
title: "Event Hub Consumer Group"
description: "Event Hub Consumer Group deployment documentation"
icon: "package"
order: 100
componentName: "azureeventhubconsumergroup"
---

# Azure Event Hub Consumer Group

Creates a consumer group on an Azure event hub -- one application's independent read cursor over the hub's partitions, so many applications consume the same stream at their own pace.

## What Gets Created

When you deploy an AzureEventHubConsumerGroup resource, Planton provisions:

- **Consumer Group** -- an `azurerm_eventhub_consumer_group` on the referenced hub, with optional free-form user metadata recording whose cursor it is

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureEventHub** to create the group on (referenced through `eventHubId`); the hub's namespace must be STANDARD or above -- BASIC allows no additional consumer groups

## Quick Start

Create a file `consumer-group.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureEventHubConsumerGroup
metadata:
  name: analytics-consumer
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureEventHubConsumerGroup.analytics-consumer
spec:
  eventHubId:
    valueFrom:
      kind: AzureEventHub
      name: telemetry-stream
      fieldPath: status.outputs.event_hub_id
  consumerGroupName: analytics
  userMetadata: "owner=data-platform; app=stream-analytics"
```

Deploy:

```shell
planton apply -f consumer-group.yaml
```

Give every consuming application its OWN group -- shared cursors collide. The name `$Default` is reserved: Azure creates that group on every hub automatically, and it cannot be declared here. `consumerGroupName` is fixed at creation; renaming replaces the group and resets its consumers' stored offsets. Azure enforces tier quotas at apply time (STANDARD: 20 groups per hub).

## Key Outputs

| Output | Purpose |
|--------|---------|
| `consumer_group_id` | The Azure Resource Manager ID of the group |
| `consumer_group_name` | What consumer applications pass to their SDK client alongside the hub name |

## Related Resources

- [Azure Event Hub](/docs/catalog/azure/event-hub) -- the parent stream
- [Azure Event Hub Namespace](/docs/catalog/azure/event-hub-namespace) -- the container whose tier sets the group quota
- [Azure Event Hub Authorization Rule](/docs/catalog/azure/event-hub-authorization-rule) -- SAS credentials for consumers
