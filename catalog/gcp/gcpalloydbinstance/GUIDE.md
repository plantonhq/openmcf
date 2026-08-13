# GcpAlloydbInstance Guide

The judgment this guide protects: instances are the scaling and
connectivity dial of an AlloyDB cluster — cheap to add, cheap to stop,
and safe to change — EXCEPT for the three identity fields that recreate
them. Know which side of that line an edit falls on.

## READ_POOL is the default for a reason

The cluster kind already bundles the primary; this kind's everyday job
is read pools. Size with `readPoolConfig.nodeCount` and scale by editing
the count — node changes apply in place. The node count alone decides
availability: 1 node is zonal, 2+ nodes spread across zones for regional
read HA. Do not set `availabilityType` on a read pool — the API derives
it from the count and does not store a sent value (live-verified: the
stored object omits the field, so an explicit value would drift on every
plan; the spec refuses it up front). `availabilityType` belongs to
PRIMARY and SECONDARY instances, which exist for advanced topologies: a
SECONDARY instance serves a secondary (DR) cluster, and destroying one
is refused by the API — the secondary CLUSTER's `deletionPolicy: FORCE`
is the designed teardown path.

## The recreate line

`cluster`, `instanceId`, and `instanceType` are the identity — changing
any of them replaces the instance. So does `allocatedIpRangeOverride`
(the PSA range is chosen at IP-assignment time). Everything else —
machine size, flags, pooling, public IP, even `gceZone` — updates in
place; a zone change live-migrates.

## Stop/start is a spec edit

`activationPolicy: NEVER` stops the compute (billing stops; storage and
configuration survive); `ALWAYS` restarts it. Stop read pools before the
primary. This is the cheapest way to park a non-production cluster
overnight without losing its shape.

## Managed connection pooling earns its keep at high churn

`connectionPoolConfig.enabled: true` puts AlloyDB's built-in pooler in
front of the instance — the right call for serverless and web workloads
that open thousands of short-lived connections. Flags use the documented
convention: drop the "connection-pooling-" prefix and use underscores
("connection-pooling-pool-mode" becomes flag key `pool_mode`). Steady
long-lived connections (classic JVM pools) gain little; leave it off.

## Public IP is three decisions, not one

`enablePublicIp` opens the listener, `authorizedExternalNetworks` decides
who may reach it (required — the spec rejects networks without the IP),
and `sslMode: ENCRYPTED_ONLY` plus `requireConnectors` decide what a
connection must prove. Ship all of them together or none.

## Teardown discipline

`DELETE` (the default) removes the instance; the cluster and its data
survive. `PREVENT` suits the read pool a production service depends on —
losing serving capacity is an incident even when no data is lost.
`ABANDON` keeps the instance running (and billing) while dropping it
from management.
