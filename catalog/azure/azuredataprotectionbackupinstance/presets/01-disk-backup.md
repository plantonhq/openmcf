# Disk Backup

This preset protects one managed disk with the modern Data Protection service: scheduled incremental snapshots, kept in a dedicated snapshot resource group, retained per the referenced disk policy.

## When to Use

- Stateful workloads on managed disks (databases on VMs, self-managed message brokers) that need point-in-time disk recovery
- Faster, cheaper disk-level protection than whole-VM backup when only the data disk matters
- Estates standardizing on the modern Data Protection vault over the classic Recovery Services vault

## Key Configuration Choices

- **A dedicated snapshot resource group** -- snapshots live OUTSIDE the disk's own group, so deleting the app's group never deletes its recovery points
- **The policy's variant must be `disk`** -- Azure rejects a policy of any other datasource type
- **Grants precede the instance** -- the vault identity needs "Disk Backup Reader" on the disk and "Disk Snapshot Contributor" on the snapshot group; compose AzureRoleAssignment resources referencing the vault's `system_assigned_identity_principal_id` output

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-backup-vault>` | The AzureDataProtectionBackupVault holding the backups | The vault component's name |
| `<your-disk-policy>` | An AzureDataProtectionBackupPolicy with the `disk` variant, on the same vault | The policy component's name |
| `<your-managed-disk>` | The AzureManagedDisk being protected | The disk component's name |
| `<your-snapshot-resource-group>` | The AzureResourceGroup where snapshots are stored | The resource-group component's name |

The instance is free; snapshot storage bills per the policy's retention.
