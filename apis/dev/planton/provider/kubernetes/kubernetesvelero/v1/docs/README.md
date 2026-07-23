# KubernetesVelero: Research and Design

## Introduction

Velero is the cluster's recovery path: it backs up Kubernetes resources
(and optionally volume data) to an object store, restores them on
demand, and runs scheduled backups — covering cluster loss, migration,
and fat-finger deletions. This component installs it from the official
Helm chart (`velero` at `https://vmware-tanzu.github.io/helm-charts`;
pinned default chart 12.1.0, which ships Velero 1.18.1 — chart and app
versions move separately, and the chart pin governs). Each object-store
backend loads the official provider plugin at the version paired with
that release: `velero-plugin-for-aws` / `-gcp` /
`-microsoft-azure` at v1.14.2 for Velero 1.18.

## Upstream Architecture

An installation is:

1. **The Velero server Deployment** — the controller reconciling
   Backup, Restore, Schedule, BackupStorageLocation, and
   VolumeSnapshotLocation resources against the object store.
2. **A provider plugin init container** — the server discovers plugin
   binaries from a shared `plugins` directory the init container
   populates at `/target`; the active backend arm renders exactly one.
3. **Optionally, the node-agent DaemonSet** — required for file-system
   backup AND for CSI snapshot data movement. It runs privileged
   host-path mounts of the kubelet pod directory: that is how it reads
   volume data directly from pods.

One installation per cluster: the CRDs and node-agent are
cluster-scoped, and one server owns the backup records in the store. The
release name is fixed to `velero`, which also pins the chart's derived
names — the fullname collapses to `velero` (the release name contains
the chart name), making the server's ServiceAccount deterministically
`velero-server` (the chart's serverServiceAccount helper appends
`-server`). That determinism is why the name is a stack output: it is
the subject every cloud-side keyless binding is written against.

## Engine vs Backup Declarations

This component installs the engine and, optionally, the RECURRING
schedules — declared in the spec because they are operating posture, not
day-2 actions (the chart renders each as a Schedule resource, named
`velero-<schedule name>`, from its schedules map). Ad-hoc Backups and
Restores remain day-2 operations against the installed server. The
split matters for lifecycle: "back up everything nightly, keep 30 days"
belongs to the installation; "restore namespace X from Tuesday" does
not.

## CRD Lifecycle: DR Safety by Helm's Own Contract

The chart ships its CRDs in the crds/ directory, which Helm installs on
first install and — by Helm's own contract for crds/-directory CRDs —
NEVER upgrades or deletes on its own. The chart compensates in both
directions, and the spec types both knobs:

- **`crds.upgrade_on_install`** (chart default true) runs the chart's
  CRD-upgrade job, an init job that re-applies the pinned CRDs on
  install/upgrade — what keeps them current across chart upgrades.
- **`crds.cleanup_on_uninstall`** (chart default false) is the
  DESTRUCTIVE inverse: deleting the CRDs — and with them every Backup /
  Restore / Schedule / BackupStorageLocation record in the cluster — on
  uninstall. Upstream warns the switch is meant for CI systems, not
  production. Either way the stored objects in the bucket are never
  touched: after a full reinstall pointing at the same bucket, Velero
  re-syncs the backup records from the store.

The keep-by-default posture is the DR-safety property: removing the
component never destroys the record of what can be restored.

## The Backend Oneof: Three Arms, One Contract

`backup_storage.backend` is a required oneof; each arm configures BOTH
the default BackupStorageLocation (named `default`, marked as the
default location) and the matching provider plugin init container.
`prefix` is item-level on the BSL — the multi-cluster pattern: one
bucket, one prefix per cluster.

- **`s3`** — real AWS S3, or ANY S3-compatible store via `s3_url` (the
  endpoint override; `http://minio.minio.svc:9000` for in-cluster MinIO
  — the self-contained arm) with `force_path_style` (bucket in the path
  instead of the subdomain — required by MinIO and most compatibles).
  `region` is required (compatibles accept a conventional value;
  "minio" by convention). `kms_key_id` covers real-S3 server-side
  encryption; `ca_cert` (an item-level BSL key, not a config entry)
  verifies self-signed endpoints.
- **`gcs`** — a GCS bucket; no VolumeSnapshotLocation config is needed
  for GCP (the chart lists only optional keys).
- **`azure_blob`** — storage account + container + resource group +
  subscription; Velero's generic "bucket" is the blob CONTAINER on
  Azure, and the VolumeSnapshotLocation carries resourceGroup +
  subscriptionId.

Credential posture per arm, spec-enforced as at-most-one:

