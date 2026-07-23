# Kubernetes MySQL

Declares a production-grade MySQL cluster reconciled by the Percona
Operator for MySQL (XtraDB Cluster / Galera) — synchronous multi-primary
replication with automated recovery, HAProxy or ProxySQL query routing,
scheduled XtraBackup backups with point-in-time recovery, TLS, and
declarative users. One resource per MySQL cluster; applications connect
through the proxy Service (`<name>-haproxy` or `<name>-proxysql`).

## What Gets Created

- **Namespace** (optional) — created and owned when `create_namespace`
  is set
- **PerconaXtraDBCluster** (`pxc.percona.com/v1`, named
  `metadata.name`) — the Galera cluster; the operator derives database
  pods (`<name>-pxc-0..N`), proxy Services, and the system-users Secret
  from it
- **Credential Secrets** — every declared password/key materializes as a
  deterministic Secret (`<name>-user-<username>`, `<name>-backup-
  <storage>`); PVC backup storages need none

## Prerequisites

- The Percona MySQL operator on the cluster
  (KubernetesPerconaMysqlOperator) — it must watch the database's
  namespace (the default posture watches its OWN namespace: install
  the operator there, or widen its watch)
- For backups: S3/S3-compatible object storage, Azure Blob, or a PVC —
  dedicate one storage to PITR when enabled
- For organization-trusted TLS: a cert-manager ClusterIssuer or Issuer
  referenced by `tls.issuer`

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesMysql
metadata:
  name: orders-db
spec:
  namespace:
    value: orders
  create_namespace: true
  instances: 3
  storage:
    size: 100Gi
  users:
    - name: app
      dbs: [orders]
      hosts: ["%"]
      grants: [SELECT, INSERT, UPDATE, DELETE]
      password: <set-a-strong-password>
```

The operator brings up three Galera nodes and three HAProxy replicas;
applications connect at the exported `kube_endpoint` using the root
password from `<name>-secrets` (or a declared user's Secret).

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the cluster runs in |
| `cluster_name` | PerconaXtraDBCluster name (equals `metadata.name`) |
| `primary_service` | Write Service (`<name>-haproxy` or `<name>-proxysql`) |
| `replicas_service` | Read Service (`<name>-haproxy-replicas`) — empty unless HAProxy replicas Service is enabled |
| `kube_endpoint` | In-cluster connection host (`<primary_service>.<namespace>.svc.cluster.local:3306`) |
| `port_forward_command` | Workstation access when no exposure is composed |
| `root_password_secret` | `{name, key}` of the operator-managed root credential (`<name>-secrets`, key `root`) |

## Next Steps

Compose exposure when the database must be reachable from outside — a
KubernetesService of type LoadBalancer or a Gateway TCP route against
the `primary_service` output (this component never embeds one). Move to
three database nodes and three proxy replicas for production Galera
quorum. Declare backups with a nightly schedule and prove restore early
— a backup you have never restored is a hope, not a plan.
