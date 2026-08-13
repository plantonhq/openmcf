# Blob Storage via Managed Identity

This preset creates the safest blob storage connection: the storage account's blob endpoint plus the factory's own managed identity -- no connection string, no SAS token, no key to rotate or leak.

## When to Use

- Any blob storage the factory's identity can be granted a role on (the default for lakehouse plumbing)
- Prefer this over `connectionString` / `sasUri` whenever possible -- those carry secret material this form simply does not have

## Key Configuration Choices

- **`serviceEndpoint` by reference** -- wires the AzureStorageAccount's `primary_blob_endpoint` output, so the connection follows the account
- **`useManagedIdentity: true`** -- grant the factory's identity **Storage Blob Data Contributor** (write) or **Storage Blob Data Reader** (read-only) on the account; the connection saves without the grant and fails at run time
- **Exactly one connection form** -- adding a connection string alongside the endpoint is rejected by validation (Azure stores only one)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-storage-account>` | The Planton name of your `AzureStorageAccount` resource | Planton console (or replace `valueFrom` with `value:` and the blob endpoint URL) |

## Related Presets

- **Key Vault Connection** -- the secrets backbone, for connections that DO need credentials.
- **SQL Database with Key Vault Secrets** -- the database equivalent of this posture.