- **S3**: IRSA (the role-arn annotation on `velero-server`; no Secret
  exists at all) XOR access keys (materialized into the credentials
  Secret in the AWS plugin's `cloud` shared-credentials format). A
  second CEL rule encodes a subtle truth: an S3-compatible endpoint
  authenticates with access keys — IRSA only mints AWS credentials.
- **GCS**: Workload Identity email (the `iam.gke.io/gcp-service-account`
  annotation, plus the BSL's `serviceAccount` config pointing the plugin
  at the same identity) XOR the JSON key file as the `cloud` file.
- **Azure Blob**: workload identity (the
  `azure.workload.identity/client-id` annotation plus the
  `azure.workload.identity/use` pod label — the AKS webhook only injects
  the federated token into labeled pods; a `cloud` file still carries
  the non-secret subscription/resource-group/cloud-name parameters, per
  the Azure plugin's own contract) XOR a service principal (AZURE_*
  lines in the credentials Secret). Workload identity requires the
  client ID — the federated credential is addressed by it (CEL-enforced).

Neither posture set means ambient node credentials (instance profile,
GCE default service account, managed identity).

## Volume Data: CSI Snapshots vs File-System Backup

Volume data travels one of two ways, and the spec models both:

- **CSI snapshots** (`volume_snapshots.enable_csi` — the server's
  EnableCSI feature flag): Velero snapshots PVCs through the CSI
  snapshot API instead of provider-native instance snapshots — the
  modern path on EKS/GKE/AKS with CSI drivers.
  `default_snapshot_move_data` engages Velero's data mover: snapshot
  data moves INTO the backup store instead of staying in the cloud
  provider, making backups portable across clusters and regions. It
  rides CSI snapshots and the node-agent (CEL-enforced pairing on the
  CSI side).

  CSI snapshots have two cluster prerequisites Velero does NOT install:
  the external snapshot controller (a managed add-on on EKS/GKE/AKS,
  often absent on self-managed clusters) and a `VolumeSnapshotClass`
  for each volume's CSI driver labeled
  `velero.io/csi-volumesnapshot-class: "true"` — that label is how
  Velero selects the class. Without them a backup still reports
  Completed while the volumes were never snapshotted (they ride
  fs-backup if enabled, or carry no data at all). After the first
  CSI-enabled backup, confirm the backup's resource list actually
  contains VolumeSnapshot entries before trusting it for recovery.
- **File-system backup** (`fs_backup` — the kopia uploader through the
  node-agent DaemonSet): reads volume data directly from pods, working
  on ANY volume type including clusters without snapshot support.
  `default_volumes_to_fs_backup` backs up ALL pod volumes without
  per-pod `backup.velero.io/backup-volumes` annotations, and requires
  the node-agent (CEL-enforced — the agent is what reads volume data).
  The agent must run on EVERY node whose volumes need backup — its
  tolerations field exists precisely for dedicated-pool taints.

`volume_snapshots.enabled` (chart default true) controls the default
VolumeSnapshotLocation itself; clusters without any snapshot provider
(kind, bare metal, the MinIO posture) disable it and carry volume data
via file-system backup.

## Server Posture

`server` types the velero-server flags every real installation ends up
setting: `default_backup_ttl` (Velero default 720h),
`default_item_operation_timeout` (4h — large-volume uploads),
`garbage_collection_frequency` (1h), log level/format, and
`restore_only_mode` — the DR-standby stance: a second cluster reading
the primary's store with backups, schedules, and GC blocked, ready to
restore. `prometheus.enabled` matches the chart default (true; rendered
only on explicit opt-out — DR you cannot observe is DR you cannot
trust); the ServiceMonitor is opt-in and requires the Prometheus
operator CRDs.

## Typed Surface vs Escape Hatch

The typed spec covers namespace and lifecycle, chart version, CRD
lifecycle, the three backend arms with their credential postures, both
volume-data paths, schedules, server tuning, deployment
sizing/scheduling, and telemetry.

Deliberately unmodeled as typed fields (all reachable via
`helm_values`):

- **Additional BackupStorageLocations / VolumeSnapshotLocations** — the
  spec models the DEFAULT pair; secondary locations are an advanced
  posture riding the same chart lists
- **Backup hooks and per-workload annotations** — properties of the
  workloads, not the installation
- **Image overrides, RBAC tuning, extra volumes, plugin-config
  ConfigMaps** — the chart's operational long tail
- **`configuration.backupSyncPeriod` and the rest of the server's flag
  long tail** — the chart's configuration passthrough carries them

## Install Semantics

Both engines install a REAL Helm release, atomically, with cleanup on
fail and a 600s timeout. Typed values render first, the credentials
document (when a declared-credential arm is active) and `helm_values`
merge after with Helm `-f` semantics. Secret material is handled
coarsely on purpose: when any arm carries actual secret content (an S3
secret key, a GCP key JSON, an Azure client secret) the rendered values
are masked in state wholesale — nothing secret can leak through a missed
path; the Azure workload-identity `cloud` file carries only non-secret
identifiers and does not trigger masking. The module (not Helm) owns
namespace creation via `create_namespace`.

## Outputs

`namespace`, `release_name` (fixed `velero`), `service_account_name`
(the chart-derived `velero-server` — the subject every cloud-side
keyless binding is written against), and `backup_storage_location_name`
(fixed `default` — what Backup and Schedule resources reference through
storageLocation).

## E2E

The behavioral facts are properties of the platform, not of any one test
run:

- The readiness proof is the BackupStorageLocation reporting Available —
  the server validated it can reach the store with the configured
  credentials; a wrong bucket, endpoint, or key surfaces here.
- The MinIO posture is the deterministic self-contained proof: an
  in-cluster S3-compatible store, path-style addressing, access keys,
  snapshots off, volume data via file-system backup — the full
  backup/restore cycle with no cloud account.
- A backup/restore round trip (back up a namespace, delete it, restore
  it) is the behavioral proof; with the node-agent deployed, volume
  contents ride along.
- Uninstall keeps the CRDs and every backup record by Helm's own
  crds/-directory contract; the stored objects survive regardless, and
  a reinstall against the same bucket re-syncs the records.
- The ServiceMonitor arm fails the release on clusters without the
  Prometheus operator CRDs, by design.
