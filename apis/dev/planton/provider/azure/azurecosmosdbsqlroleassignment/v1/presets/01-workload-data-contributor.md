# Workload Data Contributor

This preset is the standard keyless application grant: the built-in
Data Contributor role bound to a workload's managed identity across
the whole account. The app authenticates as its identity (no keys, no
connection-string secrets) and reads and writes every database in the
account -- the right shape when one workload owns the account.

## When to Use

- The account's single owning application, authenticating through a
  user-assigned managed identity
- Completing the keyless posture after disabling the account's key
  authentication
- Any workload that would otherwise hold the account's primary key

## Key Configuration Choices

- **Built-in by well-known ID** -- Data Contributor
  (`...0002`) exists in every account; its ID is a literal composed on
  the account's ARM ID, no AzureCosmosdbSqlRoleDefinition needed
- **The OBJECT ID, referenced** -- `principalId` rides the identity's
  `principal_id` output; passing a client ID deploys successfully but
  grants nothing (no directory object carries it)
- **Account-wide scope** -- right for a single-owner account; in
  shared accounts, narrow to `{account-id}/dbs/{db}` and prefer one
  grant per workload per database

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cosmosdb-account-resource-name>` | The AzureCosmosdbAccount's Planton resource name | Your Cosmos composition |
| `<subscription-id>` / `<resource-group-name>` / `<account-name>` | The coordinates of the account's ARM ID (for the built-in role's literal) | The account's `cosmosdb_account_id` output |
| `<identity-resource-name>` | The AzureUserAssignedIdentity the workload runs as | Your identity composition |
