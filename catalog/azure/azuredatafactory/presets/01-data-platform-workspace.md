# Data Platform Workspace

This preset creates the standard environment factory -- a system-assigned identity for data-store access and the managed virtual network enabled up front, ready for pipelines and private endpoints to onboard against it.

## When to Use

- One factory per environment, many pipelines inside it (the recommended pattern -- the factory is near-free at rest)
- Platform-owned workspace posture with team-owned pipelines

## Key Configuration Choices

- **System-assigned identity** -- linked services authenticate as the factory; grant the `identity_principal_id` output on your data stores
- **Managed virtual network ON** -- free until used, required for managed private endpoints, and the one setting whose removal REPLACES the factory; enabling it up front keeps every door open
- **No git binding yet** -- bind GitHub or Azure DevOps later (an in-place update) once the authoring repo exists

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `my-org-data-platform` | The factory's name (globally unique across Azure) | Your naming convention -- org-prefixed |
| `eastus` | The Azure region | Your region strategy |

## Related Presets

- **02 Private ETL Workspace** -- the locked-down variant: public access off, a managed private endpoint to the data lake
