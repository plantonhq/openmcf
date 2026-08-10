# GcpAlloydbInstance Guide

The judgment this guide protects: instances are the scaling and
connectivity dial of an AlloyDB cluster — cheap to add, cheap to stop,
and safe to change — EXCEPT for the three identity fields that recreate
them. Know which side of that line an edit falls on.

## READ_POOL is the default for a reason

The cluster kind already bundles the primary; this kind's everyday job
is read pools. Size with `readPoolConfig.nodeCount` (a 1-node pool can
only be ZONAL; 2+ nodes unlock REGIONAL) and scale by editing the count —
node changes apply in place. PRIMARY and SECONDARY types exist for
advanced topologies: a SECONDARY instance serves a secondary (DR)
cluster, and destroying one is refused by the API — the secondary
CLUSTER's `deletionPolicy: FORCE` is the designed teardown path.

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
