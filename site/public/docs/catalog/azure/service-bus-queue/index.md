---
title: "Service Bus Queue"
description: "Service Bus Queue deployment documentation"
icon: "package"
order: 100
componentName: "azureservicebusqueue"
---

# Azure Service Bus Queue

Creates a queue inside an Azure Service Bus namespace -- reliable point-to-point messaging with FIFO delivery, sessions, duplicate detection, and a built-in dead-letter sub-queue.

## What Gets Created

When you deploy an AzureServiceBusQueue resource, Planton provisions:

- **Service Bus Queue** -- an `azurerm_servicebus_queue` on the referenced namespace (via its ARM id), with your chosen size, lock, delivery, session, and dead-letter dials

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureServiceBusNamespace** to create the queue in (referenced through `namespaceId`)

## Quick Start

Create a file `queue.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureServiceBusQueue
metadata:
  name: orders-queue
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureServiceBusQueue.orders-queue
spec:
  namespaceId:
    valueFrom:
      kind: AzureServiceBusNamespace
      name: orders-bus
      fieldPath: status.outputs.namespace_id
  queueName: orders
```

Deploy:

```shell
planton apply -f queue.yaml
```

Unset dials keep Azure's defaults: 1 GB size (multi-tenant), 1-minute lock, 10 delivery attempts, unbounded TTL. Three dials are fixed at creation -- `partitioningEnabled`, `requiresDuplicateDetection`, and `requiresSession` -- so decide them up front.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `queue_id` | The scope for queue-level data-plane role assignments and the parent for queue-scoped SAS rules |
| `queue_name` | What SDK clients, connection strings, and function bindings reference |
| `namespace_name` | The namespace/queue pair, without a second reference |

## Related Resources

- [Azure Service Bus Namespace](/docs/catalog/azure/service-bus-namespace) -- the parent namespace
- [Azure Service Bus Authorization Rule](/docs/catalog/azure/service-bus-authorization-rule) -- queue-scoped SAS credentials
- [Azure Role Assignment](/docs/catalog/azure/role-assignment) -- keyless data-plane grants on `queue_id`
