# KubernetesVelero Pulumi Module

Installs Velero from the official Helm chart (`velero` at
`https://vmware-tanzu.github.io/helm-charts`) as a real Helm release. The
typed spec renders into chart values in `module/values.go`; the
`helm_values` escape hatch merges LAST over them with Helm `-f` semantics
(maps deep-merge, later document wins, lists replace) — the exact semantic
twin of the Terraform module's `helm_release` with
`values = [typed, credentials, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true; otherwise the namespace must already
   exist
2. **Helm Release** — named `velero`, FIXED: Velero's CRDs and node-agent
   are cluster-scoped and one server owns the backup records in the
   object store, so one installation per cluster is an upstream
   constraint and the release name never derives from `metadata.name`.
   The fixed name also collapses the chart's `velero.fullname` to
   `velero` (the release name contains the chart name), which is what
   makes every derived name deterministic.

## Backend Rendering (the heart of the module)

Exactly one `backup_storage.backend` arm (spec-enforced oneof). Each arm
renders three things — nothing from inactive arms leaks into the values:

- **Plugin initContainer** (`initContainers[0]`): the provider plugin
  image (default `velero/velero-plugin-for-aws:v1.14.2` /
  `-gcp:v1.14.2` / `-microsoft-azure:v1.14.2`, override via
  `backup_storage.plugin_image`) mounting the shared `plugins` dir at
  `/target` — the shape the chart's `values.yaml` documents verbatim.
- **Default BackupStorageLocation**
  (`configuration.backupStorageLocation[0]`): name `default`,
  `default: true`, the arm's provider and bucket (Azure: the blob
  container), `prefix` when set. `caCert` is an ITEM-LEVEL key (the chart
  template renders it under the BSL's `objectStorage.caCert`), not a
  `config` entry. Provider-specific `config` keys: `region`/`s3Url`/
  `s3ForcePathStyle`/`kmsKeyId` (aws — the template QUOTES every config
  value, so the path-style flag renders as the string `"true"`),
  `serviceAccount` (gcp workload identity), `resourceGroup`/
  `storageAccount`/`subscriptionId` (azure).
- **Credential posture**:
  - S3 IRSA / GCS workload identity → identity annotation on
    `serviceAccount.server.annotations`, `credentials.useSecret: false` —
    no Secret exists at all
  - S3 access keys → `credentials.secretContents.cloud` in the AWS
    shared-credentials format (`[default]` + key id + secret key, as the
    chart values.yaml documents); GCS key → the JSON verbatim
  - Azure workload identity → the client-id annotation + the
    `azure.workload.identity/use` pod label (the AKS webhook only injects
    tokens into labeled pods) PLUS a `cloud` file with the non-secret
    identifiers (`AZURE_SUBSCRIPTION_ID`, `AZURE_RESOURCE_GROUP`,
    `AZURE_CLOUD_NAME=AzurePublicCloud`) — unlike AWS/GCP the Azure
    plugin always reads the file
  - Azure service principal → the full `AZURE_*` environment-file lines
    including the client secret
  - Neither → `credentials.useSecret: false` (ambient node credentials)

## Other Rendering Notes

- **Chart-default-true flags render only on opt-out**: `upgradeCRDs`,
  `snapshotsEnabled`, `metrics.enabled`. The chart's crds/-directory CRDs
  are Helm-native keep-on-uninstall — backup records survive removing the
  component unless `crds.cleanup_on_uninstall` explicitly opts into the
  chart's destructive cleanup job.
- **VolumeSnapshotLocation** renders while snapshots are enabled
  (default): the arm's provider with the provider-required config
  (`region` for aws, `resourceGroup`+`subscriptionId` for azure, nothing
  for gcp — per the chart values.yaml comments).
- **`schedules` is a MAP keyed by schedule name** — the chart names each
  rendered Schedule object `velero-<key>` (`velero.fullname` + key). The
  three optional template booleans (`includeClusterResources`,
  `snapshotVolumes`, `defaultVolumesToFsBackup`) render presence-aware:
  unset means "Velero decides", which is different from false.
- **`configuration.defaultVolumesToFsBackup` is a server flag**, not a
  nodeAgent key — it lives under `configuration` even though the spec
  groups it with `fs_backup`.
- **Secret masking is coarse**: when the values carry actual secret
  material (S3 secret key, GCP key JSON, Azure client secret) the WHOLE
  values map is wrapped as a Pulumi secret — nothing can leak through a
  missed path.

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and waits (600s
timeout) for the server — and the node-agent DaemonSet / upgradeCRDs job
when enabled — to become ready. A Velero that never comes up (a
ServiceMonitor rendered without the Prometheus operator CRDs, an
unpullable plugin image) fails THIS deploy with a readiness timeout
instead of surfacing later as backups that silently never run.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Kubernetes namespace Velero was installed into |
| `release_name` | Helm release name (fixed `velero` — one installation per cluster) |
| `service_account_name` | The chart-derived `velero-server` ServiceAccount (helper `velero.serverServiceAccount`: `<fullname>-server` with fullname `velero`) — the subject cloud-side keyless bindings (IRSA trust policies, GCP WI bindings, Azure federated credentials) are written against |
| `backup_storage_location_name` | Name of the default BackupStorageLocation (`default`) — what Backup/Schedule resources reference through `storageLocation` |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → Helm release (with coarse secret masking)
  → output exports
- `module/values.go`: typed-spec → chart values rendering (CRD lifecycle,
  the three backend arms with plugin/BSL/credentials, volume snapshots
  and the VSL, file-system backup, the schedules map, server flags,
  deployment sizing, telemetry), escape-hatch merge
- `module/locals.go`: resolved namespace and chart version — kept in
  lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity, pinned default chart version
  (12.1.0 = Velero 1.18.1), the fixed release name, the chart-derived
  server service-account name, plugin image defaults
