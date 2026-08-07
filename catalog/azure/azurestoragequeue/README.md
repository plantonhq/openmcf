# AzureStorageQueue

A Storage queue inside an AzureStorageAccount: the simple, massive-scale
message buffer of Azure storage. Producers enqueue up to 64 KB messages
and workers poll-and-delete them -- the classic work-queue /
load-leveling pattern Functions queue triggers, WebJobs, and worker
pools consume.

## When to Use

Use AzureStorageQueue when you need:

- **A work queue** -- decouple producers from a worker pool at
  storage-account cost and capacity
- **A Functions trigger source** -- the queue trigger is the canonical
  serverless work-dispatch pattern
- **Load leveling** -- absorb bursts and let workers drain at their own
  pace
- **A role-assignment scope** -- grant Storage Queue Data Message
  Sender/Processor on `queue_id` for least-privilege producer/consumer
  splits

For pub/sub topics, sessions, ordering guarantees, or dead-lettering
semantics, Service Bus is the richer sibling; Storage queues win on
cost, capacity, and operational simplicity.

## Key Configuration

- `storage_account_id` -- the parent account, referenced from an
  AzureStorageAccount's output (fixed at creation)
- `queue_name` -- 3-63 lowercase letters/digits/hyphens, unique within
  the account; becomes the URL path segment
- `metadata` -- free-form key/value pairs on the queue (not secrets;
  not Azure tags)

## Composition

```yaml
storageAccountId:
  valueFrom:
    kind: AzureStorageAccount
    name: app-storage
    fieldPath: status.outputs.storage_account_id
```

Client URLs compose from the ACCOUNT's endpoint plus this queue's name:
`{primary_queue_endpoint}{queue_name}`.

## Documentation

- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
