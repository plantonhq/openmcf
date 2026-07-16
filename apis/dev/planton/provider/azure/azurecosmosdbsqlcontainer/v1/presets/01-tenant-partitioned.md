# Tenant-Partitioned Container

This preset creates the workhorse production container: a single
high-cardinality partition key, autoscale throughput that follows the
traffic curve, and an indexing policy that stops paying write RU for a
bulky payload subtree.

## When to Use

- Multi-tenant applications where /tenantId (or /userId, /deviceId)
  routes most queries
- Workloads with a daily traffic curve -- autoscale bills 10% of the
  ceiling when idle
- Documents carrying large blobs or nested payloads that queries never
  filter on

## Key Configuration Choices

- **`partitionKeyPaths: [/tenantId]`** -- high cardinality, even
  distribution, present in query filters; fixed at creation, so choose
  deliberately
- **`autoscaleMaxThroughput: 4000`** -- dedicated capacity scaling
  400-4000 RU/s; mutually exclusive with fixed `throughput`
- **`indexingPolicy`** -- includes `/*` explicitly (a declared policy
  replaces Azure's default wholesale) and excludes `/payload/*`; the
  policy updates in place, so tuning can come later

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-app-data` | The AzureCosmosdbSqlDatabase's Planton resource name | Your Cosmos composition |
| `orders` | The container name (1-255 chars, unique in the database) | Your data taxonomy |

## Downstream Wiring

Container-scoped data-plane grants target the container's ARM id:

```yaml
# Scope for a Cosmos DB data-plane role assignment
scope:
  valueFrom:
    kind: AzureCosmosdbSqlContainer
    name: my-orders
    fieldPath: status.outputs.sql_container_id
```
