# AzureManagedRedis -- Design Research

## The Resource

Azure Managed Redis (AMR) is the `Microsoft.Cache/redisEnterprise` ARM
family at API version `2025-07-01`: a CLUSTER (compute, load balancer,
network, TLS endpoint) plus its DEFAULT DATABASE (the Redis process --
eviction, clustering, modules, persistence, authentication), mapped
1-to-1. The component maps onto `azurerm_managed_redis` (azurerm v4.x,
`internal/services/managedredis/managed_redis_resource.go`),
parity-verified against pulumi-azure v6 (`managedredis.ManagedRedis`).
Azure has deprecated `azurerm_redis_enterprise_cluster`/`_database` in
favor of this resource, and classic Azure Cache for Redis is being
retired -- AMR is the target for new Redis deployments.

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `cluster_name` | 3-63 letters/digits/hyphens, letter/digit ends, no consecutive hyphens (the provider's exact validator as CEL); ForceNew |
| `resource_group_name` | `resource_group` | `StringValueOrRef` → AzureResourceGroup |
| `location` | `region` | AMR is not yet in every region |
| `sku_name` | same | Closed 44-value enum (the provider's exact vocabulary; legacy `Enterprise_*`/`EnterpriseFlash_*` SKUs are excluded by the provider). ARM wire values (`Balanced_B0` style) mapped row-by-row in both modules |
| `high_availability_enabled` | same | `optional bool` default true (azurerm's default); ForceNew |
| `customer_managed_key` | same | `key_vault_key_id` → `AzureKeyVaultKey.key_id` (VERSIONED -- the provider validates a versioned nested-item ID; rotation = updating the reference); `user_assigned_identity_id` → the UAI ARM id |
| `identity` | same | SystemAssigned/UserAssigned with the ids-match-type CEL |
| `public_network_access` | `public_network_access_enabled` | The provider's Enabled/Disabled string modeled as the catalog's bool grain; both modules map bool → string |
| `default_database` (block) | same (required message) | The provider requires it at create ("`default_database` must be provided when creating a new resource"); a database-less cluster is an administrative transient, not a declarative target |
| `.access_keys_authentication_enabled` | same | Plain bool, default FALSE -- keyless-first (the reverse of classic Redis) |
| `.client_protocol` | same | Closed enum (Encrypted/Plaintext), unspecified deploys Encrypted; both modules send the default explicitly |
| `.clustering_policy` | same | Closed enum (EnterpriseCluster/OSSCluster/NoCluster), unspecified deploys OSSCluster; changing it recreates the DATABASE |
| `.eviction_policy` | same | Closed 8-value enum, unspecified deploys VolatileLRU |
| `.geo_replication_group_name` | same | The provider's exact name validator as CEL; group membership is linked by AzureManagedRedisGeoReplication; changing it recreates the DATABASE |
| `.module` | `modules` | ≤4; names are the provider's closed set (RediSearch/RedisJSON/RedisBloom/RedisTimeSeries) as an in-list CEL; uniqueness CEL added (ARM rejects duplicates) |
| `.persistence_append_only_file_backup_frequency` | same | Presence enables AOF; "1s" is the only value Azure offers (the deprecated `always` is filtered by the provider) |
| `.persistence_redis_database_backup_frequency` | same | Presence enables RDB; 1h/6h/12h |
| `tags` | same | User tags merged over the platform-derived tags |

## Validation Rules Mirrored from the Provider's Apply-Time Validators

The provider enforces these in CustomizeDiff -- all statically checkable
and mirrored as CELs so they fail at validation, not apply:

- `default_database` required (modeled as a required message -- see
  above).
- Geo-replication requires `BALANCED_B3`+ (the provider excludes B0/B1
  from `SKUsSupportingGeoReplication`).
- A geo-replicated database allows only the RediSearch/RedisJSON
  modules.
- RediSearch requires `NO_EVICTION` and `ENTERPRISE_CLUSTER` (checked
  against the EFFECTIVE clustering policy: unspecified means OSSCluster,
  which also fails -- matching the provider's behavior on defaults).
- AOF XOR RDB persistence; both conflict with geo-replication (the
  provider's `ConflictsWith` pairs).

Left to apply time (cross-resource, not statically checkable):

- **In-place SKU scaling** -- the provider calls the live
  ListSkusForScaling API to decide update-vs-replace; documented on the
  field.
- **The CMK identity must also be attached through the identity block**
  -- an ARM pairing across two spec fields whose values may be
  references (CEL cannot dereference `StringValueOrRef` sub-fields --
  the protovalidate-java constraint); documented on both fields.

## Fold Verdict

**The default database FOLDS into the cluster** -- Azure's own grain:
the mapping is 1-to-1 (the provider's source folds them deliberately
after users expected clusters to work by themselves), the database has
no life without its cluster, and nothing FK-references the database as
a separate node (grants and geo-links reference the CLUSTER and derive
the database path). If Azure ever ships multi-database support it will
arrive as a new child resource, and a standalone kind lands then.

## Recorded Skips (with reasons)

- **`minimum_tls_version`** -- the provider accepts exactly one value
  (`TLS12`) and hardcodes it in its create path; a one-value knob is a
  constant, not configuration. Neither module sends it. Lands as a
  field if Azure ever adds a second floor.
- **ARM `deferUpgrade` and writable `port` (database)** -- present in
  the ARM API but expressible by NEITHER engine (not in the azurerm
  schema, hence not in the bridge); a field neither engine can deploy is
  unshippable. Land when the provider models them.
- **ARM read-backs `redisVersion` and `redundancyMode`** -- not exposed
  as attributes by either engine, so they cannot be honest outputs.
- **The `flush databases` action** -- an imperative operation, not
  declarative state.
- **Module `version`** -- output-only in the provider; the module list
  is configuration, versions are Azure's.
- **Data source / classic `zones`** -- AMR manages zone placement
  through `high_availability_enabled`; there is no zones argument.

## Operational Behavior Worth Knowing

- **Provisioning runs tens of minutes** (the provider budgets 45 min
  create / 30 min delete): the cluster polls to Running, then the
  database is created.
- **Database-recreating changes**: `clustering_policy`,
  `geo_replication_group_name`, and the module set delete and recreate
  the DATABASE in place -- data is lost and the endpoint is briefly
  unavailable, but the cluster (and its hostname) survives.
- **The endpoint is `{name}.{region}.redis.azure.net` on port 10000**
  -- not classic Redis's 6379/6380.
- **Keyless connections**: username = the granted object ID, password =
  an Entra token (`az account get-access-token --scope
  https://redis.azure.com/.default` for humans; the identity SDK for
  workloads).
- **Network isolation is Private Link only** -- no VNet injection and
  no IP firewall (both were classic-Redis mechanisms); disable public
  network access and attach an AzurePrivateEndpoint.
