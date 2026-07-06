# GCP Bigtable Table — Deep Dive

## The problem this resource solves

A Bigtable instance without tables stores nothing, and the table layer is where the decisions that determine cost and performance actually live: which column families exist, how long their data is retained, where the keyspace pre-splits, whether changes stream out, and whether the table can be deleted at all. Leaving those decisions to application code scatters them across services; leaving them to the console makes them invisible to review. This resource makes the table a first-class, reviewable infrastructure node — many per instance, each with an independent lifecycle.

## The container/content split

`GcpBigtableInstance` owns capacity and topology (clusters, zones, nodes, storage type, CMEK). `GcpBigtableTable` owns structure and retention. Teams add and remove tables without touching the instance; a table's GC policy change never disturbs its neighbors. Rows and cells remain application territory — infrastructure declares families, applications write columns into them freely (columns are created on write, never declared).

## Why GC policies are folded in — and why they matter

The provider models garbage collection as a separate resource, exactly one per (table, column family). That cardinality is the tell: a GC policy has no independent life — it cannot be referenced by anything else, cannot outlive its family, and is meaningless without it. Splitting it into its own kind would manufacture a glue node. So the spec folds a `gcPolicy` into each column family, and the modules manage the provider's per-family policy objects — policy changes never touch the table object or its data.

Operationally, GC is the load-bearing decision: **Bigtable never deletes old cell versions on its own.** A family without a policy accumulates every version of every write forever — the most common source of surprise Bigtable bills. Three expressible shapes:

- `maxAge` — drop cells older than a duration (time-series retention).
- `maxVersions` — keep only the newest N versions per cell.
- `mode: UNION|INTERSECTION` + both — combine the conditions (UNION collects when either is met; the common "cap age and versions" policy).

For deeper nested trees, `gcRules` carries the API's raw JSON rule tree, mutually exclusive with the typed fields.

One sharp edge is taught in the spec: on replicated (multi-cluster) instances, a change that EXPANDS what is eligible for collection is rejected by Bigtable unless `ignoreWarnings` is set — a deliberate safety gate against surprise data loss, mirrored as a per-policy flag.

## Aggregate column families

A family with `type: intsum` (or `intmin`/`intmax`/`inthll`) declares server-side aggregate cells: writes merge atomically into the existing value at the server, eliminating read-modify-write races and application-side counter logic. Usage metering and leaderboards stop needing transactions. Anything beyond the shortcuts is expressible as a raw JSON Type.

## Pre-splits: a creation-time-only decision

`splitKeys` seeds tablet boundaries so initial bulk load distributes across servers instead of hammering one tablet. The provider marks it ForceNew — changing it REPLACES the table and its data. The spec documents this loudly; after creation, splitting is Bigtable's own dynamic behavior plus operational tooling, not IaC.

## Change streams and automated backups

`changeStreamRetention` (1-7 days) turns on a CDC feed consumable by Dataflow — the integration point for event-driven pipelines downstream of Bigtable. `automatedBackupPolicy` (retention + frequency) is Bigtable's built-in table-level backup; both are mutable, and disabling follows the API's set-to-zero contract, handled by the modules.

## Deletion protection: API-side, defaulting PROTECTED

The provider exposes Bigtable's own `deletion_protection` string (PROTECTED/UNPROTECTED) — an API-side guard that blocks deletion by ANY client, console included, not just IaC. GCP's default is UNPROTECTED; the spec defaults PROTECTED, applying the catalog's data-bearing-kind convention, and both engines send the value explicitly so destroy behavior never depends on the engine. Destroying a table is a two-step: set UNPROTECTED, apply, destroy.

## Structured row-key schemas

`rowKeySchema` (the API's Type JSON) declares how row keys decompose into typed fields — the enabler for Bigtable's SQL-facing features and richer change-stream semantics. The API does not support in-place schema updates: the documented path (clear, apply, set new schema, apply) is recorded in the field comment.

## Mutability profile

| Surface | Mutability |
|---|---|
| `tableName`, `instance`, `splitKeys` | Immutable (ForceNew — splits REPLACE the table) |
| `columnFamilies` (add/remove) | Mutable |
| per-family `gcPolicy` | Mutable in place |
| `changeStreamRetention`, `automatedBackupPolicy` | Mutable |
| `deletionProtection` | Mutable (the unprotect-then-destroy path) |
| `rowKeySchema` | Clear-then-set only |

## IAM and API prerequisites

- `bigtableadmin.googleapis.com` enabled (both modules enable it with `disable_on_destroy=false`).
- `roles/bigtable.admin` (or the granular `bigtable.tables.*` permissions) on the project.

## Deliberately not modeled (recorded reasons)

- **`automated_backup_policy.locations`** — absent from the released google provider 6.x line (schema-probe verified); Enterprise-Plus zone pinning on newer lines.
- **App profiles** (`google_bigtable_app_profile`) — real multi-cluster routing control (single-cluster pinning, Data Boost isolation), an instance-scoped Tier-2 candidate on concrete pull.
- **Authorized/logical/materialized views, schema bundles** — newer view-layer resources; Tier-2 on concrete pull.
- **Table IAM trio** — resource-scoped IAM stays unmodeled catalog-wide.
- **`deletion_policy`** — client-side Terraform lever conflicting with Planton-managed destroy (catalog-wide decision).
