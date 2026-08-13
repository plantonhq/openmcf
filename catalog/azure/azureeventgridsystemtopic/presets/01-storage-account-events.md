# Storage Account Events

This preset creates a system topic on a storage account -- the most common source, exposing BlobCreated/BlobDeleted (and queue) events for pipelines to subscribe to.

## When to Use

- Blob-triggered processing (ingest pipelines, thumbnail generation, virus scanning hand-offs)
- Auditing object lifecycle without polling the account

## Key Configuration Choices

- **Topic type `Microsoft.Storage.StorageAccounts`** -- the storage service's event catalog (blob and queue events; table events do not exist)
- **Region follows the source** -- the topic must sit in the storage account's region; the value here is the account's, not a free choice
- **One topic per account** -- Azure allows a single system topic per source per type; teams sharing the account share this topic and attach their own subscriptions

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `<your-storage-account>` | The Planton name of your `AzureStorageAccount` resource | Planton console (or replace `valueFrom` with `value:` and the account's ARM ID) |
| `eastus` | The storage account's own region | The account's overview blade |
