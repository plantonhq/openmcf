---
title: "Postgres"
description: "Postgres deployment documentation"
icon: "package"
order: 100
componentName: "kubernetespostgres"
---

# Kubernetes Postgres

Declares a production-grade PostgreSQL cluster reconciled by
CloudNativePG — instances and streaming replication with automated
failover, storage with an optional dedicated WAL volume, PostgreSQL
configuration and synchronous replication, bootstrap (fresh initdb,
restore from backup, physical streaming, or logical import), declarative
roles, continuous WAL archiving plus scheduled base backups to
S3/GCS/Azure-Blob/S3-compatible stores, TLS, and monitoring. One
resource per PostgreSQL cluster; applications connect through the
operator-managed Services, which re-point automatically on failover.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **Cluster** (`postgresql.cnpg.io/v1`, named `metadata.name`) — the
  PostgreSQL cluster; CloudNativePG derives everything from it: instance
  pods (`<name>-1`, ...), the traffic Services (`<name>-rw` / `-ro` /
  `-r`), the credential Secrets (`<name>-app`, `<name>-superuser`), and
  the cluster's ServiceAccount
- **ObjectStore** (`barmancloud.cnpg.io/v1`, when backups are declared)
  — the Barman Cloud backup destination; a second one
  (`<name>-recovery-source`) when the bootstrap restores from a backup
- **ScheduledBackup** (one per declared schedule, named
  `<name>-<schedule>`) — the periodic base backups point-in-time
  recovery needs
- **Credential Secrets** — every declared password/key materializes as a
  deterministic Secret (`<name>-app-provided`, `<name>-role-<role>`,
  `<name>-backup-creds`, ...); keyless backup arms need none

## Prerequisites

- The CloudNativePG operator on the cluster
  (KubernetesCloudNativePgOperator) — with `barman_cloud_plugin.enabled`
  when backups are declared
- For backups: an object store (S3 bucket, GCS bucket, Azure Blob
  container, or an S3-compatible endpoint such as MinIO) — one
  destination path per cluster
- For keyless backups: cloud-side identity (IRSA role, GCP Workload
  Identity binding, or Azure federated credential) written against the
  cluster's own ServiceAccount, wired via `workload_identity`
- For organization-trusted TLS: a cert-manager-issued certificate
  (KubernetesCertificate) referenced by `certificates.server_tls_secret`

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesPostgres
metadata:
  name: orders-db
spec:
  namespace:
    value: orders
  create_namespace: true
  instances: 3
  storage:
    size: 100Gi
  bootstrap:
    initdb:
      database: orders
      data_checksums: true
  backup:
    object_store:
      destination_path: s3://acme-pg-backups/orders-db
      s3:
        region: us-west-2
        keyless: true # IRSA via workload_identity below
    retention_policy: 30d
    schedules:
      - name: daily
        schedule: "0 0 2 * * *" # six-field cron, seconds first
        immediate: true
  workload_identity:
    eks:
      role_arn:
        value: arn:aws:iam::111111111111:role/orders-db-backups
```

The operator brings up one primary and two replicas; applications
connect at the exported `kube_endpoint` with the `<name>-app` Secret's
credentials, WAL archiving starts immediately, and the first base backup
runs on creation.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | Name of the Cluster resource (equals `metadata.name`) |
| `rw_service` | Read-write Service (`<name>-rw`) — the current primary |
| `ro_service` | Read-only Service (`<name>-ro`) — replicas only |
| `r_service` | Any-instance read Service (`<name>-r`) |
| `kube_endpoint` | In-cluster connection host (`<name>-rw.<namespace>.svc.cluster.local:5432`) |
| `port_forward_command` | Workstation access when no exposure is composed |
| `username_secret` / `password_secret` | `{name, key}` of the application credential — the Secret also carries ready-made `uri` / `jdbc-uri` |
| `superuser_secret_name` | Superuser Secret — empty unless superuser access is enabled |

## Next Steps

Compose exposure when the database must be reachable from outside — a
KubernetesService of type LoadBalancer or a Gateway TCP route against
the `rw_service` output (this component never embeds one), with external
hostnames added to `certificates.server_alt_dns_names`. Move to
`scheduling.anti_affinity_type: required` and synchronous replication
for production. Prove recovery early: bootstrap a clone from the
cluster's own backups with `bootstrap.recovery` — a backup you have
never restored is a hope, not a plan.
