# Azure Storage Queue

Creates a Storage queue inside an AzureStorageAccount -- the simple, massive-scale message buffer of Azure storage. Producers enqueue, workers poll-and-delete -- the work-queue pattern Functions queue triggers and worker pools consume.

## What Gets Created

When you deploy an AzureStorageQueue resource, Planton provisions:

- **Storage Queue** -- an `azurerm_storage_queue` on the referenced account (via its ARM id -- the control-plane path), with your metadata

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureStorageAccount** to create the queue in (referenced through `storageAccountId`)

## Quick Start

Create a file `queue.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureStorageQueue
metadata:
  name: work-items
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureStorageQueue.work-items
spec:
  storageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: my-app-storage
      fieldPath: status.outputs.storage_account_id
  queueName: work-items
```

Deploy:

```shell
planton apply -f queue.yaml
```

Messages hold up to 64 KB; a queue holds as many as the account's capacity allows. When you need pub/sub topics, sessions, ordering, or dead-lettering semantics, reach for Service Bus instead -- Storage queues win on cost, capacity, and simplicity.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `queue_id` | The ARM id data-plane role assignments (Storage Queue Data Message Sender/Processor) scope to |
| `queue_name` | What SDK clients, Functions queue triggers, and app settings reference |
| `storage_account_name` | The account/queue pair, without a second reference |

Client URLs compose from the ACCOUNT's endpoint output plus this queue's name: `{primary_queue_endpoint}{queue_name}`.

## Related Resources

- [Azure Storage Account](/docs/catalog/azure/azurestorageaccount) -- the parent account
- [Azure Role Assignment](/docs/catalog/azure/azureroleassignment) -- queue-scoped data-plane grants
