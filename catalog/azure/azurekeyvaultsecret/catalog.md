# Azure Key Vault Secret

Stores a secret (password, API key, connection string) in an Azure Key Vault as a versioned, access-controlled data-plane object. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Key Vault secret** -- the versioned secret with its content-type hint, activation/expiry attributes, and tags

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureKeyVault** -- the vault the secret lives in (referenced by `keyVaultId`).
- **The value's source** -- a managed secret or another resource's output the sensitive `value` field references; avoid literals in manifests.

### Azure Subscription

- **Data-plane permissions** -- the deploying credential needs secret permissions on the vault (the "Key Vault Administrator" or "Key Vault Secrets Officer" RBAC role, or secret permissions in a legacy access policy). Subscription ownership alone is NOT enough.
- **Secrets are effectively free** -- Key Vault bills fractions of a cent per 10,000 operations.
- **Soft delete holds names** -- a deleted secret's name stays reserved for the vault's retention window unless purged (the module purges on destroy when the vault allows it).

## Deploy

### Console

Open the deployment store, find **Azure Key Vault Secret**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Basic Secret** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f secret.yaml
```

## After Deploy

Consumers that should follow value updates reference the `versionless_id` output; compliance-pinned consumers reference `secret_id` (one frozen version). The value itself is never an output -- readers fetch it from the vault at runtime with their own permissions.
