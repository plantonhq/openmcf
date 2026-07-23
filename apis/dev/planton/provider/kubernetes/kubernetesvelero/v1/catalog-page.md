# Kubernetes Velero

Installs Velero — cluster backup and disaster recovery — from the
official Helm chart, with a typed spec over the chart's meaningful
configuration surface. Velero backs up Kubernetes resources (and
optionally volume data) to an object store, restores them on demand, and
runs scheduled backups — the recovery path for cluster loss, migration,
and fat-finger deletions. The backend oneof picks where backups land
(S3/S3-compatible including in-cluster MinIO, Google Cloud Storage, or
Azure Blob) and loads the matching provider plugin. One installation per
cluster.

## What Gets Created

- **Namespace** (optional) — the installation namespace, created and
  owned when `create_namespace` is set (`velero` is the upstream
  convention)
- **Helm Release** (`velero`) — the Velero server Deployment with the
  backend's provider plugin as an init container, the default
  BackupStorageLocation (named `default`), the default
  VolumeSnapshotLocation (unless disabled), the Velero CRDs (Backup,
  Restore, Schedule, ... — surviving uninstall by Helm's own contract
  for crds/-directory CRDs), any declared Schedule resources, and — when
  file-system backup or CSI data movement is enabled — the node-agent
  DaemonSet

## Prerequisites

- An object store to land backups in: an S3 bucket, an S3-compatible
  endpoint (MinIO, Ceph RGW, ...), a GCS bucket, or an Azure Blob
  container
- Cloud-side identity for the keyless postures: an IRSA role, GCP
  Workload Identity binding, or Azure federated credential written
  against the `velero-server` service account; S3-compatible endpoints
  authenticate with access keys
- For CSI snapshots: a snapshot-capable CSI driver; clusters without
  snapshot support disable snapshots and carry volume data via
  file-system backup
- With `prometheus.service_monitor`: the Prometheus operator CRDs — the
  release fails to install without them

## Quick Start

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesVelero
metadata:
  name: velero
spec:
  namespace:
    value: velero
  createNamespace: true
  backupStorage:
    s3:
      bucket: my-cluster-backups
      region: us-west-2
      irsaRoleArn: arn:aws:iam::111111111111:role/velero
    prefix: prod-cluster
  schedules:
    - name: daily-cluster
      schedule: "0 2 * * *"
      ttl: 720h
```

The server validates the storage location on startup; backups then run
nightly and are garbage-collected after 30 days.

## Stack Outputs

| Output | Description |
|---|---|
| `namespace` | Installation namespace |
| `release_name` | Helm release name (always `velero`) |
| `service_account_name` | Always `velero-server` — the subject cloud-side keyless bindings are written against |
| `backup_storage_location_name` | Always `default` — what Backup and Schedule resources reference through storageLocation |

## Next Steps

Decide the volume-data path: `enable_csi` on clusters with
snapshot-capable CSI drivers (add `default_snapshot_move_data` to make
backups portable across clusters), or `fs_backup.deploy_node_agent` for
everything else. Prefix per cluster when several clusters share one
bucket. Prove the loop early — back up a namespace, delete it, restore
it — and keep a DR-standby cluster reading the same store with
`server.restore_only_mode` when recovery time matters. Removing the
component never deletes the backup records or the stored objects by
default.
