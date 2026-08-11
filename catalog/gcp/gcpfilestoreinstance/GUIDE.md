# GcpFilestoreInstance Guide

The judgment this guide protects: almost everything about a Filestore
instance is decided at create time — name, location, tier, protocol,
network attachment, CMEK key, LDAP, and replication are all
replace-the-instance-and-its-data events. The spec is wide, but the
mutable surface is narrow: capacity (up only), export options, labels,
protection flags, and the replica-state lever.

## Pick the tier for the failure mode you can live with

`BASIC_HDD`/`BASIC_SSD` (and their `STANDARD`/`PREMIUM` aliases) are
single-zone with no growth story beyond capacity — right for dev and
scratch. `ZONAL` is the modern single-zone tier and the only place
IOPS tuning (`performanceConfig`) starts making sense; `REGIONAL`
survives a zone loss; `ENTERPRISE` adds the replication story. The
legacy `HIGH_SCALE_SSD` exists for estates that already chose it —
prefer `ZONAL` for new high-performance shares. Capacity minimums
differ by tier (1 TiB most tiers, 2.5 TiB BASIC_SSD, 10 TiB
HIGH_SCALE_SSD) and capacity only ever grows.

## The network attachment is a composition, not a field

`DIRECT_PEERING` works out of the box. `PRIVATE_SERVICE_ACCESS` (the
Shared VPC requirement) rides an existing service-networking
connection: deploy `GcpGlobalAddress` (VPC_PEERING) +
`GcpServiceNetworkingConnection` on the VPC FIRST. PSC
(`PRIVATE_SERVICE_CONNECT`) adds `pscEndpointProject` — the consumer
project hosting the endpoint — and is the mode where per-export
`nfsExportOptions[].network` becomes REQUIRED (client IPs are not
otherwise attributable to a network).

## NFSv4.1 unlocks LDAP identity; NFSv3 is numeric matching

With NFSv3, user identity is whatever UID/GID the client presents —
fine on a trusted network, meaningless across teams. `protocol:
NFS_V4_1` (modern tiers only) plus `ldap` maps names through your
directory: `domain` + `servers` (all DNS names or all IPs, never
mixed), with `groupsOu`/`usersOu` as lookup-narrowing hints on large
directories. LDAP is immutable — retrofitting identity mapping onto a
share full of numerically-owned files is a migration, not a toggle.

## The DR story has three separate levers — know which one you need

`initialReplication` is CREATE-time: this instance joins as the STANDBY
of an existing ACTIVE peer (backups cannot be taken from a standby).
`desiredReplicaState` is the RUNTIME lever on that relationship: PAUSED
freezes the standby at a point in time (e.g. before risky maintenance
on the active), READY resumes. The restore sources are the RECOVERY
levers: `sourceBackup` restores from Filestore's own backups,
`sourceBackupdrBackup` from a Backup and DR vault — one or the other,
at create time, into a share whose capacity covers the backup.

## Two destroy guards, deliberately layered

`deletionProtectionEnabled` is the server-side guard: a protected
instance cannot be destroyed until the flag is flipped false and
applied — a visible two-step plan diff. `deletionPolicy` acts after it,
client-side: `PREVENT` as a second independent lock for shares whose
data exists nowhere else, `ABANDON` to hand the running instance to
another management plane. A file server is often the ONLY home of its
data — default production presets to protection on.

## What is deliberately absent

Filestore backups and snapshots are their own lifecycle (their own
resources and schedules), not instance spec. The deprecated provider
`zone` argument is covered by `location`.
