# Image With Version

This preset creates the image definition AND publishes its first release: version `1.0.0`, built from an OS disk snapshot, replicated to one region. The versions list is the authoritative release history -- add entries to publish, remove them to unpublish.

## When to Use

- The image pipeline's first release is ready (a prepared OS disk has been snapshotted)
- Migrating an existing image's releases into declarative management

## Key Configuration Choices

- **Semver version names** (`1.0.0`) -- three dot-separated numeric segments; the name is the version's permanent identity under the image
- **Exactly one source per version** -- this preset uses a disk snapshot (the chart-native chain: managed disk -> AzureDiskSnapshot -> version); VHD blobs and managed images/VMs are the alternatives
- **`regionalReplicaCount: 1`** -- right for trickle deployments; one replica serves ~20 concurrent VM creations
- **Full replication** (the default) -- durable per-region copies; reserve Shallow for the dev inner loop

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `<your-compute-gallery>` | The Planton name of your `AzureComputeGallery` resource | Planton console (or replace `valueFrom` with `value:` and a literal gallery name) |
| `<your-disk-snapshot>` | The Planton name of the `AzureDiskSnapshot` holding the prepared OS disk | Planton console (or replace `valueFrom` with `value:` and a literal snapshot ARM ID) |
| `myorg` / `ubuntu` / `22-04-lts-gen2` | The image's permanent publisher/offer/SKU identity | Your naming convention |

## Related Presets

- **Linux Gen2 Trusted Launch** -- the definition alone, before the first release lands
