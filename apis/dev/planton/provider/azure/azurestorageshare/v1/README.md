# AzureStorageShare

An Azure Files share inside an AzureStorageAccount: the SMB/NFS file
system unit of Azure storage. VMs, AKS pods, and container apps mount
shares for shared POSIX-style state -- lift-and-shift app data, user
profiles, CI caches, shared content -- and Azure bills, throttles,
tiers, and snapshots at the share level.

## When to Use

Use AzureStorageShare when you need:

- **A mountable shared file system** -- SMB for Windows and most Linux
  workloads, NFS (premium) for POSIX semantics
- **Lift-and-shift file storage** -- apps that expect a drive letter or
  mount point, not an object API
- **Kubernetes shared volumes** -- the Azure Files CSI driver mounts
  shares as ReadWriteMany volumes
- **A role-assignment scope** -- grant Storage File Data SMB Share
  Reader/Contributor on `rbac_scope_id` for share-level access

## Key Configuration

- `storage_account_id` -- the parent account, referenced from an
  AzureStorageAccount's output (fixed at creation); the account's kind
  gates NFS and the PREMIUM tier (FileStorage) and quotas above 5120 GB
  (large_file_share_enabled)
- `share_name` -- 3-63 lowercase letters/digits/hyphens, unique within
  the account; becomes the mount path segment
- `quota_gb` -- the provisioned maximum size (1-102400 GB); premium
  bills it whether used or not
- `enabled_protocol` -- SMB (default) or NFS; fixed at creation
- `access_tier` -- TransactionOptimized/Hot/Cool (standard) or Premium
  (FileStorage); unset lets Azure pick the account-kind default
- `acls` -- stored access policies anchoring revocable SAS tokens

## Composition

```yaml
storageAccountId:
  valueFrom:
    kind: AzureStorageAccount
    name: app-storage
    fieldPath: status.outputs.storage_account_id
```

Mount paths compose from the ACCOUNT's endpoint plus this share's name:
`{primary_file_endpoint}{share_name}`.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
