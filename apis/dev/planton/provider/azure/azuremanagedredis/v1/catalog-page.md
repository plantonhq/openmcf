# Azure Managed Redis

Azure's current-generation Redis service, built on Redis Enterprise.
Managed Redis is the successor to the retiring classic Azure Cache for
Redis and the target for new deployments: keyless by default (Entra
tokens instead of access keys), with Redis modules, active multi-primary
geo-replication, and customer-managed-key encryption.

## What Gets Created

When you deploy an AzureManagedRedis resource, Planton provisions:

- **Managed Redis cluster** -- an `azurerm_managed_redis` instance:
  compute, load balancer, and TLS endpoint at the chosen SKU, plus its
  default database (the Redis process) with your authentication,
  clustering, eviction, module, and persistence configuration

## Prerequisites

- **Azure credentials** configured via environment variables or Planton
  provider config
- **An AzureResourceGroup** to deploy into
- For customer-managed keys: an **AzureKeyVaultKey** and an
  **AzureUserAssignedIdentity** with wrap/unwrap access to it

## Quick Start

Create a file `cache.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureManagedRedis
metadata:
  name: app-cache
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureManagedRedis.app-cache
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: app-rg
      fieldPath: status.outputs.resource_group_name
  clusterName: app-cache
  skuName: BALANCED_B1
  defaultDatabase:
    accessKeysAuthenticationEnabled: false
    evictionPolicy: ALL_KEYS_LRU
```

Deploy:

```shell
planton pulumi up --manifest cache.yaml
```

Grant each consuming identity data-plane access with an
AzureManagedRedisAccessPolicyAssignment -- with keys off (the default),
grants are how clients connect at all.

## Spec Highlights

- `sku_name` -- tier family and memory size in one value: BALANCED
  (general purpose), COMPUTE_OPTIMIZED (max ops/sec), MEMORY_OPTIMIZED
  (large datasets), FLASH_OPTIMIZED (RAM + NVMe tiering); many changes
  apply in place
- `high_availability_enabled` -- a replica and the zone-redundant SLA
  (default true)
- `default_database.modules` -- RediSearch, RedisJSON, RedisBloom,
  RedisTimeSeries (RediSearch requires EnterpriseCluster clustering and
  NoEviction)
- `default_database.geo_replication_group_name` -- join an active
  (multi-primary) geo-replication group; link members with
  AzureManagedRedisGeoReplication
- `customer_managed_key` -- BYO encryption key from Key Vault
- Persistence -- AOF (every write) or RDB (periodic snapshots),
  mutually exclusive, not available on geo-replicated databases

## Stack Outputs

| Output | Description |
| --- | --- |
| `managed_redis_id` | The cluster's ARM ID -- what geo-replication and grants reference |
| `managed_redis_name` | The instance name |
| `region` | The Azure region |
| `resource_group_name` | The resource group |
| `hostname` | `{name}.{region}.redis.azure.net` -- all a keyless client needs, with the port |
| `database_id` | The default database's ARM ID (the data-plane scope) |
| `port` | The database port (10000) |
| `primary_access_key` | SECRET; empty under the keyless default |
| `secondary_access_key` | SECRET; empty under the keyless default |
| `identity_principal_id` | The system-assigned identity's principal, when enabled |
