---
title: "Container App Environment Storage"
description: "Container App Environment Storage deployment documentation"
icon: "package"
order: 100
componentName: "azurecontainerappenvironmentstorage"
---

# Azure Container App Environment Storage

Registers an Azure Files share as a named storage resource on a Container App Environment -- the bridge that lets Container Apps and Jobs mount persistent file shares as volumes.

## What Gets Created

When you deploy an AzureContainerAppEnvironmentStorage resource, Planton provisions:

- **Environment storage registration** -- an `azurerm_container_app_environment_storage` on the referenced environment, addressing the share over SMB (account name + access key) or NFS (file endpoint)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureContainerAppEnvironment** to register the storage on (referenced through `containerAppEnvironmentId`)
- **An Azure Files share** to register (an `AzureStorageShare` on an `AzureStorageAccount`); NFS shares additionally require the environment to be VNet-injected

## Quick Start

Create a file `storage.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureContainerAppEnvironmentStorage
metadata:
  name: app-data
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureContainerAppEnvironmentStorage.app-data
spec:
  containerAppEnvironmentId:
    valueFrom:
      kind: AzureContainerAppEnvironment
      name: my-env
      fieldPath: status.outputs.environment_id
  storageName: app-data
  shareName:
    valueFrom:
      kind: AzureStorageShare
      name: my-share
      fieldPath: status.outputs.share_name
  accessMode: READ_WRITE
  accountName:
    valueFrom:
      kind: AzureStorageAccount
      name: my-account
      fieldPath: status.outputs.storage_account_name
  accessKey:
    valueFrom:
      kind: AzureStorageAccount
      name: my-account
      fieldPath: status.outputs.primary_access_key
```

Deploy:

```shell
planton apply -f storage.yaml
```

Apps and jobs then mount the registration through a volume:

```yaml
volumes:
  - name: data
    storageType: AZURE_FILE
    storageName:
      valueFrom:
        kind: AzureContainerAppEnvironmentStorage
        name: app-data
        fieldPath: status.outputs.storage_name
```

Everything except the SMB `access_key` is ForceNew -- key rotation is the one in-place update; any other change recreates the registration and briefly breaks volume mounts that reference it.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `storage_id` | The registration's ARM ID |
| `storage_name` | What app and job volumes reference in `storage_name` |

## Related Resources

- [Azure Container App Environment](/docs/catalog/azure/container-app-environment) -- the environment the registration lives on
- [Azure Container App](/docs/catalog/azure/container-app) / [Azure Container App Job](/docs/catalog/azure/container-app-job) -- the workloads that mount it
- [Azure Storage Share](/docs/catalog/azure/storage-share) -- the Azure Files share being registered
