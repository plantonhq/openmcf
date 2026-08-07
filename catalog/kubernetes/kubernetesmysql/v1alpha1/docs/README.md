# KubernetesMysql: Research and Design

## Introduction

KubernetesMysql declares one production-grade MySQL cluster reconciled by
the Percona Operator for MySQL based on Percona XtraDB Cluster
(KubernetesPerconaMysqlOperator must be on the cluster; module pinned to
operator release v1.20.0). The spec renders a `pxc.percona.com/v1`
PerconaXtraDBCluster custom resource — Galera synchronous multi-primary
replication, automated full-cluster-crash recovery, HAProxy or ProxySQL
query routing, scheduled XtraBackup backups with point-in-time recovery,
TLS, and declarative users.

## Upstream Architecture

The operator manages a Galera cluster (PXC) behind a proxy layer.
Applications connect through the proxy Service, never a database pod:

- **Database pods** `<name>-pxc-0..N` — every node holds the full dataset;
  a committed transaction is certified on every node before the client
  sees success.
- **Proxy Services** `<name>-haproxy` (writes, port 3306; reads via 3307)
  or `<name>-proxysql` when ProxySQL is chosen. HAProxy is the upstream
  default (3 replicas).
- **System-users Secret** `<name>-secrets` — operator-generated passwords
  for root and internal accounts (key `root` is the admin password).
- **Per-user Secrets** — declarative users with declared passwords get
  `<name>-user-<username>`; the operator watches and rotates on change.

## Credential Materialization

Nothing sensitive ever appears inline in a rendered custom resource:

- `<name>-user-<username>` (key `password`) — declared application users.
- `<name>-backup-<storage>` — S3 keys (`AWS_ACCESS_KEY_ID` /
  `AWS_SECRET_ACCESS_KEY`) or Azure keys (`AZURE_STORAGE_ACCOUNT_NAME` /
  `AZURE_STORAGE_ACCOUNT_KEY`). PVC storages need no credentials.

Names are deterministic so both IaC engines agree byte-for-byte.

## Proxy and Exposure

The proxy oneof chooses HAProxy (default, omitted block) or ProxySQL
(stateful — requires its own PVC). Exposure is composed, never
embedded: `expose_primary` / `expose_replicas` service types and
annotations exist for managed-cloud LoadBalancer recipes; reach the
database from outside by composing a first-class exposure kind against
the exported `primary_service`.

## Backups

Backup storages are a named list (each name unique, referenced by
schedules and PITR by name). Supported backends: S3/S3-compatible
(MinIO, etc.), Azure Blob,
or a PVC (`type: filesystem`). PITR continuously ships binlogs to a
**dedicated** storage — never share it with base backups (upstream
requirement). Schedules use five-field cron.

## Version Pins

| What | Pin | Notes |
|---|---|---|
| Operator CR contract | `crVersion: 1.20.0` | Module constant |
| Database image | `percona/percona-xtradb-cluster:8.4.8-8.1` | Overridden by `image_name` |
| HAProxy | `percona/haproxy:2.8.18-1` | Module constant |
| ProxySQL | `percona/proxysql2:2.7.3-1.3` | Module constant |
| XtraBackup | `percona/percona-xtrabackup:8.4.0-5.1` | Module constant |
| Log collector | `percona/fluentbit:5.0.6-1` | Module constant |
| Upgrade apply | `disabled` | Version changes via `image_name` only |

## Unsafe Opt-In

Galera quorum requires three nodes. The `unsafe` block opts into
development topologies the operator otherwise rejects: single-node
clusters, disabled TLS, single proxy replica, backups against unhealthy
clusters.

## IaC Twins

Pulumi (typed crd2pulumi SDK) and Terraform (kubectl_manifest + null-
prune locals) render identical CR bodies and credential Secrets. Keep
`locals.go` / `cluster.go` and `locals.tf` in lockstep.
