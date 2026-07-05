---
title: "Work Queue"
description: "This preset creates a plain work queue -- the decoupling primitive: producers enqueue up to 64 KB messages, a worker pool polls and deletes them at its own pace."
type: "preset"
rank: "01"
presetSlug: "01-work-queue"
componentSlug: "storage-queue"
componentTitle: "Storage Queue"
provider: "azure"
icon: "package"
order: 1
---

# Work Queue

This preset creates a plain work queue -- the decoupling primitive:
producers enqueue up to 64 KB messages, a worker pool polls and deletes
them at its own pace.

## When to Use

- Background job dispatch (thumbnails, exports, notifications)
- Load leveling between a bursty producer and steady workers
- Any Functions queue-trigger source

## Key Configuration Choices

- **The queue is deliberately minimal** -- visibility timeouts, batch
  sizes, and retry behavior are data-plane choices producers and
  consumers make at runtime, not infrastructure
- **`metadata.purpose`** -- a self-documenting marker for operators
  browsing the account

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `<queue-name>` | 3-63 lowercase letters/digits/hyphens | Your naming convention |
| `<work-type>` | What flows through this queue | Your job taxonomy |

## Downstream Wiring

Split producer and consumer permissions at the queue boundary:

```yaml
# Producer grant (send only)
scope:
  valueFrom:
    kind: AzureStorageQueue
    name: my-work-items
    fieldPath: status.outputs.queue_id
roleDefinitionName: Storage Queue Data Message Sender
```
