# Azure Cache for Redis

Creates a fully managed Redis instance -- caching, session state,
leaderboards, and pub/sub with sub-millisecond latency -- with the full
production surface: tier/size selection, Microsoft Entra (token)
authentication up to a fully keyless posture, VNet injection,
clustering, RDB/AOF persistence, patch windows, and IP firewall rules.

> **Retirement notice:** Azure is retiring classic Azure Cache for Redis
> in favor of Azure Managed Redis -- ARM has begun rejecting NEW cache
> creations region by region (observed first on Premium creations, while
> Basic/Standard elsewhere still succeed). Existing caches keep running
> and this kind manages them fully; choose **Azure Managed Redis** for
> new deployments.

## What Gets Created

When you deploy an AzureRedisCache resource, Planton provisions:

- **Redis Cache** -- an `azurerm_redis_cache` at the chosen tier and
  size, with engine configuration, identity, zones, and networking
- **Firewall Rules** -- one `azurerm_redis_firewall_rule` per entry in
  `firewallRules`, allow-listing IPv4 ranges on the public endpoint

## Prerequisites

- **Azure credentials** configured via environment variables or Planton
  provider config
- **An AzureResourceGroup** to create the cache in
- **For VNet injection (Premium)**: an AzureSubnet dedicated to Redis
- **For managed-identity persistence**: grant the cache's identity
  "Storage Blob Data Contributor" on the persistence storage account

## Quick Start

Create a file `cache.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRedisCache
metadata:
  name: app-cache
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureRedisCache.app-cache
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      kind: AzureResourceGroup
      name: app-rg
      fieldPath: status.outputs.resource_group_name
  cacheName: my-app-cache
  skuName: STANDARD
  capacity: 1
  redisConfiguration:
    activeDirectoryAuthenticationEnabled: true
    maxmemoryPolicy: allkeys-lru
```

Deploy:

```shell
planton apply -f cache.yaml
```

Provisioning takes 15-40 minutes -- Redis is the slowest-provisioning
service in the Azure catalog. Deletion also runs several minutes; the
globally unique name becomes reusable once the delete completes.

## Key Outputs

| Output | Purpose |
|--------|---------|
| `redis_cache_id` | The reference target for linked servers, access policies/grants, and private endpoints |
| `hostname` / `ssl_port` | The endpoint (`{name}.redis.cache.windows.net:6380`); all a keyless client needs |
| `region` | The linked-server location seam for geo-replication |
| `primary_access_key` / `secondary_access_key` | The shared keys (secret-bearing); both faces stay live for zero-downtime rotation |
| `primary_connection_string` / `secondary_connection_string` | Ready-to-use connection strings (secret-bearing) |
| `identity_principal_id` | The system-assigned identity -- RBAC grant target for managed-identity persistence |

## Related Resources

- [Azure Redis Linked Server](/docs/catalog/azure/azureredislinkedserver) -- geo-replication / DR pairing
- [Azure Redis Cache Access Policy](/docs/catalog/azure/azurerediscacheaccesspolicy) -- custom data-plane permission sets
- [Azure Redis Cache Access Policy Assignment](/docs/catalog/azure/azurerediscacheaccesspolicyassignment) -- Entra grants (the keyless story)
- [Azure Resource Group](/docs/catalog/azure/azureresourcegroup) -- the container
- [Azure Subnet](/docs/catalog/azure/azuresubnet) -- VNet injection (Premium)
