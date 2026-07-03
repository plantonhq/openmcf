# GcpCloudSqlDatabase — Design & Research

## What this component is

`GcpCloudSqlDatabase` models one logical database (`google_sql_database`) inside a Cloud SQL instance. It exists as a first-class kind because databases pass the split test decisively: an instance hosts many of them, each application's database has its own lifecycle (previews and staging datasets come and go), and none of that should ever require touching — or risking — the instance node.

## Design notes

- **Instance by name** — the `instance` ref resolves `GcpCloudSql.status.outputs.instance_name`. The provider's `instance` argument takes the bare instance name (not a self link), which is exactly what the output carries.
- **Charset/collation are engine-specific** and deliberately NOT cross-validated against the engine here: the database kind cannot see its instance's engine pre-deploy (the ref may be unresolved), and the API validates the combination authoritatively at create time. The field comments teach the per-engine rules instead:
  - MySQL: any supported pair; `utf8mb4` + `utf8mb4_0900_ai_ci` is the modern default (the legacy 3-byte `utf8` corrupts astral-plane characters).
  - PostgreSQL: `UTF8` only at creation; collation is an OS locale (`en_US.UTF8`).
  - SQL Server: charset ignored; collation is a SQL Server collation name.
- **No API enablement** — the hosting instance cannot exist without `sqladmin.googleapis.com` enabled; enabling it again here would only add churn on every database create/destroy.
- **Immutability** — name and instance are ForceNew; charset/collation update in place (the API alters the database).

## Deliberately unmodeled (with reasons)

- **`deletion_policy` (ABANDON)** — a client-side lever that conflicts with managed destroy semantics (the platform owns teardown ordering).
- **Grants/privileges** — what a user may do inside a database is schema territory, applied by migrations or application tooling; the Cloud SQL Admin API does not model per-database grants.

## Composition map

- `instance` ← `GcpCloudSql.status.outputs.instance_name`.
- Pairs naturally with `GcpCloudSqlUser` — one database + one user per application is the standard shape on a shared instance.

## Deployment methods

| Method | When to use | Notes |
|--------|-------------|-------|
| **Planton manifest** | Production composition inside infra charts | Reference the instance by `valueFrom.kind: GcpCloudSql`; charset/collation mistakes fail fast at the API |
| **Terraform module** | Standalone TF workflows | Plain-string `instance` after CLI ref resolution; no API enablement (the instance already enabled `sqladmin`) |
| **Pulumi module** | Standalone Pulumi workflows | Same contract as Terraform; exports `database_name` + `self_link` |
| **gcloud / console** | One-off admin tasks | Fine for bootstrap; not composable — no FK graph, no chart wiring |

## Operational guidance

- Dropping the resource drops the database and its data; the instance's backups/PITR are the recovery path.
- Charset/collation mistakes surface as API errors at create time — the preset examples carry known-good pairs per engine.
