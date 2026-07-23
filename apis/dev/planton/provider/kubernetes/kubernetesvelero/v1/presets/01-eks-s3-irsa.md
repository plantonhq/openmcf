# EKS S3 IRSA

This preset installs Velero on EKS in the standard production posture:
backups land in an S3 bucket reached keylessly via IRSA, volume data is
captured through CSI snapshots (the modern path with the EBS CSI
driver), and a nightly schedule keeps 30 days of backups. Velero is one
installation per cluster — its CRDs and node-agent are cluster-scoped,
and one server owns the backup records in the store.

## When to Use

- Any production EKS cluster that needs backup / disaster recovery
- The 30-second choice: this is the standard first Velero installation
  on AWS

## Key Configuration Choices

- **S3 + IRSA (`backupStorage.s3.irsaRoleArn`)** — the keyless posture
  on EKS: no stored access keys; the role's trust policy must allow the
  cluster's OIDC provider and Velero's service account
- **`prefix: <cluster-name>`** — the multi-cluster pattern: one bucket,
  one prefix per cluster, so several clusters share a store without
  colliding
- **CSI snapshots (`volumeSnapshots.enableCsi: true`)** — PVCs snapshot
  through the CSI snapshot API instead of provider-native instance
  snapshots; requires a snapshot-capable CSI driver (the EBS CSI driver
  on EKS). The default VolumeSnapshotLocation stays enabled (chart
  default)
- **Nightly schedule with `ttl: 720h`** — backups at 02:00 every night,
  garbage-collected after 30 days (Velero's default TTL, stated
  explicitly)
- **`crds.cleanupOnUninstall: false`** (chart default) — backup records
  survive removing the component; upstream reserves the cleanup switch
  for CI systems. The stored objects in the bucket are never touched
  either way
- **`namespace: velero` + `createNamespace: true`** — the upstream
  convention, in a namespace this resource creates and owns

## Placeholders to Replace

| Placeholder                                     | Description                                 | Where to Find                       |
| ----------------------------------------------- | ------------------------------------------- | ----------------------------------- |
| `<velero-backups-bucket>`                       | S3 bucket for backups                       | S3 console or `AwsS3Bucket` outputs |
| `<aws-region>`                                  | Region of the bucket                        | Your AWS account                    |
| `arn:aws:iam::123456789012:role/velero-backups` | IRSA role ARN — replace account id and name | IAM console                         |
| `<cluster-name>`                                | Per-cluster prefix within the bucket        | Your cluster naming                 |

## Related Presets

- **02-gke-gcs-workload-identity** — GKE clusters backing up to GCS
- **03-minio-self-contained** — S3-compatible / on-prem stores (MinIO)
