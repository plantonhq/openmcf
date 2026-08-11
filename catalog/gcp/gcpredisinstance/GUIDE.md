# GcpRedisInstance Guide

The judgment this guide protects: this is the LEGACY Memorystore for
Redis — VPC-peered, AUTH-string, vertically sized. It remains the right
choice for existing estates and AUTH-based clients; for new deployments
weigh `GcpMemorystoreInstance` (Valkey, PSC, sharded) first, because
almost everything structural here is frozen at create time.

## The connectivity decision is the create-time decision

`connectMode`, `authorizedNetwork`, `reservedIpRange`, and
`transitEncryptionMode` are all immutable. `DIRECT_PEERING` is the
simple default; `PRIVATE_SERVICE_ACCESS` is REQUIRED for Shared VPC and
lets the instance consume a range you allocated — but it needs the
composition deployed first: `GcpGlobalAddress` (purpose VPC_PEERING) +
`GcpServiceNetworkingConnection` on the network, and `reservedIpRange`
then names the ALLOCATION, not a CIDR. Getting this wrong is a
replace-the-instance event, cache warm-up included.

## Size for the tier you will need, not the one you have

BASIC is a single node with no SLA — a restart loses everything and
serves nothing meanwhile. STANDARD_HA needs ≥5 GiB and is the floor for
anything production. Read replicas bolt onto STANDARD_HA only, and
enabling them on an EXISTING instance needs `secondaryIpRange` (the
original /29 has no room) — the one scale-out lever that works in
place.

## Maintenance windows coordinate systems, not just this instance

`maintenanceWindow` pins day + hour + minute (UTC); use `minute` to
stagger this cache's window against the database it fronts — a cache
failover DURING a database maintenance window is the compounding
failure. Write `description` for the humans who will wonder why the
window is Tuesday 03:30 ("after the nightly batch"). `maintenanceVersion`
is the self-service patch lever: set a newer available version to apply
now instead of waiting for GCP's rollout.

## Persistence is disaster insurance, not durability

RDB snapshots (STANDARD_HA) bound the data loss of a full failover to
the snapshot period — they do not make Redis a database. Anchor
`rdbSnapshotStartTime` into a low-traffic window; snapshot I/O competes
with serving.

## The two destroy guards do different jobs

`deletionProtection` defaults TRUE (both engines send it explicitly):
destroying the instance fails until you set it false and apply — the
two-step that turns a fat-fingered teardown into a visible plan diff.
`deletionPolicy` acts after that gate: `PREVENT` for a second
independent lock on caches whose loss would stampede the backing store,
`ABANDON` to hand the running instance to another management plane.
Remember what a cache destroy really costs: not the data (it is a
cache) but the cold-start load on whatever it was protecting.

## What is deliberately absent

Sharded, cluster-protocol Redis is a different resource
(`google_redis_cluster`) and a different product generation — that need
is served by `GcpMemorystoreInstance`, not by widening this kind.
