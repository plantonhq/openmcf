# GcpFirestoreDatabase Guide

The judgment this guide protects: almost everything that matters about a
Firestore database is decided at creation and cannot be changed —
location, name, edition, CMEK. The mutable knobs are the minority. Treat
this manifest as a decision record, not a tuning surface.

## The immutable four, decided once

- **Location**: multi-region (`nam5`, `eur3`) buys availability at higher
  write latency and cost; a single region is cheaper and faster. There is
  no migration — only export/import into a new database.
- **Name**: `(default)` is what every client library connects to when no
  database ID is given; one per project. Named databases isolate domains
  or tenants inside one project — but Firestore reserves a deleted
  database's name for several minutes, so automation that
  destroys-and-recreates the same name back-to-back will collide with the
  ghost (unique-per-run names sidestep it entirely).
- **Edition**: `STANDARD` for classic Firestore. `ENTERPRISE` (Native type
  only) unlocks the data-access-mode switches below — choose it up front
  if MongoDB compatibility is anywhere on the roadmap.
- **CMEK**: the key must live in the location-matched KMS region (`us` for
  `nam5`, `europe` for `eur3`); Google-managed encryption otherwise.

## The ENTERPRISE mode switches: single-protocol databases by API decree

An ENTERPRISE database is single-protocol, and not by convention — the
control plane rejects a create that enables both data-access modes ("Only
one of Firestore or MongoDB Compatible Data Access Mode can be enabled"),
so every database is either a Firestore-protocol or a MongoDB-protocol
database, decided at create. A MongoDB-dedicated database is
`mongodbCompatibleDataAccessMode: DATA_ACCESS_MODE_ENABLED` with the
classic API explicitly disabled; pair it with
`MONGODB_COMPATIBLE_API`-scoped `GcpFirestoreIndex` indexes, or queries
fail at runtime. `realtimeUpdatesMode` gates live query snapshots and
rides the CLASSIC Firestore API: enabling it requires
`firestoreDataAccessMode: DATA_ACCESS_MODE_ENABLED` explicitly (unset does
not count — the API rejects that too), which means a MongoDB-dedicated
database can never serve realtime updates. If you need both MongoDB
drivers and realtime subscriptions, that is two databases.

## PITR is not a backup strategy

`pointInTimeRecoveryEnablement` extends the version window from 1 hour to
7 days — recovery from a bad write, inside the live database. It does not
survive database deletion and does not reach past a week. Real protection
is a `GcpFirestoreBackupSchedule` (up to 14 weeks, survives deletion).
Production posture is both: PITR for oops-queries, scheduled backups for
disasters.

## Three deletion levers, three different questions

- `deleteProtectionState` (GCP-side): "may ANY client delete this
  database?" — the guard against console accidents.
- `deletionPolicy: PREVENT` (IaC-side): "may THIS tool's destroy delete
  it?" — the guard against pipeline accidents.
- `deletionPolicy: ABANDON`: "stop managing it, keep it running" — the
  honest exit when a database outlives its IaC.

The default `DELETE` keeps destroys real (the raw provider would default
to ABANDON and quietly orphan every destroyed database). Production
databases deserve `DELETE_PROTECTION_ENABLED` + `PREVENT`; ephemeral
environments keep the default so teardown actually tears down.

## Conventions and gotchas

- `type` is technically mutable (Native ↔ Datastore Mode) but is a
  data-model migration, not a config change — plan it like one.
- Firestore has no labels; `resourceManagerTags` is the only
  org-governance surface, and it is create-time-only (mutating replaces
  the database).
- Concurrency mode is a per-type default most workloads should not touch;
  `OPTIMISTIC_WITH_ENTITY_GROUPS` exists solely for legacy Datastore
  semantics.

## On the diagram

The database is the hub of the Firestore family: `GcpFirestoreIndex` nodes
(query shapes) and `GcpFirestoreBackupSchedule` nodes (protection cadence)
hang off it via its `database_name` output, making the whole posture —
what is queryable, what is protected — reviewable at a glance.

## Pairs well with

- `GcpFirestoreIndex` — every multi-field query shape, composed
  many-per-database.
- `GcpFirestoreBackupSchedule` — daily + weekly protection cadence.
- `GcpKmsKey` — CMEK, location-matched.
- `GcpProject` — the owning project (`projectId`).
