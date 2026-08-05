# AzureMssqlElasticPool

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureMssqlElasticPoolSpec** defines the configuration for creating an
Azure SQL elastic pool on an AzureMssqlServer logical server.

An elastic pool is a shared-compute container: databases that join it
(AzureMssqlDatabase with sku_name "ElasticPool" + elastic_pool_id)
draw from the pool's eDTUs or vCores instead of carrying their own
SKU -- the right economics for many small databases with
non-overlapping usage peaks (SaaS tenant-per-database patterns).

**SKU model** (`sku_name` + `capacity`):
- DTU pools: "BasicPool", "StandardPool", "PremiumPool" with capacity
  in eDTUs (e.g. BasicPool 50-1600, StandardPool 50-3000, PremiumPool
  125-4000)
- vCore pools: "GP_Gen5", "GP_Fsv2", "GP_DC" (General Purpose),
  "BC_Gen5", "BC_DC" (Business Critical), "HS_Gen5", "HS_PRMS",
  "HS_MOPRMS" (Hyperscale) with capacity in vCores
The service tier and hardware family ARM wants alongside the name are
pure functions of it, so both engines derive them -- a
name/tier/family mismatch is unrepresentable.

**Sizing contract**: per_database_settings bounds what any ONE database
may consume (min guaranteed, max allowed); max_size_gb XOR
max_size_bytes caps the pool's total storage.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMssqlElasticPool
metadata:
  name: test-mssql-pool
spec:
  serverId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Sql/servers/test-mssql-server
  region: eastus
  poolName: tenant-pool
  # A vCore sku exercises the derived tier + family (GeneralPurpose /
  # Gen5) and the license enum mapping.
  skuName: GP_Gen5
  capacity: 4
  perDatabaseSettings:
    minCapacity: 0.25
    maxCapacity: 2
  maxSizeGb: 100
  zoneRedundant: false
  licenseType: LICENSE_INCLUDED
  maintenanceConfigurationName: SQL_Default
  tags:
    team: data
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.serverId` | `string \| valueFrom` | yes |  | AzureMssqlServer (`status.outputs.server_id`) |
| `spec.region` | `string` | yes |  |  |
| `spec.poolName` | `string` | yes |  |  |
| `spec.skuName` | `string` | yes |  |  |
| `spec.capacity` | `int32` | yes |  |  |
| `spec.perDatabaseSettings` | `AzureMssqlElasticPoolPerDatabaseSettings` | yes |  |  |
| `spec.perDatabaseSettings.minCapacity` | `double` |  |  |  |
| `spec.perDatabaseSettings.maxCapacity` | `double` | yes |  |  |
| `spec.maxSizeGb` | `double` |  |  |  |
| `spec.maxSizeBytes` | `int64` |  |  |  |
| `spec.zoneRedundant` | `bool` |  |  |  |
| `spec.enclaveType` | `enum` |  |  |  |
| `spec.licenseType` | `enum` |  |  |  |
| `spec.highAvailabilityReplicaCount` | `int32` |  |  |  |
| `spec.maintenanceConfigurationName` | `string` |  | `SQL_Default` |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.serverId

`string | valueFrom` · required

The logical server the pool is created on, by ARM ID. References an
AzureMssqlServer's server_id output; the server's name and resource
group are derived from it. Fixed at creation.

- references: AzureMssqlServer (`status.outputs.server_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureMssqlServer, name: <that resource's name>, fieldPath: status.outputs.server_id}} -- a bare string does not parse

### spec.region

`string` · required

The Azure region the pool is created in -- MUST match the parent
server's region (ARM rejects a mismatch; a pool cannot live in a
different region than its server). Changing it replaces the pool.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.poolName

`string` · required

The pool's name, unique within the server. Changing the name
replaces the pool.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.skuName

`string` · required

The pool SKU: a DTU pool (BasicPool/StandardPool/PremiumPool) or a
vCore pool ({GP|BC|HS}_{Gen5|Fsv2|DC|PRMS|MOPRMS}). The service tier
and hardware family are derived from it on both engines.

- rule: sku_name must be BasicPool, StandardPool, PremiumPool, GP_Gen5, GP_Fsv2, GP_DC, BC_Gen5, BC_DC, HS_Gen5, HS_PRMS, or HS_MOPRMS
- rule: {"required":true}

### spec.capacity

`int32` · required

