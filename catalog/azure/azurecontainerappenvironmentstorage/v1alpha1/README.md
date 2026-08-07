# AzureContainerAppEnvironmentStorage

Register an Azure Files share as a named storage resource on a Container App Environment, so apps and jobs in the environment can mount it as a persistent volume.

## Overview

Container Apps and Jobs cannot mount file shares directly. The share is first registered on the environment as a storage resource, and workloads then declare `AZURE_FILE` / `NFS_AZURE_FILE` volumes referencing the registration by `storage_name`. One registration can back volumes in any number of apps and jobs in the environment.

## Key Features

- **Two share protocols**: SMB (account name + access key -- the common case) or NFS (the account's file endpoint; requires a VNet-injected environment)
- **Access mode control**: READ_ONLY for shared configuration and reference data, READ_WRITE for working storage
- **In-place key rotation**: the SMB `access_key` is the one updatable field -- rotate storage keys without recreating the registration
- **Composable**: `share_name`, `account_name`, and `access_key` all default to `AzureStorageShare` / `AzureStorageAccount` outputs

## When to Use

- Persistent working storage shared across replicas of an app (upload staging, caches that must survive restarts)
- Distributing reference data or shared configuration to many apps read-only
- Batch jobs that hand work products to other workloads through a shared file system

## Spec Highlights

| Field | Notes |
| --- | --- |
| `storage_name` | The handle app/job volumes reference. Max 32 lowercase alphanumerics/hyphens. ForceNew |
| `share_name` | The Azure Files share; defaults to an `AzureStorageShare` reference. ForceNew |
| `access_mode` | READ_ONLY or READ_WRITE. ForceNew |
| `account_name` + `access_key` | The SMB path (exactly one protocol per registration) |
| `nfs_server_url` | The NFS path (`{account}.file.core.windows.net`); requires VNet injection |

## Outputs

| Output | Purpose |
| --- | --- |
| `storage_id` | The registration's ARM ID |
| `storage_name` | The composition seam -- what app and job volumes reference |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureContainerAppEnvironmentStorage
metadata:
  name: app-data
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

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
