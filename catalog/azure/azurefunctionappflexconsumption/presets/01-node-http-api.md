# Node HTTP API

This preset deploys a Node.js function app on the Flex Consumption plan with one always-ready instance for HTTP triggers -- an event-driven API that answers instantly on the hot path and bills per execution everywhere else.

## When to Use

- HTTP-triggered APIs and webhooks where the first request must not cold-start
- Event-driven backends that idle most of the day and spike unpredictably
- Teams moving off classic Consumption for the per-instance memory and concurrency dials

## Key Configuration Choices

- **`alwaysReady: http × 1`** -- one warm instance serves every HTTP trigger; it is the app's only idle cost (drop the block for pure pay-per-execution)
- **`maximumInstanceCount: 100`** -- the fan-out ceiling is the cost lever for spikes; Azure enforces that always-ready counts stay within it
- **`STORAGE_ACCOUNT_CONNECTION_STRING`** -- the simplest deployment-storage mode; the key references the storage account's output and travels no further (Azure manages the derived connection string as a hidden app setting)
- **The health-check pair** -- path and eviction time require each other; the eviction time travels as an app setting Azure manages for you

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console |
| `<your-flex-service-plan>` | The Planton name of your FC1-SKU `AzureServicePlan` | Planton console |
| `<your-storage-account>` | The Planton name of the `AzureStorageAccount` holding the deployment container | Planton console |

## Related Presets

- **Identity-Secured Worker** -- the credential-free storage-auth shape
- **Entra-Protected API** -- the platform-authenticated shape
