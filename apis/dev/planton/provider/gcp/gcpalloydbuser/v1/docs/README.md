# GcpAlloydbUser — Design & Research

## What this component is

`GcpAlloydbUser` models one database user (`google_alloydb_user`) on an AlloyDB cluster — the AlloyDB analogue of `GcpCloudSqlUser`.

## The two auth families

- **ALLOYDB_BUILT_IN** — username + password; password is mutable (rotates in place).
- **ALLOYDB_IAM_USER** — passwordless IAM authentication; `user_id` is the IAM principal email.

CEL enforces: IAM users must not set a password.

## Design notes

- **Cluster by full path** — `cluster` ref resolves `GcpAlloydbCluster.status.outputs.cluster_id`.
- **Secret-by-default** — `password` is `(sensitive)`; never exported in outputs.
- **API enablement** — enables `alloydb.googleapis.com` so users deploy on fresh projects without relying on other AlloyDB kinds.

## Deliberately unmodeled

- **`password_wo` / `password_wo_version`** — write-only HCL mechanics; Planton secret pipeline covers this.
- **`deletion_policy` (ABANDON/PREVENT)** — conflicts with managed destroy semantics.

## Composition map

- `cluster` ← `GcpAlloydbCluster.status.outputs.cluster_id`
- IAM `userId` ← service account email from `GcpServiceAccount` for keyless workloads
