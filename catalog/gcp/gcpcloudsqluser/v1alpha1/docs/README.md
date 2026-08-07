# GcpCloudSqlUser — Design & Research

## What this component is

`GcpCloudSqlUser` models one database user (`google_sql_user`) on a Cloud SQL instance. Users are first-class nodes because credentials have their own lifecycle — created per application, rotated on schedule, revoked on decommission — and none of that should require touching the instance.

## The two auth families

- **BUILT_IN** — classic username + password. The password is the one mutable field: updating it rotates the credential in place. Per-user hardening lives in `passwordPolicy` (lockout after failed attempts, expiration, MySQL current-password verification), layered on top of the instance-level password validation policy.
- **IAM types** (`CLOUD_IAM_USER`, `CLOUD_IAM_SERVICE_ACCOUNT`, `CLOUD_IAM_GROUP`) — passwordless database authentication through IAM: no credential to store, leak, or rotate. The user name IS the IAM principal email (MySQL stores it truncated before the `@` — the `user_name` output reflects the stored form clients authenticate with). On PostgreSQL, the instance must first set `cloudsql.iam_authentication = "on"`; principals also need `roles/cloudsql.instanceUser` and connect through the Auth Proxy / connectors.

The spec validates the family boundary pre-deploy: IAM types must not set a password; `passwordPolicy` is BUILT_IN-only.

## Design notes

- **Instance by name** — the `instance` ref resolves `GcpCloudSql.status.outputs.instance_name`.
- **Secret-by-default** — `password` is `(sensitive)`: encrypted in IaC state (Pulumi `ToSecret`; write-only to the API), never exported in outputs.
- **`host` is MySQL-only** (classic `user@host` semantics) and cannot be engine-validated here (cross-kind); the field comment teaches it.
- **No API enablement** — the hosting instance cannot exist without `sqladmin.googleapis.com` enabled.
- **Immutability** — name, instance, type, host are ForceNew; password and password policy update in place.

## Deliberately unmodeled (with reasons)

- **`password_wo` / `password_wo_version`** — write-only fields are an HCL plan-display mechanic; Planton's secret pipeline already keeps the credential out of rendered surfaces, so modeling both forms would create two ways to say one thing.
- **`database_roles`, `iam_email`** — not on the released 6.x provider line.
- **`deletion_policy` (ABANDON)** — client-side lever that conflicts with managed destroy semantics.
- **Grants** — what the user may do inside each database is schema territory (migrations), not Admin-API territory.

## Composition map

- `instance` ← `GcpCloudSql.status.outputs.instance_name`.
- `userName` for IAM SA users ← `GcpServiceAccount.status.outputs.email` — the keyless-workload story: WIF issues the workload's identity, this kind maps that identity into the database.

## Operational guidance

- One user per application; never share the instance admin user.
- PostgreSQL users that own objects cannot be dropped until ownership is reassigned — plan decommissions accordingly.
- Prefer IAM types wherever the workload already has a GCP identity; passwords are the fallback, not the default.
