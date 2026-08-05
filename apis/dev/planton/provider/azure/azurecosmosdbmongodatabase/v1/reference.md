# AzureCosmosdbMongoDatabase

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureCosmosdbMongoDatabaseSpec** defines the configuration for
creating a MongoDB API database inside an Azure Cosmos DB account: the
namespace collections live in and the boundary for SHARED throughput.

Databases are many-per-account with independent lifecycles, which is
why they are a first-class kind referencing the account rather than a
list folded into the account's spec. The parent is fixed at creation:
a database cannot move between accounts.

**Throughput model**: a database may provision throughput (fixed RU/s
or autoscale) that all its collections SHARE, or provision nothing and
let each collection bring its own dedicated throughput. On serverless
accounts (ENABLE_SERVERLESS capability) neither field may be set.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureCosmosdbMongoDatabase
metadata:
  name: test-cosmos-mongo-database
spec:
  cosmosdbAccountId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DocumentDB/databaseAccounts/planton-hack-cosmos-mongo
  databaseName: app-data
  # Exercises the shared fixed-throughput seam (mutually exclusive with
  # autoscale -- the spec rejects both together).
  throughput: 400
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.cosmosdbAccountId` | `string \| valueFrom` | yes |  | AzureCosmosdbAccount (`status.outputs.cosmosdb_account_id`) |
| `spec.databaseName` | `string` | yes |  |  |
| `spec.throughput` | `int32` |  |  |  |
| `spec.autoscaleMaxThroughput` | `int32` |  |  |  |

## Field Details

### spec.cosmosdbAccountId

`string | valueFrom` · required

The Cosmos DB account the database lives in, by ARM ID. References
an AzureCosmosdbAccount's cosmosdb_account_id output so the account
and its databases compose in one manifest set. The account must be a
MONGO_DB-kind account with the ENABLE_MONGO capability. Fixed at
creation.

- references: AzureCosmosdbAccount (`status.outputs.cosmosdb_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureCosmosdbAccount, name: <that resource's name>, fieldPath: status.outputs.cosmosdb_account_id}} -- a bare string does not parse

### spec.databaseName

`string` · required

The database's name -- unique within the account, 1-255 characters
(Azure's only constraint on Cosmos entity names). Changing the name
replaces the database AND everything in it.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"255"}}

### spec.throughput

`int32` · optional (explicit presence)

Fixed provisioned throughput in RU/s, SHARED by every collection in
the database that does not bring its own. Minimum 400, in increments
of 100. Mutually exclusive with `autoscale_max_throughput`; leave
both unset for serverless accounts or when every collection
provisions its own throughput.

- rule: throughput must be set in increments of 100 RU/s
- rule: {"int32":{"gte":400}}

### spec.autoscaleMaxThroughput

`int32` · optional (explicit presence)

Autoscale ceiling in RU/s: Azure scales the database's shared
throughput between 10% of this value and this value. Minimum 1000,
in increments of 1000. Mutually exclusive with `throughput`.

- rule: autoscale_max_throughput must be set in increments of 1000 RU/s
- rule: {"int32":{"gte":1000}}

## Validation Rules

- `cosmosdb_mongo_db_throughput_xor`: throughput and autoscale_max_throughput are mutually exclusive

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureCosmosdbMongoDatabase, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.mongo_database_id` | `string` | The Azure Resource Manager ID of the database -- what AzureCosmosdbMongoCollection's mongo_database_id references, and the management-plane scope for database-level operations. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.DocumentDB/databaseAccounts/{account}/mongodbDatabases/{name} |
| `status.outputs.mongo_database_name` | `string` | The database's name -- what MongoDB drivers reference inside the account's connection. |
| `status.outputs.cosmosdb_account_name` | `string` | The name of the Cosmos DB account the database lives in, parsed from the resolved account ID -- saves consumers a second reference when they need the account/database pair. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.cosmosdbAccountId` | AzureCosmosdbAccount | `status.outputs.cosmosdb_account_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureCosmosdbMongoCollection | `spec.mongoDatabaseId` | `status.outputs.mongo_database_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
