# MinIO Self-Contained

This preset installs Velero against an S3-COMPATIBLE object store — an
in-cluster MinIO by default — with no cloud account involved anywhere:
the store is reached by endpoint URL with path-style addressing and
access keys, volume snapshots are disabled (an S3-compatible store
provides no snapshot provider), and ALL volume data travels through
kopia file-system backup. This is the self-contained / on-prem posture:
it works on bare metal, kind, and air-gapped clusters exactly as it does
on cloud ones.

## When to Use

- On-prem / bare-metal clusters without a cloud object store
- Air-gapped environments where backups must stay inside the cluster or
  data center
- Any cluster backing up to an S3-compatible store (MinIO, Ceph RGW,
  DigitalOcean Spaces, Backblaze B2, ...)

## Key Configuration Choices

- **`s3Url` + `forcePathStyle: true`** — the S3-compatible arm: the
  endpoint URL points at the store (in-cluster MinIO service DNS here),
  and path-style addressing (bucket in the path instead of the
  subdomain) is required by MinIO and most S3-compatible stores
- **`region: minio`** — S3-compatible stores use whatever region value
  they expect; MinIO accepts any, `minio` by convention
- **Access keys, not IRSA** — an S3-compatible endpoint authenticates
  with access keys (IRSA only mints AWS credentials; the spec enforces
  this pairing). `accessKeyId: velero` and `<minio-secret-key>` are
  PLACEHOLDER values — replace both with your MinIO credentials (access
  key = username, secret key = password). The secret key is stored as a
  managed secret, not in plain text
- **`volumeSnapshots.enabled: false`** — no snapshot provider exists for
  the store, so the default VolumeSnapshotLocation is not created
- **`fsBackup` with `defaultVolumesToFsBackup: true`** — the node-agent
  DaemonSet (kopia) reads volume data directly from pods and works on
  ANY volume type; with the default flipped on, every pod volume is
  backed up without per-pod annotations
- **Self-signed TLS?** — if your endpoint uses https with a private CA,
  set `backupStorage.s3.caCert` (base64 PEM bundle)

## Placeholders to Replace

| Placeholder                   | Description                                        | Where to Find                 |
| ----------------------------- | -------------------------------------------------- | ----------------------------- |
| `velero-backups`              | Bucket name in the store (create it first)         | Your MinIO console            |
| `http://minio.minio.svc:9000` | Endpoint URL of your S3-compatible store           | Your MinIO/store installation |
| `velero`                      | Access key id (MinIO username) — placeholder value | Your MinIO credentials        |
| `<minio-secret-key>`          | Secret access key (MinIO password)                 | Your MinIO credentials        |
| `<cluster-name>`              | Per-cluster prefix within the bucket               | Your cluster naming           |

## Related Presets

- **01-eks-s3-irsa** — real AWS S3 with keyless IRSA and CSI snapshots
- **02-gke-gcs-workload-identity** — GCS with Workload Identity
