# Incremental Disk Backup

This preset creates an incremental Copy-mode snapshot of a managed disk -- one deliberate point-in-time copy that stores only the delta since the disk's previous snapshot, on standard storage.

## When to Use

- A pre-change backup before risky maintenance (delete it once the change settles)
- The image pipeline's handoff artifact (snapshot the prepared OS disk, then publish it as a gallery image version)
- Any single deliberate copy -- for SCHEDULED, retention-managed backups use the Recovery Services kinds instead

## Key Configuration Choices

- **`createOption: Copy` + `sourceResourceId`** -- the working pair for disk sources; Azure (not the schema) validates the pairing, so get it right in the manifest
- **`incrementalEnabled: true`** -- the backup-chain default; fixed at creation, and the first incremental of a disk stores the full disk
- **Public network defaults** -- fine for build artifacts; for sensitive disks add `networkAccessPolicy: AllowPrivate`/`DenyAll` and turn public access off

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `<your-managed-disk>` | The Planton name of the `AzureManagedDisk` to snapshot | Planton console (or replace `valueFrom` with `value:` and a literal disk ARM ID) |
| `app-disk-snap` | The snapshot's name (up to 80 letters, numbers, dashes, underscores -- no dots) | Your naming convention |
| `eastus` | The Azure region -- create the snapshot where its source lives | The source disk's region |

## Related Presets

- **Import VHD** -- creates the snapshot from a VHD blob instead of a disk
