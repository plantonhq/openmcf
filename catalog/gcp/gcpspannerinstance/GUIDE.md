# GcpSpannerInstance Guide

The judgment this guide protects: the instance configuration (`config`)
is the one decision you cannot walk back — it pins geographic topology
for every database that will ever live on the instance. Everything else
(capacity, edition, autoscaling, backup defaults) updates in place with
no downtime, so err small and grow.

## The config choice is forever; the capacity choice is not

`regional-*` gives the lowest write latency in one region;
multi-region configs (`nam6`, `nam-eur-asia1`, ...) buy the 99.999%
SLA (on ENTERPRISE_PLUS) at higher cost and higher write latency
(quorum spans regions). Changing config means a new instance and a
data migration. Capacity is the opposite: nodes ↔ processing units ↔
autoscaling all switch in place, online. Start with 100-500 processing
units for anything unproven — 1 node (1000 PUs) is already ~10k QPS of
reads.

## Autoscaling: bound it like you mean it

The autoscaler operates strictly inside `autoscalingLimits` — an
undersized max is a ceiling you will hit during your first traffic
spike, an oversized min is a bill you pay every idle hour. Google's
recommended CPU targets are 65% (regional) and 45% (multi-region —
failover headroom matters). The `totalCpuUtilizationPercent` target
adds background work (change streams, backups) to the signal — use it
when background load is a real share of the instance. On multi-region
instances, `asymmetricAutoscalingOptions` scale a read-heavy replica
region independently: per-replica bounds (nodes or processing units —
per-replica PU bounds must be multiples of 1000), per-replica CPU
targets, or disabling a CPU signal for one replica entirely.

## The edition ladder is one-way in practice

STANDARD covers single-region basics. ENTERPRISE unlocks asymmetric
autoscaling and incremental backups; ENTERPRISE_PLUS unlocks the
multi-region SLA. Upgrades apply in place; downgrades require first
disabling every higher-edition feature — treat a downgrade as a
project, not a field edit.

## FREE_INSTANCE is a development lane, not a small production tier

One per billing account, ~10 GB, no capacity/edition/automatic-backup
configuration (the spec rejects all three pre-deploy). The upgrade to
PROVISIONED works in place — the reverse does not exist, so the free
slot is spent once.

## Destroy semantics are layered

`forceDestroy` governs the backups: false (default) makes destroy fail
while any database on the instance holds one — the last restore point
never rides along with a stack teardown silently. `deletionPolicy`
governs the instance itself: `PREVENT` for the instance a whole
topology depends on, `ABANDON` to hand a running instance to another
management plane. The two compose: a DELETE policy still stops on
backups unless forceDestroy is armed.

## What is deliberately absent

Databases and backup schedules are their own composable kinds
(`GcpSpannerDatabase`, `GcpSpannerBackupSchedule`) — the instance is
capacity and topology only. Custom instance configs and instance
partitions are separate provider resources, recorded as deferred
compositions.
