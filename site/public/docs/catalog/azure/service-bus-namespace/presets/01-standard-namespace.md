---
title: "Standard Namespace"
description: "This preset creates a STANDARD-tier Service Bus namespace -- the full-featured multi-tenant tier (queues, topics, subscriptions, sessions, duplicate detection) that fits most production workloads."
type: "preset"
rank: "01"
presetSlug: "01-standard-namespace"
componentSlug: "service-bus-namespace"
componentTitle: "Service Bus Namespace"
provider: "azure"
icon: "package"
order: 1
---

# Standard Namespace

This preset creates a STANDARD-tier Service Bus namespace -- the
full-featured multi-tenant tier (queues, topics, subscriptions,
sessions, duplicate detection) that fits most production workloads.

## When to Use

- The default starting point for any messaging domain
- Workloads that don't need VNet isolation, CMK, geo-DR, or messages
  beyond 256 KB (those are PREMIUM features)

## Key Configuration Choices

- **`sku: STANDARD`** -- multi-tenant with the full feature set;
  BASIC↔STANDARD changes update in place, but moving to PREMIUM later
  replaces the namespace
- **`tags`** -- ARM tags are Azure's governance surface; user tags merge
  over the platform's identity tags

## Values to Customize

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `eastus` | The namespace's region (e.g. eastus) | Your region strategy |
| `my-messaging-rg` | The AzureResourceGroup's Planton resource name | Your foundation composition |
| `myorg-app-bus` | 6-50 chars, letters/numbers/hyphens, globally unique | Your naming convention |
| `order-processing` | What this namespace carries | Your data taxonomy |

## Downstream Wiring

Entities reference the namespace's ARM id:

```yaml
# On an AzureServiceBusQueue
namespaceId:
  valueFrom:
    kind: AzureServiceBusNamespace
    name: my-app-bus
    fieldPath: status.outputs.namespace_id
```
