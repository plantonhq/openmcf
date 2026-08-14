# Key Vault Connection

This preset creates the factory's secrets backbone: the Key Vault connection every other linked service's `keyVaultPassword` / `keyVaultConnectionString` blocks resolve through. Deploy it FIRST -- it is the reason no other connection needs a pasted secret.

## When to Use

- Once per factory, before any connection that carries credentials
- Whenever a rotation story matters: secrets rotate in the vault; manifests never change

## Key Configuration Choices

- **No secret material at all** -- Data Factory reaches the vault as its managed identity; grant that identity get/list on the vault's secrets (the connection saves fine without the grant and fails at run time)
- **A stable `name`** (`secrets-vault`) -- every other connection references this name in its Key-Vault-sourced secret blocks; renaming it later means touching them all

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-key-vault>` | The Planton name of your `AzureKeyVault` resource | Planton console (or replace `valueFrom` with `value:` and the vault's ARM ID) |

## Related Presets

- **Blob Storage via Managed Identity** -- the no-secret storage connection.
- **SQL Database with Key Vault Secrets** -- a database connection whose credentials live in this vault.
