# Azure Storage Share

Creates an Azure Files share inside an AzureStorageAccount -- the SMB/NFS file system unit of Azure storage. VMs, AKS pods, and container apps mount shares for shared POSIX-style state; quota, protocol, tier, and stored access policies are set per share.

## What Gets Created

When you deploy an AzureStorageShare resource, Planton provisions:

- **Azure Files Share** -- an `azurerm_storage_share` on the referenced account (via its ARM id -- the control-plane path), with your chosen quota, protocol, tier, stored access policies, and metadata

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureStorageAccount** to create the share in (referenced through `storageAccountId`); NFS shares and the PREMIUM tier need a FileStorage (premium file) account

## Quick Start

Create a file `share.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureStorageShare
metadata:
  name: team-files
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureStorageShare.team-files
spec:
  storageAccountId:
    valueFrom:
      kind: AzureStorageAccount
      name: my-app-storage
      fieldPath: status.outputs.storage_account_id
  shareName: team-files
  quotaGb: 500
```

Deploy:

```shell
planton apply -f share.yaml
```

The share speaks SMB by default -- what Windows mounts natively and Linux mounts via cifs. For Linux workloads that need POSIX semantics, set `enabledProtocol: NFS` (requires a premium FileStorage account).

## Key Outputs

| Output | Purpose |
|--------|---------|
| `share_id` | The management ARM id -- what ARM reads and policy target |
| `rbac_scope_id` | The scope Azure Files data-plane role assignments (Storage File Data SMB Share Reader/Contributor) target -- a different segment than the management id |
| `share_name` | What mount commands, CSI volume definitions, and app settings reference |
| `storage_account_name` | The account/share pair, without a second reference |

Mount paths compose from the ACCOUNT's endpoint output plus this share's name: `{primary_file_endpoint}{share_name}`.

## Related Resources

- [Azure Storage Account](/docs/catalog/azure/azurestorageaccount) -- the parent account
- [Azure Role Assignment](/docs/catalog/azure/azureroleassignment) -- share-scoped data-plane grants
