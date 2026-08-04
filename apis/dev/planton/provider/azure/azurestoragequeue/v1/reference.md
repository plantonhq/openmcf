# AzureStorageQueue

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureStorageQueueSpec** defines the configuration for creating a
Storage queue inside an Azure Storage Account: the simple, massive-
scale message buffer of Azure storage. Producers enqueue up to 64 KB
messages and workers poll-and-delete them -- the classic work-queue /
load-leveling pattern Functions queue triggers, WebJobs, and worker
pools consume. (For pub/sub topics, sessions, ordering guarantees, or
dead-lettering semantics, Service Bus is the richer sibling; Storage
queues win on cost, capacity, and operational simplicity.)

Queues are many-per-account with independent lifecycles, which is why
they are a first-class kind referencing the account rather than a list
folded into the account's spec. The parent is fixed at creation: a
queue cannot move between accounts.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureStorageQueue
metadata:
  name: test-storage-queue
spec:
  storageAccountId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/plantonhackstorage
  queueName: hack-work-items
  # Exercises the metadata pass-through.
  metadata:
    purpose: hack-validation
    consumer: worker-pool
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.queueName` | `string` | yes |  |  |
| `spec.metadata` | `map<string, string>` |  |  |  |

## Field Details

### spec.storageAccountId

`string | valueFrom` · required

The storage account the queue lives in, by ARM ID. References an
AzureStorageAccount's storage_account_id output so the account and
its queues compose in one manifest set. Fixed at creation.

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.queueName

`string` · required

The queue's name: 3-63 lowercase letters, digits, and hyphens;
starts and ends with a letter or digit; no consecutive hyphens.
Unique within the account (it becomes the URL path segment:
https://{account}.queue.core.windows.net/{name}). Changing the name
replaces the queue.

- rule: queue_name must be 3-63 lowercase letters, digits, and hyphens, starting and ending with a letter or digit, with no consecutive hyphens
- rule: {"required":true,"string":{"minLen":"3","maxLen":"63"}}

### spec.metadata

`map<string, string>`

Free-form metadata key/value pairs stored on the queue (visible to
anyone who can read queue properties -- not for secrets). Keys must
be valid C# identifiers per Azure's rule; lowercase is canonical
(Azure lowercases keys on read).

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureStorageQueue, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.queue_id` | `string` | The Azure Resource Manager ID of the queue. Role assignments (Storage Queue Data Contributor/Message Processor/Message Sender) scope to it for queue-level data access. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Storage/storageAccounts/{account}/queueServices/default/queues/{name} |
| `status.outputs.queue_name` | `string` | The queue's name -- what SDK clients, Functions queue triggers, and app settings reference within the account. |
| `status.outputs.storage_account_name` | `string` | The name of the storage account the queue lives in, parsed from the resolved account ID -- saves consumers a second reference when they need the account/queue pair. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
