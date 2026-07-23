---
title: "GKE GCS Workload Identity"
description: "This preset installs Velero on GKE with backups landing in a Google Cloud Storage bucket reached keylessly via Workload Identity, and the kopia-based node-agent deployed for file-system backup of..."
type: "preset"
rank: "02"
presetSlug: "02-gke-gcs-workload-identity"
componentSlug: "velero"
componentTitle: "Velero"
provider: "kubernetes"
icon: "package"
order: 2
---

# GKE GCS Workload Identity

This preset installs Velero on GKE with backups landing in a Google
Cloud Storage bucket reached keylessly via Workload Identity, and the
kopia-based node-agent deployed for file-system backup of volume data.
No service-account key is stored anywhere in the manifest.

## When to Use

- GKE clusters that need backup / disaster recovery to GCS
- Clusters whose volume data should travel via file-system backup (works
  on any volume type) rather than CSI snapshots

## Key Configuration Choices

- **GCS + Workload Identity
  (`workloadIdentityServiceAccountEmail`)** — the keyless posture on
  GKE: Velero's Kubernetes service account is annotated with the GCP
  service account, which needs storage permissions on the bucket and a
  WI binding to Velero's KSA. The email in the YAML is a placeholder
  value — the field requires the `…gserviceaccount.com` format, so
  replace the whole address with your GSA
- **Node-agent (`fsBackup.deployNodeAgent: true`)** — the kopia-based
  DaemonSet reads volume data directly from pods and works on ANY
  volume type; it runs privileged host-path mounts of the kubelet pod
  directory (that is how it reads volume data)
- **`defaultVolumesToFsBackup: false`** — volumes are backed up per-pod
  via the `backup.velero.io/backup-volumes` annotation; flip to `true`
  to fs-backup every pod volume without annotations
- **`prefix: <cluster-name>`** — one bucket, one prefix per cluster
- **Nightly schedule with `ttl: 720h`** and **CRDs kept on uninstall**
  (chart defaults) — same DR-safety posture as the S3 preset

## Placeholders to Replace

| Placeholder                                          | Description                                | Where to Find                            |
| ---------------------------------------------------- | ------------------------------------------ | ---------------------------------------- |
| `<velero-backups-bucket>`                            | GCS bucket for backups                     | GCS console or `GcpGcsBucket` outputs    |
| `velero-backups@my-project.iam.gserviceaccount.com`  | GSA email for Workload Identity (placeholder value — replace entirely) | IAM console — your Velero service account |
| `<cluster-name>`                                     | Per-cluster prefix within the bucket       | Your cluster naming                      |

## Related Presets

- **01-eks-s3-irsa** — EKS clusters backing up to S3 with CSI snapshots
- **03-minio-self-contained** — S3-compatible / on-prem stores (MinIO)
