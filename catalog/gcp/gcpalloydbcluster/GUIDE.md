# GcpAlloydbCluster Guide

The judgment this guide protects: this one node is a database AND the
compute that serves it. Its destroy path has three independent levers,
its connectivity choice is permanent, and its restore surface exists only
at create time — decide those three things before the first apply.

## PSA or PSC — decided once, forever

Connectivity is exactly one of `network` (VPC peering via Private Service
Access) or `pscConfig.pscEnabled` (Private Service Connect), and both are
immutable. PSA needs the VPC prepared first (a `GcpGlobalAddress` with
VPC_PEERING purpose plus a `GcpServiceNetworkingConnection`); PSC needs
nothing on the VPC but moves connection control to
`primaryInstance.pscInstanceConfig` (consumer allowlists, auto-created
endpoints). Migrating between the two models is a new cluster plus a
restore — plan it like a migration, not an edit.

## Three destroy levers, in evaluation order

1. `deletionPolicy` is checked FIRST: `PREVENT` fails the destroy before
   anything else; `ABANDON` drops the cluster from management — and
   bypasses deletion_protection entirely.
2. `deletionProtection` (default TRUE) is the everyday guard: while it is
   true, destroy fails until you flip it false and APPLY that change
   first. Two deliberate steps between an engineer and data loss.
3. `FORCE` as the deletionPolicy also deletes any instances still in the
   cluster — the only way to destroy a SECONDARY cluster that has a
   secondary instance, and a footgun everywhere else.

The bundled primary has its own `primaryInstance.deletionPolicy`
(DELETE/PREVENT/ABANDON) because AlloyDB treats it as its own resource.

## The primary is bundled; everything else is composed

This kind always creates the cluster's PRIMARY instance — its sizing,
connectivity (public IP, PSC, PSA range override), pooling, and
stop/start lever all live under `primaryInstance`. Read pools and
secondary-cluster instances are separate `GcpAlloydbInstance` nodes; a
cluster-per-application with one right-sized primary plus read pools
added under load is the shape that scales.

## Stopping the primary without losing it

`primaryInstance.activationPolicy: NEVER` stops the compute (billing
stops with it) while storage and configuration survive; `ALWAYS` brings
it back. Stop read pools before the primary, start the primary before
the pools — AlloyDB enforces the ordering.

## Restores are new clusters

The four `restore*` sources (backup, PITR from a source cluster's
continuous backup, Backup-and-DR backup, Backup-and-DR PITR) are
mutually exclusive and create-time only. A point-in-time recovery is
therefore always a NEW cluster seeded from the old one's stream — verify
it, repoint applications, then retire the damaged cluster. Keep
`continuousBackupConfig` enabled (the default): a 14-day PITR window is
the cheapest insurance a database has.

## Backups have their own labels and their own keys

`automatedBackupPolicy.labels` is how backup storage shows up in cost
reports separately from the cluster; the policy's
`encryptionKmsKeyName` lets backups rotate keys independently of the
cluster's CMEK. Both matter to compliance reviews long before they
matter to engineers.
