---
title: "MongoDB"
description: "MongoDB deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesmongodb"
---

# Kubernetes MongoDB

Declares a production-grade MongoDB cluster reconciled by the Percona
Operator for MongoDB — replica sets with automated failover, optional
sharding (mongos routers + config servers), scheduled backups with
point-in-time recovery via Percona Backup for MongoDB, TLS, and
declarative users. The server is Percona Server for MongoDB, a fully
MongoDB-compatible open-source distribution: every driver, tool, and
query works unchanged. One resource per MongoDB cluster; applications
connect through the operator-managed Services.

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **PerconaServerMongoDB** (`psmdb.percona.com/v1`, named
  `metadata.name`) — the MongoDB cluster; the operator derives
  everything from it: member pods (`<name>-<rs>-0..N`), the
  per-replica-set headless Services (`<name>-<rs>`), the mongos Service
  (`<name>-mongos`, sharding only), and the system-users Secret
  (`<name>-secrets`, operator-generated passwords for the built-in
  accounts)
- **Credential Secrets** — every declared password/key materializes as
  a deterministic Secret (`<name>-user-<username>` for declarative
  users, `<name>-backup-<storage>` for backup-store credentials);
  keyless backup arms need none

## Prerequisites

- The Percona MongoDB operator on the cluster
  (KubernetesPerconaMongoOperator) — it must watch the database's
  namespace (the default posture watches its OWN namespace: install
  the operator there, or widen its watch)
- For backups: S3, an S3-compatible endpoint (MinIO, Ceph RGW, ...),
  GCS, or Azure Blob storage
- For organization-trusted TLS: a cert-manager ClusterIssuer or Issuer
  referenced by `tls.issuer`

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesMongodb
metadata:
  name: orders-db
spec:
  namespace:
    value: percona-mongo # where the operator watches
  replica_sets:
    - name: rs0
      size: 3
      storage:
        size: 50Gi
  users:
    - name: app
      roles:
        - name: readWrite
          db: orders
```

The operator brings up a three-member replica set with automated
failover; applications connect at the exported `kube_endpoint` (with
`?replicaSet=rs0` so the driver follows failovers), using the declared
user's operator-managed Secret or the admin credential from
`<name>-secrets`.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | PerconaServerMongoDB name (equals `metadata.name`) |
| `service` | The Service applications connect to — `<name>-mongos` when sharded, otherwise the first replica set's headless Service (`<name>-<rs>`) |
| `kube_endpoint` | In-cluster connection host (`<service>.<namespace>.svc.cluster.local:27017`) |
| `replica_set` | The first replica set's name (the driver's `replicaSet` parameter) — empty for sharded clusters |
| `port_forward_command` | Workstation access when no exposure is composed |
| `admin_password_secret` | `{name, key}` of the database-admin password (the operator-managed `<name>-secrets` Secret, key `MONGODB_DATABASE_ADMIN_PASSWORD`) |

## Next Steps

Compose exposure when the database must be reachable from outside — a
first-class exposure kind against the exported service name (this
component never embeds one; the per-set `expose` block exists for the
managed-cloud LoadBalancer recipes). Move to `tls.mode: requireTLS`
for production. Declare a backup block with a nightly task and
point-in-time recovery, and prove restore early — a backup you have
never restored is a hope, not a plan. Shard only when a single replica
set's write volume or working set no longer fits: enable `sharding`
and every declared replica set becomes a shard behind mongos.
