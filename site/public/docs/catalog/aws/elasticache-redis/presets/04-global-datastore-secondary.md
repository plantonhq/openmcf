---
title: "Global Datastore Secondary (Cross-Region DR)"
description: "This preset joins an EXISTING ElastiCache global datastore as a read-only secondary in another region — the disaster-recovery and read-locality pattern. The secondary inherits its engine, engine..."
type: "preset"
rank: "04"
presetSlug: "04-global-datastore-secondary"
componentSlug: "elasticache-redis"
componentTitle: "ElastiCache Redis"
provider: "aws"
icon: "package"
order: 4
---

# Global Datastore Secondary (Cross-Region DR)

This preset joins an EXISTING ElastiCache global datastore as a read-only secondary in another region — the disaster-recovery and read-locality pattern. The secondary inherits its engine, engine version, node type, shard layout, encryption settings, and parameter group from the global primary, so this manifest deliberately sets none of them: it carries only the join id, the local topology (2 nodes with automatic failover across AZs), and the local networking.

## When to Use

- Cross-region disaster recovery with sub-second replication lag and promotable secondaries
- Read locality: serve users in a second region from a local read endpoint
- Compliance topologies that require a warm replica outside the primary region

## Key Configuration Choices

- **`globalReplicationGroupId`** — the join path to the global datastore. The global datastore itself is created alongside the PRIMARY replication group (outside this component today); this field attaches a secondary to it.
- **Inherited settings stay unset** — engine, engineVersion, nodeType, numNodeGroups, the encryption flags, and the parameter-group fields are inherited from the primary. Validation rejects them at authoring time, and AWS's provider rejects their presence at plan time — an explicit false is just as illegal as an explicit true.
- **`numCacheClusters: 2` + failover + multi-AZ** — the secondary's LOCAL node count and resilience; the keyspace layout still comes from the primary.
- **Local networking only** — subnets and security groups are per-region resources and always belong to the secondary itself.

## Placeholders to Replace

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `ldgnf-payments-global` | The global datastore id to join (rename to yours) | `ldgnf-payments-global` |
| `<private-subnet-id-az1>` | Private subnet AZ1 in the secondary region | `subnet-0a1b2c3d4e5f6g7h8` |
| `<private-subnet-id-az2>` | Private subnet AZ2 in the secondary region | `subnet-1a2b3c4d5e6f7g8h9` |
| `<security-group-id>` | SG allowing port 6379 in the secondary region | `sg-0123456789abcdef0` |

## Common Additions

- Add `preferredCacheClusterAzs` to pin the secondary's nodes to specific AZs
- Add `snapshotRetentionLimit` + `snapshotWindow` for secondary-side backups
- Add `notificationTopicArn` for failover and maintenance alerts in the DR region

## Related Presets

- **02-redis-ha-cluster** — the non-clustered HA shape this secondary mirrors locally
- **03-redis-clustered-production** — the primary-side production topology a global datastore is typically built from
