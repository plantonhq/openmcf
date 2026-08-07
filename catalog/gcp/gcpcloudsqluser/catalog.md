# GCP Cloud SQL User

Creates a database user on an existing Google Cloud SQL instance — a classic password (BUILT_IN) user, or a passwordless IAM-authenticated user, service account, or group. Users as first-class resources give every application its own principal on shared infrastructure: created, rotated, and revoked independently, never sharing the admin credential.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cloud SQL User** -- a `google_sql_user` on the referenced instance: BUILT_IN (username + password, with an optional per-user password policy) or one of the three IAM types (CLOUD_IAM_USER, CLOUD_IAM_SERVICE_ACCOUNT, CLOUD_IAM_GROUP — passwordless, authenticated through IAM)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project.
- **Org Secret** (BUILT_IN users) -- store the password as an org secret and reference it as `$secret/<slug>`; the runner resolves it just-in-time and the platform rejects plaintext.

### GCP Project

- **A Cloud SQL instance** -- the [GcpCloudSql](/cloud-catalog/gcp-cloud-sql) instance the user authenticates against. Reference it via ValueFromRef so the pipeline deploys the instance first.
- **IAM authentication flag** (IAM-typed users on PostgreSQL) -- the instance must carry the database flag `cloudsql.iam_authentication: "on"` before IAM users can be created.

## Deploy

### Console

Open the deployment store, find **GCP Cloud SQL User**, and click **Create**. The wizard walks two decisions: which instance the user lives on, then who the user is and how it authenticates. The [Presets](#presets) tab offers **Application User** (BUILT_IN) and **IAM Service Account User** (passwordless) starting points.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudSqlUser
metadata:
  name: orders-app
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  instance:
    valueFrom:
      kind: GcpCloudSql
      name: orders-db-prod
      fieldPath: status.outputs.instance_name
  userName: orders-app
  password: $secret/orders-app-db-password
```

```shell
planton apply -f user.yaml
```

## Key Configuration

**Authentication type** -- Immutable. `BUILT_IN` (the default) is a classic password user. The IAM types authenticate through Google IAM instead: nothing to store, leak, or rotate, and access is revoked by revoking the IAM role — the strongest posture for workloads that already run as a service account (GKE Workload Identity, Cloud Run).

**User name** -- Immutable, max 128 characters. For IAM types it is the principal's full email; on MySQL, GCP stores IAM user names truncated before the `@` (the `user_name` output reflects the stored name).

**Password** -- BUILT_IN only, and MUTABLE: updating the value rotates the credential in place, making this resource the rotation lever. Always a `$secret/<slug>` reference. Subject to the instance's password validation policy.

**Host scope** -- MySQL only: classic `user@host` scoping (e.g. `10.0.0.%`). Immutable; leave empty for PostgreSQL/SQL Server.

**Per-user password policy** -- BUILT_IN only: failed-attempts lockout, password expiration, and (MySQL) current-password verification — layered on the instance-wide policy.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpCloudSql** | `instance` | `status.outputs.instance_name` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `user_name` | The stored user name (IAM names may be truncated on MySQL) | Application connection strings |
| `instance_name` | The hosting instance's name | Cross-checking attachment |

## Works With

- [**GCP Cloud SQL**](/cloud-catalog/gcp-cloud-sql) -- the instance this user authenticates against
- [**GCP Cloud SQL Database**](/cloud-catalog/gcp-cloud-sql-database) -- the database the application connects to with these credentials
- [**GCP Service Account**](/cloud-catalog/gcp-service-account) -- the workload identity behind a CLOUD_IAM_SERVICE_ACCOUNT user
