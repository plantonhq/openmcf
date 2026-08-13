# Starter Capacity

This preset creates the smallest Fabric capacity (F2) -- the right starting point for almost every environment, because the SKU scales up and down in place while the meter runs per hour from minute one.

## When to Use

- Standing up Fabric for the first time (development or a first production workload)
- Any environment where usage evidence should drive sizing, not anticipation

## Key Configuration Choices

- **F2** -- the smallest, cheapest step; move up when the Capacity Metrics app shows throttling (F64 is the tier that unlocks Copilot and free-license Power BI sharing)
- **One administrator** -- the platform group that manages the capacity from the Fabric side; the spec requires at least one at all times
- **Region matters** -- Fabric is not available in every Azure region; check availability before choosing

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `<admin-upn-or-object-id>` | The capacity administrator -- an Entra user principal name (e.g. `admin@contoso.com`) or a service principal's object ID | Entra ID -> Users / App registrations |
| `myorgfabric` | The capacity's name (3-63 lowercase letters and numbers, starting with a letter) | Your naming convention |
| `eastus` | The Azure region | Your region strategy (against the Fabric availability list) |

## Related Presets

- None yet -- workspace management lives in Microsoft's dedicated `fabric` provider, outside this kind.
