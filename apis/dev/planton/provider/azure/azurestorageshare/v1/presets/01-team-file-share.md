# Team File Share

This preset creates a standard SMB file share -- the general-purpose
shared drive: Windows mounts it natively, Linux mounts it via cifs, and
Azure's default TransactionOptimized tier balances storage and
operation costs.

## When to Use

- Shared team or application file storage on a standard account
- Lift-and-shift apps that expect a mounted drive, not an object API
- Kubernetes ReadWriteMany volumes via the Azure Files CSI driver

## Key Configuration Choices

- **`quotaGb: 500`** -- the provisioned ceiling; grows in place later.
  Standard accounts cap at 5120 GB unless the account enables
  `largeFileShareEnabled`
- **`enabledProtocol: SMB`** -- fixed at creation; works on every
  account kind (NFS would require a premium FileStorage account)
- **Tier left unset** -- Azure applies TransactionOptimized on standard
  accounts; set `accessTier: COOL` for rarely-read archival shares

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `team-files` | 3-63 lowercase letters/digits/hyphens | Your naming convention |
| `<data-domain>` | What lives in this share | Your data taxonomy |

## Downstream Wiring

Scope a data-plane grant to just this share (note: Azure Files grants
target `rbac_scope_id`, not the management id):

```yaml
# On an AzureRoleAssignment
scope:
  valueFrom:
    kind: AzureStorageShare
    name: my-team-files
    fieldPath: status.outputs.rbac_scope_id
roleDefinitionName: Storage File Data SMB Share Contributor
```
