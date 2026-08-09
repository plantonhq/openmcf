# GcpBigtableTable Guide

The judgment this guide protects: garbage collection is opt-in and
load-bearing. Bigtable never deletes old cell versions on its own — a
column family without a GC policy accumulates every write forever,
which is the most common source of surprise Bigtable bills. Give every
family a policy at creation.

## Model families around retention, not entities

A column family is the unit of retention and access, not a relational
column. Group data by how long it should live and how it is read: a
`measurements` family with `maxAge: 720h` beside a `metadata` family
with `maxVersions: 3` is the canonical shape. Combining both conditions
requires `mode` (UNION drops when either is met — the common cap;
INTERSECTION only when both are). For rule trees deeper than one level,
`gcRules` takes the API's raw JSON — one shape or the other, never both.

## Splits are a day-zero decision

`splitKeys` pre-splits the key range so initial load distributes across
tablets instead of hammering one server — and it is immutable: changing
it REPLACES the table and its data. Set splits from your known key
distribution at creation, or leave them out and let Bigtable rebalance
organically; never plan to "fix splits later" through IaC.

## GC edits on replicated instances need an explicit override

Expanding what is eligible for collection on a multi-cluster instance
is rejected by Bigtable as a data-loss safety measure until
`ignoreWarnings: true` accompanies the change. Clusters can be briefly
inconsistent while the wider policy propagates — schedule loosening
deliberately.

## Backups: cadence, retention, and (on ENTERPRISE_PLUS) placement

`automatedBackupPolicy` turns on Bigtable's built-in backups with a
frequency and retention. `locations` restricts WHERE backups may be
created (zones in `projects/{project}/locations/{zone}` form; empty
means all zones of the instance) — accepted only for tables on
ENTERPRISE_PLUS instances, so the placement decision starts at the
parent instance's edition.

## Destroy semantics: an API guard plus a client policy

`deletionProtection` defaults PROTECTED — deletion by ANY client fails
until it is explicitly UNPROTECTED, the right default for a
data-bearing resource. `deletionPolicy` then decides what a permitted
destroy does, for BOTH objects this kind manages (the table and its
per-family GC policies): `PREVENT` as a second wall, `ABANDON` to drop
them from management — also the escape hatch when a GC-policy delete
is rejected on a replicated instance.

## What is deliberately absent

Rows and cells are application territory — this resource owns the
structure applications write into. App profiles, views, and schema
bundles are separate provider resources recorded as deferred
compositions; table-scoped IAM composes through project grants. The
table has no labels surface (organize at the instance).
