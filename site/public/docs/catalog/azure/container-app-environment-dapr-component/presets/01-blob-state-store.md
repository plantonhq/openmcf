---
title: "Blob Storage State Store"
description: "This preset registers a Dapr state store backed by Azure Blob Storage. Dapr-enabled apps whose `dapr.app_id` appears in `scopes` call the Dapr state API with the component name (`statestore`) and..."
type: "preset"
rank: "01"
presetSlug: "01-blob-state-store"
componentSlug: "container-app-environment-dapr-component"
componentTitle: "Container App Environment Dapr Component"
provider: "azure"
icon: "package"
order: 1
---

# Blob Storage State Store

This preset registers a Dapr state store backed by Azure Blob Storage. Dapr-enabled apps whose `dapr.app_id` appears in `scopes` call the Dapr state API with the component name (`statestore`) and Dapr persists state as blobs -- application code never sees a connection string.

## When to Use

- Durable key-value state for Dapr-enabled Container Apps (session data, saga state, actor state)
- Teams standardizing on Dapr's state API so the backend can change without code changes

## Key Configuration Choices

- **Component type** (`state.azure.blobstorage`) -- Any Dapr state component works the same way (state.redis, state.azure.cosmosdb, ...); the metadata keys change per backend
- **Secret-backed account key** (`metadata[].secretName` -> `secrets[]`) -- The access key travels as a component secret, never a plain metadata value
- **Deliberate scoping** (`scopes`) -- Only the listed dapr app ids see the component; an empty list would expose it to every Dapr-enabled app in the environment

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<container-app-environment-id>` | ARM ID of the Container App Environment | `AzureContainerAppEnvironment` status outputs |
| `<storage-account-name>` | The storage account backing the state store | `AzureStorageAccount` status outputs (storage_account_name) |
| `<storage-account-access-key>` | Account access key Dapr uses | `AzureStorageAccount` status outputs (primary_access_key) |
| `<dapr-app-id>` | The consuming app's dapr.app_id value | The `AzureContainerApp`'s `dapr.app_id` field |

## Related Presets

- **02-servicebus-pubsub** -- The pub/sub building block on Azure Service Bus
