# Endpoint Access Association

This preset binds a machine to a Data Collection Endpoint (DCE) so agent CONFIGURATION access flows through the endpoint -- the shape for machines in locked-down networks reaching Azure Monitor over private link.

## When to Use

- Machines whose outbound access to Azure Monitor's public configuration endpoints is blocked (private-link estates, forced tunneling)
- At most one per machine -- endpoint associations are singular by Azure's design

## Key Configuration Choices

- **No `name`** -- Azure mandates the fixed name `configurationAccessEndpoint` for endpoint associations; leave the field unset and both engines apply it (setting anything else is rejected)
- **The endpoint is a literal ARM id** -- the Data Collection Endpoint is not yet a catalog kind; paste its ARM id
- **This association carries no rule** -- rule bindings ride their own associations (preset 01); a machine typically carries this endpoint association PLUS one or more rule associations

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-vm>` | The Planton name of your `AzureVirtualMachine` resource | Planton console (or replace `valueFrom` with `value:` and a literal VM ARM ID) |
| `<your-data-collection-endpoint-arm-id>` | The Data Collection Endpoint's ARM resource ID | The DCE's Properties page in the Azure portal |
