# Standard Backup Vault

This preset creates the everyday production vault: geo-redundant backup storage with cross-region restore, Microsoft-managed encryption, all alert switches at their all-on defaults. The right starting point for a region's VM and file-share backups.

## When to Use

- The first vault a region's environment needs -- one home for its backup protections
- Standard production posture without customer-managed-key requirements
- Anywhere the paired-region restore capability is worth its modest storage premium

## Key Configuration Choices

- **`storageModeType: GeoRedundant`** -- a copy in the paired region; decide before the first protection (redundancy locks once items are protected)
- **`crossRegionRestoreEnabled: true`** -- restore in the paired region during an outage; note disabling later REPLACES the vault (one-way)
- **No `monitoring` block** -- all five alert switches default ON service-side; configure the block only to turn specific classes off
- **No immutability yet** -- run `Unlocked` once retention settings have settled (see the second preset)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The AzureResourceGroup the vault lives in | Your resource group resource's name |

The vault is free at rest -- cost starts when items are protected.
