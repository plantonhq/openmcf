# Private Gallery

This preset creates a private (RBAC-only) Compute Gallery -- the right starting point for an organization's golden-image library. The sharing posture is permanent, and private is almost always the correct choice.

## When to Use

- Standing up the organization's (or a team's) image library
- Any gallery whose consumers live in your own subscriptions (grant them RBAC read access)

## Key Configuration Choices

- **No sharing block** -- the gallery stays private; Groups/Community sharing is create-only, so choose deliberately before deploying
- **Package-style name** -- gallery names forbid dashes and allow dots; a name like `platform.images` reads well in every image ARM ID consumers see
- **Free at rest** -- the gallery bills nothing; image versions bill for their replica storage

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `platform.images` | The gallery's name (up to 80 letters, numbers, dots, underscores -- no dashes) | Your naming convention |
| `eastus` | The Azure region | Your region strategy |

## Related Presets

- **Community Gallery** -- publishes the gallery publicly under a generated community name
