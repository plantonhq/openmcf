---
title: "Service Bus Namespace"
description: "Service Bus Namespace deployment documentation"
icon: "package"
order: 100
componentName: "azureservicebusnamespace"
---

# Azure Service Bus Namespace

Creates an Azure Service Bus namespace -- the container and billing boundary for enterprise messaging. The namespace sets the pricing tier, network posture, encryption, and authentication mode; queues, topics, subscriptions, SAS credentials, and geo-DR pairings compose onto it as first-class resources.

## What Gets Created

When you deploy an AzureServiceBusNamespace resource, Planton provisions:

- **Service Bus Namespace** -- an `azurerm_servicebus_namespace` in the referenced resource group, with your chosen tier, optional managed identity, customer-managed-key encryption, and network rules

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureResourceGroup** to create the namespace in (referenced through `resourceGroup`)

## Quick Start

Create a file `namespace.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServiceBusNamespace
metadata:
  name: orders-bus
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureServiceBusNamespace.orders-bus
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: messaging-rg
      fieldPath: status.outputs.resource_group_name
  namespaceName: myorg-orders-bus
```

Deploy:

```shell
planton apply -f namespace.yaml
```

Unset `sku` deploys STANDARD -- the full-featured multi-tenant tier. Choose PREMIUM for dedicated capacity, VNet integration, CMK, geo-DR, or messages beyond 256 KB (set `capacity` and `premiumMessagingPartitions` with it).

## Key Outputs

| Output | Purpose |
|--------|---------|
| `namespace_id` | The parent reference for every Service Bus child resource, the namespace-wide data-plane RBAC scope, and the private-endpoint target |
| `endpoint` | What Service Bus SDKs connect to |
| `default_primary_connection_string` | The root SAS credential (full manage rights) -- quick starts and break-glass; mint scoped rules for production |
| `identity_principal_id` | Grant the system-assigned identity access on other resources (e.g. Key Vault for CMK) |

## Related Resources

- [Azure Service Bus Queue](/docs/catalog/azure/service-bus-queue) -- point-to-point messaging
- [Azure Service Bus Topic](/docs/catalog/azure/service-bus-topic) -- publish-subscribe
- [Azure Service Bus Authorization Rule](/docs/catalog/azure/service-bus-authorization-rule) -- least-privilege SAS credentials
- [Azure Service Bus Disaster Recovery Config](/docs/catalog/azure/service-bus-disaster-recovery-config) -- geo-DR pairing
