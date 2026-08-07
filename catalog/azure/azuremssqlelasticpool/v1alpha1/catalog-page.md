# Azure SQL Elastic Pool

Creates an Azure SQL elastic pool on an AzureMssqlServer logical server -- a shared-compute container that member databases (AzureMssqlDatabase with `skuName: ElasticPool`) draw from instead of carrying their own SKU. The right economics for many small databases with non-overlapping usage peaks, such as SaaS tenant-per-database fleets.

## What Gets Created

When you deploy an AzureMssqlElasticPool resource, Planton provisions:

- **Elastic Pool** -- an `azurerm_mssql_elasticpool` on the referenced server with your chosen SKU (the service tier and hardware family are derived from the SKU name, so a mismatched combination is unrepresentable), capacity, per-database bounds, storage cap, and availability posture

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An AzureMssqlServer** to create the pool on (referenced through `serverId`; the `region` must match the server's)

## Quick Start

Create a file `pool.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMssqlElasticPool
metadata:
  name: tenant-pool
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureMssqlElasticPool.tenant-pool
spec:
  serverId:
    valueFrom:
      kind: AzureMssqlServer
      name: my-sql
      fieldPath: status.outputs.server_id
  region: eastus
  poolName: tenant-pool
  skuName: StandardPool
  capacity: 100
  perDatabaseSettings:
    minCapacity: 0
    maxCapacity: 50
```

Deploy:

```shell
planton apply -f pool.yaml
```

Databases join the pool by setting `skuName: ElasticPool` and referencing `status.outputs.elastic_pool_id`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `serverId` | `StringValueOrRef` | The parent logical server. Defaults to referencing an `AzureMssqlServer`'s `server_id` output. Fixed at creation. | Required |
| `region` | `string` | Must match the parent server's region (ARM rejects a mismatch). | Required |
| `poolName` | `string` | Unique within the server. Changing it replaces the pool. | Required, ≤128 chars |
| `skuName` | `string` | `BasicPool`/`StandardPool`/`PremiumPool` (DTU) or `GP_Gen5`/`GP_Fsv2`/`GP_DC`/`BC_Gen5`/`BC_DC`/`HS_Gen5`/`HS_PRMS`/`HS_MOPRMS` (vCore). Tier + family are derived. | Closed vocabulary |
| `capacity` | `int32` | eDTUs (DTU pools) or vCores (vCore pools). ARM validates the exact ladder. | ≥1 |
| `perDatabaseSettings` | `object` | `minCapacity` (guaranteed per database; 0 oversubscribes) + `maxCapacity` (noisy-neighbor cap), in the pool's capacity unit. | Required, min ≤ max |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `maxSizeGb` | `double` | SKU default | The pool's total storage cap. Mutually exclusive with `maxSizeBytes`. |
| `maxSizeBytes` | `int64` | -- | Byte-precise storage cap. |
| `zoneRedundant` | `bool` | `false` | Spread replicas across availability zones (Premium/BC and zone-capable GP). |
| `enclaveType` | `enum` | -- | `VBS` / `DEFAULT_ENCLAVE` -- every database in the pool must share it. |
| `licenseType` | `enum` | LicenseIncluded | `BASE_PRICE` (Hybrid Benefit) / `LICENSE_INCLUDED`. vCore pools only. |
| `highAvailabilityReplicaCount` | `int32` | -- | Hyperscale pools only: readable HA replicas per database (0-4). |
| `maintenanceConfigurationName` | `string` | `SQL_Default` | Member databases inherit this window. |
| `tags` | `map(string)` | `{}` | User tags, merged over Planton-derived tags. |

## Examples

### vCore Pool with Hybrid Benefit

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMssqlElasticPool
metadata:
  name: gp-pool
spec:
  serverId:
    valueFrom:
      name: my-sql
  region: eastus
  poolName: gp-pool
  skuName: GP_Gen5
  capacity: 4
  licenseType: BASE_PRICE
  maxSizeGb: 512
  perDatabaseSettings:
    minCapacity: 0.25
    maxCapacity: 2
```

### Joining a Database to the Pool

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMssqlDatabase
metadata:
  name: tenant-42
spec:
  serverId:
    valueFrom:
      name: my-sql
  databaseName: tenant-42
  skuName: ElasticPool
  elasticPoolId:
    valueFrom:
      kind: AzureMssqlElasticPool
      name: tenant-pool
      fieldPath: status.outputs.elastic_pool_id
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `elastic_pool_id` | `string` | The pool's ARM ID -- referenced by `AzureMssqlDatabase.elasticPoolId` |
| `elastic_pool_name` | `string` | The pool's name |

## Related Components

- [AzureMssqlServer](/docs/catalog/azure/sql-server) — the parent logical server
- [AzureMssqlDatabase](/docs/catalog/azure/sql-database) — member databases joining via `elasticPoolId`
