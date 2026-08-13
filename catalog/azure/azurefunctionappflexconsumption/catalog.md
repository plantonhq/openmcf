# Azure Function App Flex Consumption

Deploys an Azure Function App on the Flex Consumption plan -- Azure's newest serverless Functions hosting model, with per-instance memory selection, a configurable scale-out ceiling, and always-ready instance pools that eliminate cold starts. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Flex Consumption Function App** -- the Microsoft.Web site with its runtime declaration, deployment-storage binding, scale configuration (instance memory, fan-out ceiling, HTTP concurrency, always-ready pools), site configuration, app settings, connection strings, managed identity, Easy Auth v2, and tags

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An FC1 service plan** -- an `AzureServicePlan` with `skuName: FLEX_CONSUMPTION_FC1`; Azure rejects flex apps on any other tier. One FC1 plan hosts many flex apps and has no idle compute cost.
- **A blob container for deployment storage** -- an `AzureStorageAccount` plus an `AzureStorageContainer`; the `storageContainerEndpoint` composes the account's blob endpoint with the container name.

### Azure Subscription

- **Regional availability** -- Flex Consumption is not offered in every region; the plan and app must share a region where it is.
- **Pick the storage authentication mode first** -- connection-string auth is the simplest (the key travels in the manifest as a reference); managed-identity auth is credential-free but needs a "Storage Blob Data Contributor" grant on the storage account before package deployments work.
- **Billing follows executions and always-ready instances** -- an idle app with no warm pools costs nothing.

## Deploy

### Console

Open the deployment store, find **Azure Function App Flex Consumption**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Node HTTP API** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f flex-function-app.yaml
```

## After Deploy

Point DNS or an upstream proxy at the `default_hostname` output. Grant the identity outputs (`identity_principal_id`) the RBAC roles the functions need -- including "Storage Blob Data Contributor" on the deployment storage account when using identity-based storage auth. Deploy function code through your CI/CD pipeline (Azure Functions Core Tools, GitHub Actions, or Azure DevOps) targeting the app.
