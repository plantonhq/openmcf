# Azure MSSQL Elastic Pool

Deploys a shared-capacity elastic pool onto an existing Azure SQL logical server (AzureMssqlServer). The pool buys one capacity budget — eDTUs for DTU pools, vCores for vCore pools — and every member database draws from it: the many-small-databases answer, one bill for the whole fleet. Databases join the pool by referencing its `elastic_pool_id` output (which requires their `sku_name` to be the literal `ElasticPool` — the database wizard holds the pairing).

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **SQL Elastic Pool** -- on the referenced logical server, in the SAME region as the server (Azure rejects a mismatch), with the chosen SKU, shared capacity, and storage cap
- **Per-Database Limits** -- the guaranteed minimum and burst maximum every member database gets
- **Fleet-Wide Posture** -- zone redundancy, the confidential enclave, and the maintenance window every member inherits
- **Azure Tags** -- resource metadata tags applied to the pool for tracking and cost allocation

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An AzureMssqlServer** -- the logical server that hosts the pool, referenced through its `server_id` output. Only databases on this same server can join.

### Azure Subscription

- **Region planning** -- the pool's `region` must EQUAL the server's region exactly. There is no cross-region pool.
- **SKU planning** -- a closed 11-value catalog: DTU pools (`BasicPool`, `StandardPool`, `PremiumPool` — capacity in eDTUs) or vCore pools (`GP_Gen5`, `GP_Fsv2`, `GP_DC`, `BC_Gen5`, `BC_DC`, `HS_Gen5`, `HS_PRMS`, `HS_MOPRMS` — capacity in vCores, Azure Hybrid Benefit licensing available).

## Deploy

### Console

Open the deployment store, find **Azure MSSQL Elastic Pool**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **general-purpose-vcore** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMssqlElasticPool
metadata:
  name: tenants-pool
  org: acme-corp
  env: prod
spec:
  serverId:
    valueFrom:
      kind: AzureMssqlServer
      name: app-sql
      fieldPath: status.outputs.server_id
  region: eastus
  poolName: tenants-pool
  skuName: GP_Gen5
  capacity: 8
  perDatabaseSettings:
    minCapacity: 0
    maxCapacity: 2
```

```shell
planton apply -f mssql-elastic-pool.yaml
```

This creates an 8-vCore General Purpose pool where every member database can burst to 2 vCores with no reserved floor. A Stack Job tracks the provisioning in real time.

### InfraChart

The pool's `serverId` reference orders it after its server; pooled databases reference the pool's `elastic_pool_id` output and order after the pool — the InfraPipeline resolves the whole chain.

## Key Configuration

These are the most important decisions when configuring an elastic pool. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Capacity economics** -- size `capacity` for the fleet's combined PEAK, not the sum of individual peaks — that gap is the pool's entire cost advantage. Twenty databases that each burst to 2 vCores but average 0.2 often fit an 8-vCore pool.

**Per-database fairness** -- `perDatabaseSettings.maxCapacity` stops one noisy member from starving the fleet; `minCapacity` (0 = no reservation, Azure's default) guarantees a floor — non-zero minimums multiply by member count and must fit inside the pool.

**Storage cap** -- ONE value in ONE unit: `maxSizeGb` or `maxSizeBytes`, never both (they are mutually exclusive). Blank keeps the SKU's included storage.

**Fleet-wide posture** -- `zoneRedundant` (Premium/Business Critical only), `enclaveType` (Always Encrypted with secure enclaves for every member), `highAvailabilityReplicaCount` (Hyperscale pools only, 0-4), and `maintenanceConfigurationName` — the window every pooled database inherits.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureMssqlServer** | `serverId` | `status.outputs.server_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `elastic_pool_id` | Azure resource ID of the pool | AzureMssqlDatabase `elastic_pool_id` — how a database joins the pool |
| `elastic_pool_name` | Name of the pool | Monitoring, dashboards |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard DTU** -- The classic eDTU pool for mixed small workloads. Start from the **standard-dtu** preset.

**General Purpose vCore** -- The everyday production pool with Hybrid Benefit licensing available. Start from the **general-purpose-vcore** preset.

**Business Critical zone-redundant** -- Local SSD, built-in replicas, and zone redundancy for the fleet. Start from the **business-critical-zr** preset.

## Works With

- [**Azure MSSQL Server**](/cloud-catalog/azure-mssql-server) -- the logical server that hosts the pool (same region, always)
- [**Azure MSSQL Database**](/cloud-catalog/azure-mssql-database) -- members join via the pool's `elastic_pool_id` output
