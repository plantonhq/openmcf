# CloudEvents Topic

This preset creates a custom topic accepting CloudEvents 1.0 -- the vendor-neutral event envelope and the right default for new integrations.

## When to Use

- New application eventing where you control the publishers
- Integrations that span clouds or vendors (every modern handler and SDK speaks CloudEvents)

## Key Configuration Choices

- **CloudEvents input schema** -- create-only; changing it later replaces the topic and its endpoint hostname, so it is chosen deliberately here
- **Region-wide unique name** -- the name becomes `{name}.{region}.eventgrid.azure.net`; replace the `my-org-` example prefix with your organization's to keep the claim collision-free and self-describing
- **Defaults left on** -- public network access and key auth stay enabled (Azure's defaults) for the fast integration path; see the locked-down preset for the hardened shape

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `my-org-app-events` | The topic's region-wide-unique name -- swap in your org prefix | Your naming convention |
