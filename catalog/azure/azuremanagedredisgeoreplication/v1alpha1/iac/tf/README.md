# AzureManagedRedisGeoReplication - Terraform Module

Terraform implementation for the AzureManagedRedisGeoReplication
deployment component.

## Resources Created

- `azurerm_managed_redis_geo_replication.main` -- the group link
  joining the referenced Managed Redis instances into an active
  (multi-primary) geo-replication group

## Variable Highlights

| Variable | Notes |
| --- | --- |
| `spec.managed_redis_id` | The member the group is managed through; ForceNew |
| `spec.linked_managed_redis_ids` | The other members (1-4); adding links, removing force-unlinks |

## Behavior Notes

- ONE resource manages the whole group (linking is reciprocal); the
  provider serializes link/unlink operations and absorbs ARM's
  out-of-band updates to every member's replication state.
- Destroy force-unlinks all members -- each keeps its own copy of the
  data and becomes independent.
- Azure enforces the cross-member contracts at link time: same group
  name, `BALANCED_B3`+ SKUs, no persistence, geo-compatible modules.

## Provider Version

`azurerm ~> 4.0`. The `azurerm_managed_redis_geo_replication` resource
landed in recent 4.x releases -- a lockfile pinned to an old 4.x
resolves a provider without the resource; re-init to the current 4.x
line.

## Usage

```hcl
module "geo_group" {
  source = "./path/to/module"

  metadata = { name = "global-cache-group" }
  spec = {
    managed_redis_id         = "/subscriptions/.../redisEnterprise/global-cache-east"
    linked_managed_redis_ids = ["/subscriptions/.../redisEnterprise/global-cache-west"]
  }
}
```
