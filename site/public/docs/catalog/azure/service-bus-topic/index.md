---
title: "Service Bus Topic"
description: "Service Bus Topic deployment documentation"
icon: "package"
order: 100
componentName: "azureservicebustopic"
---

# Azure Service Bus Topic

Creates a topic inside an Azure Service Bus namespace -- the publish-subscribe primitive. Publishers send to the topic; each subscription under it receives an independent, filtered copy of the stream.

## What Gets Created

When you deploy an AzureServiceBusTopic resource, Planton provisions:

- **Service Bus Topic** -- an `azurerm_servicebus_topic` on the referenced namespace (via its ARM id), with your chosen size, duplicate-detection, and ordering dials

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureServiceBusNamespace** on STANDARD or PREMIUM (BASIC is queue-only)

## Quick Start

Create a file `topic.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureServiceBusTopic
metadata:
  name: order-events
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureServiceBusTopic.order-events
spec:
  namespaceId:
    valueFrom:
      kind: AzureServiceBusNamespace
      name: events-bus
      fieldPath: status.outputs.namespace_id
  topicName: order-events
```

Deploy:

```shell
planton apply -f topic.yaml
```

A topic without subscriptions delivers nothing -- compose AzureServiceBusSubscription resources under it, one per consuming application. Consumer semantics (locks, delivery counts, sessions) live on the subscription.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `topic_id` | The parent reference for subscriptions and topic-scoped SAS rules; the scope for topic-level data-plane role assignments |
| `topic_name` | What SDK clients and function bindings reference |
| `namespace_name` | The namespace/topic pair, without a second reference |

## Related Resources

- [Azure Service Bus Namespace](/docs/catalog/azure/service-bus-namespace) -- the parent namespace
- [Azure Service Bus Subscription](/docs/catalog/azure/service-bus-subscription) -- the consuming views
- [Azure Service Bus Authorization Rule](/docs/catalog/azure/service-bus-authorization-rule) -- topic-scoped SAS credentials
