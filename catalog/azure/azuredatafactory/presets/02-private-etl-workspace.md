# Private ETL Workspace

This preset creates the locked-down factory: public network access off, integration running inside the managed virtual network, and a managed private endpoint carrying traffic to the data lake privately.

## When to Use

- Regulated environments where data movement must never touch the public internet
- Factories whose data stores already sit behind private endpoints

## Key Configuration Choices

- **Public access OFF** -- the factory's endpoints reject internet traffic; the managed private endpoints are its only data path
- **One managed private endpoint per data store** -- entries are create-only (a change replaces that endpoint, siblings untouched); each connection must be APPROVED on the target side before pipelines run
- **System-assigned identity** -- pair the private path with identity-based auth on the stores; grant the `identity_principal_id` output

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console (or replace `valueFrom` with `value:` and a literal group name) |
| `<your-storage-account-arm-id>` | The data lake storage account's ARM resource ID | Portal -> the storage account -> Properties -> Resource ID (or reference an AzureStorageAccount's ID output with `valueFrom`) |
| `my-org-private-etl` | The factory's name (globally unique across Azure) | Your naming convention -- org-prefixed |
| `eastus` | The Azure region | Your region strategy |

## Related Presets

- **01 Data Platform Workspace** -- the standard open-network variant for environments without private-link requirements
