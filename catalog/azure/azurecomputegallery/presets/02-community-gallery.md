# Community Gallery

This preset creates a Community-shared Compute Gallery -- a public storefront that publishes every image version in it under a generated public name built from your prefix.

## When to Use

- Publishing images to the public Azure community catalog under your publisher identity
- Open-source projects and vendors distributing prebuilt VM images

## Key Configuration Choices

- **Community sharing is permanent and public** -- the whole sharing block is create-only, and everything published in the gallery becomes publicly deployable; keep community galleries separate from internal ones
- **The prefix becomes your public brand** -- Azure generates the public community name from it (read it back from the `community_gallery_name` output); it cannot be rebranded in place
- **EULA, publisher email, and URI are shown to every consumer** before they deploy from the gallery

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `<your-eula-url>` | The end-user license agreement URL shown to consumers | Your legal/docs site |
| `myorgimages` | The public name prefix (5-16 letters and numbers) | Your branding |
| `<your-publisher-email>` | The contact email shown to consumers | Your publishing team |
| `<your-publisher-uri>` | The publisher homepage shown to consumers | Your website |

## Related Presets

- **Private Gallery** -- the RBAC-only default posture for internal image libraries
