# Container-Scoped Reader

This preset is the least-privilege read grant: the built-in Data
Reader bound to a principal on exactly ONE container. An analytics
job, a dashboard, or a debugging operator sees the documents it needs
and nothing else -- not the neighboring containers, not the other
databases, no writes anywhere.

## When to Use

- Read-side consumers (analytics, search indexers, dashboards) that
  touch one container
- Operator or support access for a specific dataset, without key
  custody
- Demonstrating the narrowest grant shape before widening deliberately

## Key Configuration Choices

- **Data Reader by well-known ID** (`...0001`) -- point reads,
  queries, and the change feed; no write actions at all
- **The scope is a literal container path** -- composed on the
  account's ARM ID (`{account-id}/dbs/{db}/colls/{container}`);
  references cannot append path suffixes, so sub-account scopes are
  always literals
- **Inheritance stops at the container** -- unlike the account- or
  database-scoped shapes, nothing else is readable

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cosmosdb-account-resource-name>` | The AzureCosmosdbAccount's Planton resource name | Your Cosmos composition |
| `<subscription-id>` / `<resource-group-name>` / `<account-name>` | The coordinates of the account's ARM ID | The account's `cosmosdb_account_id` output |
| `<identity-resource-name>` | The AzureUserAssignedIdentity to grant | Your identity composition |
| `<database-name>` / `<container-name>` | The one container this principal may read | Your Cosmos composition |
