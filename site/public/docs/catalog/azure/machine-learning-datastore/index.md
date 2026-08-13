---
title: "Machine Learning Datastore"
description: "Machine Learning Datastore deployment documentation"
icon: "package"
order: 100
componentName: "azuremachinelearningdatastore"
---

# Azure Machine Learning Datastore

Registers a datastore on an Azure Machine Learning workspace -- the saved connection that tells the workspace where data lives (a blob container, a Data Lake Gen2 filesystem, or an Azure Files share) and how to reach it. Exactly one variant block selects the flavor. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Datastore** -- an ARM child of the workspace (`.../workspaces/{ws}/dataStores/{name}`); the variant block decides whether it is a blob-container, data-lake, or file-share connection

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureMachineLearningWorkspace** -- the datastore is registered on it.
- **The storage target** -- an AzureStorageContainer, AzureStorageDataLakeGen2Filesystem, or AzureStorageShare to point at.

### Azure Subscription

- **Credentials are write-only** -- ARM never returns keys, SAS tokens, or client secrets; reference secrets rather than embedding literals.
- **Nearly everything is fixed at creation** -- the name, storage target, description, and tags all replace the datastore on change; only credentials and the service-data identity update in place.

## Deploy

### Console

Open the deployment store, find **Azure Machine Learning Datastore**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Blob Training Data** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMachineLearningDatastore
metadata:
  name: training-data
  org: acme-corp
  env: prod
spec:
  workspaceId:
    valueFrom:
      name: ml-prod
  name: training_data
  serviceDataIdentity: WORKSPACE_SYSTEM_ASSIGNED_IDENTITY
  blobStorage:
    storageContainerId:
      valueFrom:
        name: training-data-container
```

```shell
planton apply -f azure-machine-learning-datastore.yaml
```

The datastore registers in seconds.

### InfraChart

In an ML-platform chart the order is: workspace → storage container/share → **datastore**, wiring both by reference; jobs then reference the datastore by its name.

## Key Configuration

These are the most important decisions when configuring the datastore. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The variant block IS the type** -- set exactly one of `blobStorage`, `dataLakeGen2`, or `fileShare`. Each carries its own storage target reference and its own auth grammar.

**Auth strategy** -- the hardened posture is `serviceDataIdentity: WORKSPACE_SYSTEM_ASSIGNED_IDENTITY` with no embedded credentials (grant the workspace identity data access on the storage target). Key/SAS auth works everywhere but embeds a secret reference in the manifest; the FILE SHARE variant always requires exactly one credential (the provider's own contract).

**Default datastore** -- only the blob variant can claim `isDefault: true` (where job outputs land unless directed elsewhere).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureMachineLearningWorkspace** | `workspaceId` | `status.outputs.machine_learning_workspace_id` |
| **AzureStorageContainer** | `blobStorage.storageContainerId` | `status.outputs.container_id` |
| **AzureStorageDataLakeGen2Filesystem** | `dataLakeGen2.storageContainerId` | `status.outputs.filesystem_id` |
| **AzureStorageShare** | `fileShare.storageFileshareId` | `status.outputs.share_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `datastore_id` | ARM ID of the datastore | Operational tooling |
| `datastore_name` | The datastore's name | What jobs and data assets reference |
| `is_default` | Whether this is the workspace's default datastore | Chart logic |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Blob training data** -- a container datastore under workspace identity. Start from the **Blob Training Data** preset.

**Data lake filesystem** -- a Gen2 filesystem with service-principal auth. Start from the **Data Lake Filesystem** preset.

## Works With

- [**Azure Machine Learning Workspace**](/cloud-catalog/azure-machine-learning-workspace) -- the parent workspace
- [**Azure Storage Container**](/cloud-catalog/azure-storage-container) -- the blob variant's target
- [**Azure Storage Data Lake Gen2 Filesystem**](/cloud-catalog/azure-storage-data-lake-gen2-filesystem) -- the data-lake variant's target
- [**Azure Storage Share**](/cloud-catalog/azure-storage-share) -- the file-share variant's target
