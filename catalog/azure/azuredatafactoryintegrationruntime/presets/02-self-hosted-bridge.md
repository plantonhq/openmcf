# Self-Hosted Bridge

This preset creates a self-hosted integration runtime registration -- the bridge between Data Factory and machines Azure cannot reach directly (on-premises servers, private networks, other clouds). Creating it is free; the compute is yours.

## When to Use

- Copying data from on-premises databases or file shares
- Reaching sources inside networks with no Azure connectivity
- Any source or sink a firewall separates from Azure's own compute

## Key Configuration Choices

- **`selfHosted: {}`** -- an empty block is a complete registration; Azure allocates it and issues two authorization keys, surfaced as this component's SENSITIVE outputs (`primary_authorization_key` / `secondary_authorization_key`)
- **The agent comes after deploy** -- install the integration runtime agent on your machine (from Studio's Manage -> Integration runtimes page) and paste a key; the node shows Running once it joins. Two or more nodes with the same key form a highly available cluster
- **Sharing instead of reinstalling** -- when another factory already runs a bridge, add `rbacAuthorization` with that runtime's ARM ID (and grant this factory's identity Contributor on it) instead of installing a second agent fleet; a linked registration issues no keys of its own

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |

## Related Presets

- **Data Flow Compute** -- the managed engine for mapping data flows.
- **SSIS Runtime** -- pairs with this bridge as its on-premises proxy (`proxy` block).
