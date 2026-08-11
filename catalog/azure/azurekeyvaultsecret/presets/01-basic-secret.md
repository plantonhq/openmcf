# Basic Secret

This preset stores one secret in a Key Vault with the value wired by reference -- the manifest records that the secret exists and where it lives, never what it is. The example wires a storage account's connection string; swap the value reference for whatever emits yours.

## When to Use

- Database passwords, API keys, and connection strings that applications read from the vault at runtime
- The handoff point between a credential-emitting resource and its consumers
- Any secret that today lives in configuration files or pipeline variables

## Key Configuration Choices

- **The value is a reference** -- point it at another resource's sensitive output (or a managed secret); a literal in the manifest defeats the vault's purpose
- **Consumers should use `versionless_id`** -- value updates then propagate automatically; pin `secret_id` only under a compliance mandate
- **`contentType` is documentation** -- record what the value is (and its encoding for multi-line values) so consumers handle it correctly

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-key-vault>` | The AzureKeyVault the secret lives in | The vault component's name |
| `<your-storage-account>` | The example value source (an AzureStorageAccount whose connection string is stored) -- swap the whole `value.valueFrom` block for your real source | The producing component's name |

Secrets are effectively free -- Key Vault bills fractions of a cent per 10,000 operations.
