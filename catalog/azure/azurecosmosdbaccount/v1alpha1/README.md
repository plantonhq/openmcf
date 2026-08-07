# AzureCosmosdbAccount

An Azure Cosmos DB account -- the globally distributed, multi-model
database account that owns regions, consistency, network posture,
encryption, and backup for everything stored inside it. The data
containers are first-class kinds referencing the account:
AzureCosmosdbSqlDatabase / AzureCosmosdbSqlContainer for the SQL
(NoSQL) API and AzureCosmosdbMongoDatabase /
AzureCosmosdbMongoCollection for the MongoDB API.

## When to Use

Use AzureCosmosdbAccount when you need:

- **A globally distributed document database** -- single-digit-millisecond
  reads/writes with well-defined consistency, replicated to any set of
  Azure regions
- **The SQL (NoSQL) API** -- SQL-like queries over JSON documents (the
  default kind), including vector and full-text search capabilities
- **MongoDB compatibility** -- existing MongoDB drivers and tools against
  a fully managed backend (kind MONGO_DB + ENABLE_MONGO)
- **Serverless or provisioned throughput** -- pay-per-request via
  ENABLE_SERVERLESS, or provisioned RU/s on the databases and containers
  with an account-level `capacity` cap as the cost guardrail

## Key Configuration

- `account_name` -- globally unique DNS label (becomes
  `https://{name}.documents.azure.com`); fixed at creation
- `kind` -- GLOBAL_DOCUMENT_DB (SQL, default) or MONGO_DB; fixed at
  creation
- `consistency_policy` -- the account-wide read-consistency contract;
  SESSION is the recommended default, BOUNDED_STALENESS carries
  multi-region floors the spec enforces
- `geo_locations` -- the replicated regions; exactly one carries
  failover priority 0 (the write region);
  `multiple_write_locations_enabled` turns on active-active
- `capabilities` -- serverless, extra APIs (Cassandra/Gremlin/Table on a
  SQL account), MongoDB feature switches; most capability changes
  recreate the account
- `identity` + `default_identity` + `key_vault_key_id` -- the
  customer-managed-key story: a user-assigned identity that exists (and
  holds Key Vault grants) before the account does unwraps the CMK
- `backup` -- PERIODIC snapshots or CONTINUOUS point-in-time restore;
  `create_mode: RESTORE` + `restore` creates an account FROM a
  continuous-backup restore point
- `public_network_access_enabled`, `virtual_network_rules`,
  `ip_range_filter`, `local_authentication_enabled` -- the network and
  auth posture; disable local auth to force every data-plane caller
  through Entra ID

## Composition

```yaml
resourceGroup:
  valueFrom:
    kind: AzureResourceGroup
    name: data-platform
    fieldPath: status.outputs.resource_group_name
keyVaultKeyId:
  valueFrom:
    kind: AzureKeyVaultKey
    name: cosmos-cmk
    fieldPath: status.outputs.versionless_id
```

Databases and containers reference the account's
`cosmosdb_account_id` output; private endpoints use it as
`private_connection_resource_id` (subresource "Sql" or "MongoDB").

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
