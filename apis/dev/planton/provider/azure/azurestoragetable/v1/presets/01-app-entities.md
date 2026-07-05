# Application Entities Table

This preset creates a plain entities table -- the serverless key/value
store applications reach for when they need cheap, huge, schemaless
storage addressed by partition key + row key.

## When to Use

- User profiles, device state, feature flags, session data
- Any point-read workload that doesn't need secondary indexes or
  global distribution (Cosmos DB's Table API is the premium sibling)

## Key Configuration Choices

- **Partition design happens in the APPLICATION** -- the table itself
  has no schema; choose partition keys so hot workloads spread and
  entities needing atomic batch transactions share a partition
- **PascalCase name** -- the Table-storage convention (names forbid
  hyphens, unlike every other storage child)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The AzureStorageAccount's Planton resource name | Your storage composition |
| `<TableName>` | Letter-start 3-63 alphanumerics | Your naming convention |

## Downstream Wiring

Scope a data-plane grant to just this table:

```yaml
# On an AzureRoleAssignment
scope:
  valueFrom:
    kind: AzureStorageTable
    name: my-app-entities
    fieldPath: status.outputs.table_id
roleDefinitionName: Storage Table Data Contributor
```
