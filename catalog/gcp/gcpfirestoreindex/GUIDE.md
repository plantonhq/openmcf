# GcpFirestoreIndex Guide

The judgment this guide protects: every property of an index is
replace-on-change, and the replacement is not instant — Firestore builds the
new index in the background while the old one serves queries. Treat index
changes as rolling capacity operations, not config edits.

## When to use it (and when not)

Declare an index for every multi-field query shape your application ships —
Firestore rejects such queries outright until a matching index exists, and
the console's error-link workflow creates click-owned indexes nobody can
reproduce. One component per index: they are many-per-database with
independent lifecycles, and separate nodes let a query shape be added or
retired without touching its neighbors.

## The Enterprise surface needs the right database underneath

`searchConfig`, `apiScope: MONGODB_COMPATIBLE_API`, and `multikey` are
Firestore Enterprise capabilities. The database is where the edition lives:
pair with a `GcpFirestoreDatabase` whose `databaseEdition` is `ENTERPRISE`
(and `type: FIRESTORE_NATIVE` — the database spec enforces that pairing).
On a STANDARD-edition database these index shapes fail at APPLY time with a
provider error, not at validation — the spec cannot see the referenced
database's edition. `multikey` additionally requires the MongoDB-compatible
`apiScope` (validated pre-deploy); text search pairs with the same scope,
while geo search works under the default scope.

## Conventions and gotchas

- `skipWait` returns as soon as creation is REQUESTED. Use it when
  orchestrating many indexes and polling readiness yourself; do not use it
  when the next pipeline step immediately runs the query the index serves —
  the query fails until the background build completes.
- `deletionPolicy: ABANDON` unmanages an index while it keeps serving
  queries — the honest escape hatch when a collection has grown so large
  that a delete-and-rebuild cycle is an availability risk. `PREVENT` fails
  the destroy instead.
- The module always emits the `flat {}` vector marker — it is the only
  vector index layout GCP offers, so the spec deliberately has no knob for it.

## On the diagram

Each index renders as its own node hanging off the database — the set of
index nodes IS the reviewable inventory of the application's query shapes.
A click-created index exists in GCP but renders as nothing, which is exactly
the drift this component eliminates.

## Pairs well with

- `GcpFirestoreDatabase` — the `database` reference resolves its
  `database_name` output; Enterprise index shapes require the edition set
  there.
