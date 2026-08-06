---
title: "Event Hub Namespace"
description: "Event Hub Namespace deployment documentation"
icon: "package"
order: 100
componentName: "azureeventhubnamespace"
---

# Azure Event Hub Namespace

Creates an Azure Event Hubs namespace -- the container and billing boundary for high-throughput event streaming. The namespace sets the pricing tier, throughput capacity, network posture, and authentication mode; event hubs, consumer groups, SAS credentials, schema groups, geo-DR pairings, and CMK encryption compose onto it as first-class resources.

## What Gets Created

When you deploy an AzureEventHubNamespace resource, Planton provisions:

- **Event Hubs Namespace** -- an `azurerm_eventhub_namespace` in the referenced resource group, with your chosen tier, optional managed identity, auto-inflate, dedicated-cluster placement, and network rules

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureResourceGroup** to create the namespace in (referenced through `resourceGroup`)

## Quick Start

Create a file `namespace.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventHubNamespace
metadata:
  name: telemetry-hubs
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureEventHubNamespace.telemetry-hubs
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: streaming-rg
      fieldPath: status.outputs.resource_group_name
  namespaceName: myorg-telemetry-hubs
```

Deploy:

```shell
planton apply -f namespace.yaml
```

Unset `sku` deploys STANDARD -- the full-featured multi-tenant tier with the Kafka endpoint, 20 consumer groups per hub, 7-day retention, and auto-inflate. Choose PREMIUM for reserved processing units, predictable latency, and extended retention (moving in or out of PREMIUM later replaces the namespace); on PREMIUM, `capacity` means processing units (1, 2, 4, 8, or 16) instead of throughput units (1-40).

## Key Outputs

| Output | Purpose |
|--------|---------|
| `namespace_id` | The parent reference for every Event Hubs child resource, the namespace-wide data-plane RBAC scope, and the private-endpoint target |
| `namespace_name` | The host SDKs connect to (`{name}.servicebus.windows.net`) and the Kafka bootstrap host on port 9093 -- Event Hubs shares the Service Bus DNS zone |
| `default_primary_connection_string` | The root SAS credential (full manage rights) -- quick starts and break-glass; mint scoped rules for production |
| `identity_principal_id` | Grant the system-assigned identity access on other resources (e.g. Storage for capture, Key Vault for CMK) |
| `default_*_connection_string_alias` | Failover-stable credentials -- populated only when the namespace carries a geo-DR pairing |

## Related Resources

- [Azure Event Hub](/docs/catalog/azure/event-hub) -- the partitioned, replayable event stream
- [Azure Event Hub Consumer Group](/docs/catalog/azure/event-hub-consumer-group) -- independent read cursors
- [Azure Event Hub Authorization Rule](/docs/catalog/azure/event-hub-authorization-rule) -- least-privilege SAS credentials
- [Azure Event Hub Schema Group](/docs/catalog/azure/event-hub-schema-group) -- the schema registry
- [Azure Event Hub Disaster Recovery Config](/docs/catalog/azure/event-hub-disaster-recovery-config) -- geo-DR pairing
- [Azure Event Hub Cluster](/docs/catalog/azure/event-hub-cluster) -- single-tenant dedicated capacity
- [Azure Event Hub Namespace Customer Managed Key](/docs/catalog/azure/event-hub-namespace-customer-managed-key) -- BYOK encryption
