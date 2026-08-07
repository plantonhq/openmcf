---
title: "Standard DTU Pool"
description: "This preset creates a 100-eDTU Standard elastic pool -- the classic shared-compute container for many small databases whose usage peaks do not overlap (the SaaS tenant-per-database pattern)."
type: "preset"
rank: "01"
presetSlug: "01-standard-dtu"
componentSlug: "mssql-elastic-pool"
componentTitle: "MSSQL Elastic Pool"
provider: "azure"
icon: "package"
order: 1
---

# Standard DTU Pool

This preset creates a 100-eDTU Standard elastic pool -- the classic
shared-compute container for many small databases whose usage peaks do
not overlap (the SaaS tenant-per-database pattern).

## When to Use

- Fleets of small databases that would each waste most of a dedicated
  SKU
- Workloads already sized in DTUs migrating from standalone Standard
  databases

## Key Configuration Choices

- **`StandardPool` at 100 eDTUs** -- databases join by setting
  `skuName: ElasticPool` + `elasticPoolId` on their AzureMssqlDatabase
- **`maxCapacity: 50`** -- no single database can consume more than half
  the pool, capping noisy neighbors
- **`minCapacity: 0`** -- nothing is reserved per database, letting the
  pool oversubscribe (the usual economics)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<server-resource-name>` | The AzureMssqlServer's Planton resource name | Your server composition |
| `<server-region>` | The parent server's region (must match) | The server's spec |
| `<pool-name>` | The pool's name on the server | Your convention |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Join databases to the pool through its output:

```yaml
# On an AzureMssqlDatabase
skuName: ElasticPool
elasticPoolId:
  valueFrom:
    kind: AzureMssqlElasticPool
    name: <this-pool-resource-name>
    fieldPath: status.outputs.elastic_pool_id
```
