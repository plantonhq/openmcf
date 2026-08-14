# SQL Database with Key Vault Secrets

This preset creates an Azure SQL Database connection whose ENTIRE connection string lives in Key Vault -- the manifest carries an address into the vault, never a credential.

## When to Use

- Azure SQL connections authenticated by SQL credentials (username/password inside the connection string)
- Whenever the rotation story matters: rotate the secret in the vault; Data Factory picks it up on the next run

## Key Configuration Choices

- **`keyVaultConnectionString`** -- exactly one of this or an inline `connectionString` (validation enforces the choice); the vault form wins whenever available
- **`linkedServiceName` by reference** -- points at the KEY VAULT variant of this same kind (the Key Vault Connection preset); deploy that first
- **Managed identity alternative** -- for Entra-authenticated SQL, drop the secret entirely: `useManagedIdentity: true` plus a connection string without credentials, and make the factory's identity a database user

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-key-vault-connection>` | The Planton name of your Key Vault `AzureDataFactoryLinkedService` | Planton console (or replace `valueFrom` with `value:` and that connection's name) |
| `secretName` | The vault secret holding the full connection string | Your Key Vault's secrets list |

## Related Presets

- **Key Vault Connection** -- deploy it first; this preset references it by name.
- **Blob Storage via Managed Identity** -- the storage equivalent of the no-secret posture.
