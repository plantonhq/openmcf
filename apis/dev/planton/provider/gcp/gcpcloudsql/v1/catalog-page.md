# GCP Cloud SQL

Deploys a fully managed MySQL, PostgreSQL, or SQL Server instance with the production surface real workloads need: regional high availability, point-in-time recovery, explicit connectivity (public IPv4, private VPC IP, or Private Service Connect), Query Insights, password policies, managed connection pooling, customer-managed encryption, dual delete guards, and read replicas. The component automatically enables the Cloud SQL Admin API on the target project.

## What Gets Created

When you deploy a GcpCloudSql resource, Planton provisions:

- **Cloud SQL Admin API enablement** — activates `sqladmin.googleapis.com` on the target project
- **Cloud SQL instance** — a `google_sql_database_instance`: a primary, or a read replica when `masterInstanceName` is set

Databases and users are separate composable resources ([GcpCloudSqlDatabase](/docs/catalog/gcp/gcpcloudsqldatabase), [GcpCloudSqlUser](/docs/catalog/gcp/gcpcloudsqluser)) that reference this instance by name.

## Prerequisites

- **GCP credentials** configured via environment variables or Planton provider config
- **An existing GCP project** — referenced via `projectId` (or the provider's default project)
- **For private IP** — private services access on the VPC: a [GcpGlobalAddress](/docs/catalog/gcp/gcpglobaladdress) (purpose `VPC_PEERING`) + a [GcpServiceNetworkingConnection](/docs/catalog/gcp/gcpservicenetworkingconnection)
- **IAM permissions** — `roles/cloudsql.admin` on the target project

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudSql
metadata:
  name: my-postgres
spec:
  projectId:
    value: my-gcp-project-123
  instanceName: orders-db
  region: us-central1
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_16
  tier: db-custom-2-7680
  backup:
    enabled: true
    pointInTimeRecoveryEnabled: true
```

```shell
planton apply -f postgres.yaml
```

This creates a PostgreSQL 16 instance whose public IP has **no authorized networks** — reachable only through the IAM-authenticated Cloud SQL Auth Proxy or connectors.

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `instanceName` | `string` | — (required) | Cloud-side name (RFC1035). Immutable; reserved ~1 week post-delete. |
| `region` | `string` | — (required) | GCP region. Immutable. |
| `databaseEngine` | `string` | — (required) | `MYSQL`, `POSTGRESQL`, or `SQLSERVER`. |
| `databaseVersion` | `string` | — (required) | e.g. `POSTGRES_16`. Mutable (in-place upgrade). |
| `tier` | `string` | — (required) | Machine type. Mutable (in-place resize). |
| `projectId` | `StringValueOrRef` | provider default | Project that owns the instance. |
| `edition` | `string` | `ENTERPRISE` | `ENTERPRISE_PLUS` adds the data cache + 99.99% HA SLA. |
| `availabilityType` | `string` | `ZONAL` | `REGIONAL` = HA with automatic failover (requires backups). |
| `disk` | object | 10 GB PD_SSD | Type/size/auto-resize. Disks grow, never shrink. |
| `network` | object | public, no allowlist | Private VPC IP (ref → GcpVpc), public IPv4 + authorized networks, TLS posture, server CA mode, PSC. |
| `backup` | object | disabled | Daily backups, PITR (PG/SQL Server) or binary logs (MySQL), retention. |
| `maintenanceWindow` / `denyMaintenancePeriod` | object | — | Weekly window + freeze periods. |
| `insightsConfig` | object | disabled | Query Insights telemetry. |
| `passwordValidationPolicy` | object | — | Instance-level password rules. |
| `connectionPooling` | object | disabled | Managed connection pooling. |
| `databaseFlags` | map | — | Engine flags. |
| `encryptionKeyName` | `StringValueOrRef` | Google-managed | CMEK (ref → GcpKmsKey). Immutable. |
| `deletionProtection` / `deletionProtectionEnabled` | bool | `false` | Engine-side and API-side delete guards. |
| `masterInstanceName` | `StringValueOrRef` | — | Makes this a read replica of the named primary. Immutable. |
| `rootPassword` | secret | — | Initial admin password (required for SQL Server). Write-only. |

SQL Server-only: `timeZone`, `collation`, `threadsPerCore`, `sqlServerAuditConfig`, `activeDirectoryDomain` — validated per engine before deploy.

## Examples

### Production PostgreSQL on Private IP

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudSql
metadata:
  name: orders-db-prod
spec:
  projectId:
    value: my-gcp-project-123
  instanceName: orders-db-prod
  region: us-central1
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_16
  tier: db-custom-4-15360
  availabilityType: REGIONAL
  network:
    privateNetwork:
      valueFrom:
        kind: GcpVpc
        name: prod-vpc
        fieldPath: status.outputs.network_id
    sslMode: ENCRYPTED_ONLY
  backup:
    enabled: true
    pointInTimeRecoveryEnabled: true
  deletionProtection: true
  deletionProtectionEnabled: true
```

### Read Replica

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudSql
metadata:
  name: orders-db-replica
spec:
  instanceName: orders-db-replica-1
  region: us-central1
  databaseEngine: POSTGRESQL
  databaseVersion: POSTGRES_16
  tier: db-custom-2-7680
  masterInstanceName:
    valueFrom:
      kind: GcpCloudSql
      name: orders-db-prod
      fieldPath: status.outputs.instance_name
```

## Stack Outputs

| Output | Description |
|--------|-------------|
| `instance_name` | The composition key databases/users/replicas reference |
| `connection_name` | `project:region:instance` for the Auth Proxy and connectors |
| `private_ip` / `public_ip` | Instance addresses (empty when the path is not enabled) |
| `self_link` | GCP resource self link |
| `service_account_email` | Instance's Google-managed service account (grant GCS access for imports/exports) |
| `dns_name` | DNS name (PSC-enabled instances) |
| `psc_service_attachment_link` | PSC attachment consumers target with endpoints |

## Related Components

- [GcpCloudSqlDatabase](/docs/catalog/gcp/gcpcloudsqldatabase) — databases inside the instance
- [GcpCloudSqlUser](/docs/catalog/gcp/gcpcloudsqluser) — per-application users
- [GcpServiceNetworkingConnection](/docs/catalog/gcp/gcpservicenetworkingconnection) — private services access for private IP
- [GcpKmsKey](/docs/catalog/gcp/gcpkmskey) — customer-managed encryption
