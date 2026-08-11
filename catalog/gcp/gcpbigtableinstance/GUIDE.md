# GcpBigtableInstance Guide

The judgment this guide protects: Bigtable capacity decisions are
per-cluster and mostly reversible, but the cluster's PLACEMENT facts —
zone, storage type, CMEK key, scaling factor — are immutable, and on a
single-cluster instance changing any of them recreates the instance.
Design the topology first; tune the nodes forever after.

## Topology: clusters are replicas, zones are the failure domains

Every cluster is a full replica in its own zone (up to 8, across
regions). One cluster is a dev posture; production wants at least two
in different zones — the client library routes and fails over
transparently. Replication is automatic once a second cluster exists;
what you choose is WHERE the replicas sit relative to your readers.

## Three scaling modes; autoscaling is the production default

Fixed `numNodes` for predictable load; `autoscalingConfig`
(min/max/cpuTarget, optionally storageTarget) for everything else;
neither set lets Bigtable auto-allocate for the data footprint —
convenient for scratch instances, surprising for production. CPU
target accepts 10-80: lower means more headroom and more cost.
`NodeScalingFactor2X` scales in steps of two nodes for large workloads
— every bound must then be even, and the factor is immutable.

## SSD unless you can prove HDD

HDD costs less per TB and reads ~5% as fast for point lookups. It fits
large, scan-heavy, latency-tolerant archives — nothing else. The
storage type is per-cluster and immutable; a wrong HDD choice is a
data migration to fix.

## Edition and tags are org-level decisions

`edition` defaults to ENTERPRISE; ENTERPRISE_PLUS unlocks enterprise
features such as restricting table automated-backup placement
(`GcpBigtableTable.automatedBackupPolicy.locations`). The upgrade
applies in place; there is no downgrade. `resourceManagerTags` bind
org-policy/IAM tags at CREATE time only — changing them replaces the
instance, so settle tagging before production data lands.

## Destroy semantics are layered

`deletionProtection` (default TRUE) is the standing guard — destroy
fails until it is explicitly set false and applied. `forceDestroy`
clears backups that would otherwise block deletion. `deletionPolicy`
decides what a then-permitted destroy does: `PREVENT` as a second
independent wall, `ABANDON` to hand the running instance (and its
data) to another management plane.

## What is deliberately absent

Tables are `GcpBigtableTable` (many per instance, independent
lifecycles). App profiles, authorized/logical/materialized views, and
schema bundles are separate provider resources recorded as deferred
compositions. The DEVELOPMENT/PRODUCTION instance-type distinction is
deprecated in GCP — all instances are production; use a single small
cluster for dev.
