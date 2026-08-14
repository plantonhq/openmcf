# Linux Gen2 Trusted Launch

This preset creates a Gen2 Linux image definition with trusted launch SUPPORTED -- the identity and security posture that age best for new golden-image libraries. It publishes no versions yet; releases are added to the `versions` list as the image pipeline produces them.

## When to Use

- Registering a new golden image's identity before its first build lands
- Any Linux image whose consumers should be able to CHOOSE trusted launch

## Key Configuration Choices

- **`hyperVGeneration: V2` + `trustedLaunchSupported: true`** -- consumers may deploy with or without trusted launch; `trustedLaunchEnabled` would FORCE it on every consumer forever (all security flags are create-only and mutually exclusive)
- **The identifier triple is permanent** -- publisher/offer/SKU cannot change without replacing the definition and every version in it
- **Recommended sizing is advisory** -- shown to consumers, updatable in place

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `<your-compute-gallery>` | The Planton name of your `AzureComputeGallery` resource | Planton console (or replace `valueFrom` with `value:` and a literal gallery name) |
| `myorg` / `ubuntu` / `22-04-lts-gen2` | The image's permanent publisher/offer/SKU identity | Your naming convention |
| `eastus` | The Azure region of the definition | Your region strategy |

## Related Presets

- **Image With Version** -- the same definition publishing a snapshot-sourced release
