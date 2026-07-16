# Azure Event Hub

Creates an event hub inside an Azure Event Hubs namespace -- one partitioned, replayable event stream. Producers append to partitions; consumer groups read independently, each with its own offset. Capture archives every event to Blob Storage as Avro with no consumer application to run.

## What Gets Created

When you deploy an AzureEventHub resource, Planton provisions:

- **Event Hub** -- an `azurerm_eventhub` in the referenced namespace, with your partition layout, retention model, gate state, and optional capture-to-storage

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureEventHubNamespace** to create the hub in (referenced through `namespaceId`)
- **For capture**: an AzureStorageAccount and AzureStorageContainer to archive into

## Quick Start

Create a file `hub.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureEventHub
metadata:
  name: telemetry-stream
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureEventHub.telemetry-stream
spec:
  namespaceId:
    valueFrom:
      kind: AzureEventHubNamespace
      name: telemetry-hubs
      fieldPath: status.outputs.namespace_id
  eventHubName: telemetry
  partitionCount: 8
  messageRetention: 3
```

Deploy:

```shell
planton apply -f hub.yaml
```

Size `partitionCount` for peak up front on shared namespaces: the count can never be decreased, and only PREMIUM or dedicated-cluster namespaces can increase it later (Azure caps it at 32 on shared namespaces, 1024 on PREMIUM/dedicated). Set exactly one retention model -- `messageRetention` in days (tier caps: 1 BASIC, 7 STANDARD, 90 PREMIUM/dedicated) or `retentionDescription` for hour-granular windows and Kafka-style compaction.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `event_hub_id` | The parent reference for consumer groups and hub-scoped authorization rules, and the hub-level data-plane RBAC scope |
| `event_hub_name` | What producers and consumers address within the namespace -- the Kafka topic name on the Kafka endpoint |
| `partition_ids` | The partition identifiers as Azure assigned them -- what partition-aware consumers enumerate |

## Related Resources

- [Azure Event Hub Namespace](/docs/catalog/azure/azureeventhubnamespace) -- the container and billing boundary
- [Azure Event Hub Consumer Group](/docs/catalog/azure/azureeventhubconsumergroup) -- independent read cursors
- [Azure Event Hub Authorization Rule](/docs/catalog/azure/azureeventhubauthorizationrule) -- hub-scoped SAS credentials
- [Azure Storage Account](/docs/catalog/azure/azurestorageaccount) / [Azure Storage Container](/docs/catalog/azure/azurestoragecontainer) -- the capture destination