The pool's capacity: eDTUs for DTU pools (e.g. 50, 100, 1200) or
vCores for vCore pools (e.g. 2, 4, 40). ARM validates the exact
ladder per SKU and region.

- rule: {"required":true,"int32":{"gte":1}}

### spec.perDatabaseSettings

`AzureMssqlElasticPoolPerDatabaseSettings` · required

Per-database consumption bounds inside the pool.

- rule: {"required":true}
- rule: per_database_settings.min_capacity cannot exceed max_capacity

### spec.perDatabaseSettings.minCapacity

`double`

The capacity every database is GUARANTEED (reserved even while
idle). 0 reserves nothing -- the usual choice, letting the pool
oversubscribe.

- rule: {"double":{"gte":0}}

### spec.perDatabaseSettings.maxCapacity

`double` · required

The capacity any one database may consume at peak. Caps noisy
neighbors; must be at least min_capacity and at most the pool's
capacity.

- rule: {"required":true,"double":{"gt":0}}

### spec.maxSizeGb

`double` · optional (explicit presence)

The pool's total storage cap in gigabytes. Mutually exclusive with
max_size_bytes. Set it explicitly on DTU pools -- ARM's ladder pins
the cap to the capacity (e.g. BasicPool at 50 eDTUs is 4.8828125)
and provisioning rejects an unset cap on some engines; vCore pools
may leave it unset for the SKU default.

- rule: {"double":{"gt":0}}

### spec.maxSizeBytes

`int64` · optional (explicit presence)

The pool's total storage cap in bytes, for sizes finer than whole
gigabytes. Mutually exclusive with max_size_gb.

- rule: {"int64":{"gt":"0"}}

### spec.zoneRedundant

`bool`

Spread the pool's replicas across availability zones. Premium and
Business Critical pools (plus zone-capable General Purpose regions).

### spec.enclaveType

`enum`

The confidential-computing enclave for the pool's databases. Every
database in the pool must use the same enclave type, so it is set
at the pool level. Changing it is disruptive -- plan accordingly.

Allowed values (use exactly as shown):

- `azure_mssql_elastic_pool_enclave_type_unspecified` -- Not specified: no enclave configured.
- `VBS` -- Virtualization-based security enclaves (Always Encrypted with secure enclaves) for every database in the pool.
- `DEFAULT_ENCLAVE` -- ARM's explicit "Default" (no enclave) -- distinct from unspecified so an update can actively clear a previously set enclave.

### spec.licenseType

`enum`

Azure Hybrid Benefit for vCore pools: BASE_PRICE brings your own SQL
Server license; LICENSE_INCLUDED pays as you go. Unset lets Azure
default (LicenseIncluded).

Allowed values (use exactly as shown):

- `azure_mssql_elastic_pool_license_type_unspecified` -- Not specified: Azure defaults to LicenseIncluded.
- `BASE_PRICE` -- Bring your own SQL Server license with Software Assurance.
- `LICENSE_INCLUDED` -- Pay-as-you-go, license included in the hourly rate.

### spec.highAvailabilityReplicaCount

`int32` · optional (explicit presence)

Hyperscale pools only: how many readable HA replicas back each
database in the pool (0-4).

- rule: {"int32":{"lte":4,"gte":0}}

### spec.maintenanceConfigurationName

`string` · optional (explicit presence)

The maintenance window Azure patches the pool in (e.g.
"SQL_Default", "SQL_EastUS_DB_1"). Unspecified applies
"SQL_Default". Databases in the pool inherit this window.

- default: `SQL_Default`

### spec.tags

`map<string, string>`

Free-form tags applied to the pool, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag
with the same key wins.

## Validation Rules

- `mssql_pool_max_size_xor`: max_size_gb and max_size_bytes are mutually exclusive
- `mssql_pool_ha_replicas_require_hyperscale`: high_availability_replica_count requires a Hyperscale pool sku (HS_)
- `mssql_pool_license_requires_vcore`: license_type applies to vCore pools only (GP_/BC_/HS_)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMssqlElasticPool, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.elastic_pool_id` | `string` | The Azure Resource Manager ID of the elastic pool. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Sql/servers/{server}/elasticPools/{name} Referenced by AzureMssqlDatabase.elastic_pool_id. |
| `status.outputs.elastic_pool_name` | `string` | The name of the elastic pool. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.serverId` | AzureMssqlServer | `status.outputs.server_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureMssqlDatabase | `spec.elasticPoolId` | `status.outputs.elastic_pool_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
