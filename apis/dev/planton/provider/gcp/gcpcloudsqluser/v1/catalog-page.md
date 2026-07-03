# GCP Cloud SQL User

Creates a database user on a Cloud SQL instance — classic username/password (`BUILT_IN`) or passwordless IAM authentication for users, service accounts, and groups. One user per application, each with its own rotatable credential, instead of a shared admin login.

## What Gets Created

When you deploy a GcpCloudSqlUser resource, Planton provisions:

- **Cloud SQL user** — a `google_sql_user` on the referenced instance

## Prerequisites

- **An existing Cloud SQL instance** — a [GcpCloudSql](/docs/catalog/gcp/gcpcloudsql) resource (or a literal instance name)
- **For IAM users on PostgreSQL** — the instance must set `databaseFlags: {"cloudsql.iam_authentication": "on"}`
- **GCP credentials** with Cloud SQL admin permissions on the project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudSqlUser
metadata:
  name: orders-app-user
spec:
  instance:
    valueFrom:
      kind: GcpCloudSql
      name: my-postgres
      fieldPath: status.outputs.instance_name
  userName: orders-app
  password: a-strong-generated-password
```

```shell
planton apply -f user.yaml
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `instance` | `StringValueOrRef` | — (required) | The hosting instance (ref → GcpCloudSql). Immutable. |
| `userName` | `string` | — (required) | Login name (BUILT_IN) or IAM principal email (IAM types). Immutable. |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the instance. |
| `password` | secret | — | BUILT_IN credential; updating rotates in place. Never for IAM types. |
| `type` | `string` | `BUILT_IN` | `CLOUD_IAM_USER`, `CLOUD_IAM_SERVICE_ACCOUNT`, `CLOUD_IAM_GROUP` for passwordless IAM auth. |
| `host` | `string` | — | MySQL only: `user@host` scoping. |
| `passwordPolicy` | object | — | Lockout/expiry hardening for BUILT_IN users. |

## Examples

### Passwordless IAM Service Account User

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudSqlUser
metadata:
  name: ci-runner-user
spec:
  instance:
    valueFrom:
      kind: GcpCloudSql
      name: my-postgres
      fieldPath: status.outputs.instance_name
  userName: ci-runner@my-project.iam.gserviceaccount.com
  type: CLOUD_IAM_SERVICE_ACCOUNT
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `user_name` | The user name as stored by Cloud SQL (IAM users on MySQL are stored truncated before the `@`) |

## Related Components

- [GcpCloudSql](/docs/catalog/gcp/gcpcloudsql) — the instance this user lives on
- [GcpCloudSqlDatabase](/docs/catalog/gcp/gcpcloudsqldatabase) — the application database this user works in
- [GcpServiceAccount](/docs/catalog/gcp/gcpserviceaccount) — the identity behind IAM service-account users
