# KubernetesVelero Terraform Module

Installs Velero from the official Helm chart (`velero` at
`https://vmware-tanzu.github.io/helm-charts`) as a real Helm release. The
typed spec renders into chart values in `locals.tf`
(`local.typed_values`); credential material rides a SECOND,
`sensitive()`-wrapped values document; the `helm_values` escape hatch is
passed as the LAST document the provider merges (Helm `-f` semantics) —
the exact semantic twin of the Pulumi module's `buildHelmValues` +
`mergeMaps`.

## Module Behavior

- **The release name is FIXED to `velero`** — Velero's CRDs and
  node-agent are cluster-scoped and one server owns the backup records in
  the object store; one installation per cluster is an upstream
  constraint. The fixed name also collapses the chart's `velero.fullname`
  to `velero` (the release name contains the chart name), which is what
  makes the derived `velero-server` ServiceAccount name deterministic.
- **CRDs survive uninstall by default** — the chart ships its CRDs in the
  crds/ directory, which Helm installs once and (by Helm's own contract)
  never deletes on uninstall, so backup records outlive the component.
  `crds.upgrade_on_install` (chart default true) drives the chart's
  `upgradeCRDs` job that keeps them current across upgrades;
  `crds.cleanup_on_uninstall` opts into the chart's DESTRUCTIVE
  `cleanUpCRDs` job. Both render only when they differ from the chart
  default.
- **Readiness is verified at install time** — `wait` + `wait_for_jobs` +
  `atomic` + `cleanup_on_fail` with a 600s timeout. A Velero that never
  comes up (a ServiceMonitor rendered without the Prometheus operator
  CRDs, an unpullable plugin image) fails THIS apply instead of surfacing
  later as backups that silently never run.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.

## Backend Rendering

Exactly one `backup_storage` arm (spec-enforced oneof). Each arm renders
the plugin initContainer, the default BackupStorageLocation and the
credential posture — nothing from inactive arms leaks into the values:

- **Plugin initContainer**: `velero/velero-plugin-for-aws:v1.14.2` /
  `-gcp:v1.14.2` / `-microsoft-azure:v1.14.2` by default (the versions
  paired with chart 12.1.0 = Velero 1.18), override via
  `backup_storage.plugin_image`; mounts the shared `plugins` dir at
  `/target` — the shape the chart `values.yaml` documents verbatim.
- **BackupStorageLocation** (`configuration.backupStorageLocation[0]`):
  name `default`, `default: true`, provider + bucket (Azure: the blob
  container), `prefix` when set. `caCert` is an ITEM-LEVEL key (rendered
  under the BSL's `objectStorage.caCert` by the chart template), not a
  `config` entry. The chart template QUOTES every `config` value, so
  config entries are strings — `s3ForcePathStyle` renders as `"true"`.
- **Credential posture** (all XORs spec-enforced):
  - S3 IRSA / GCS workload identity → identity annotation on
    `serviceAccount.server.annotations`, `credentials.useSecret: false`
    (GCS additionally sets the BSL's `config.serviceAccount` — the
    values.yaml-documented workload-identity key)
  - S3 access keys → `cloud` file in the AWS shared-credentials format
    (`[default]` + key id + secret key, documented inline in the chart
    values.yaml); GCS key → the JSON verbatim
  - Azure workload identity → client-id annotation + the
    `azure.workload.identity/use` pod label (the AKS webhook only injects
    tokens into labeled pods) PLUS a `cloud` file with the non-secret
    identifiers (subscription, resource group,
    `AZURE_CLOUD_NAME=AzurePublicCloud`) — the Azure plugin always reads
    the file
  - Azure service principal → the full `AZURE_*` environment-file lines
  - Neither → `credentials.useSecret: false` (ambient node credentials)

## Rendering Quirks

- **`configuration` is a collector** — the BSL/VSL lists AND the
  velero-server CLI flags all live under it; `defaultVolumesToFsBackup`
  is a server flag, so it renders there even though the spec groups it
  with `fs_backup`.
- **VolumeSnapshotLocation** renders while snapshots are enabled (chart
  default): the arm's provider with the provider-required config
  (`region` for aws, `resourceGroup`+`subscriptionId` for azure, nothing
  for gcp — per the values.yaml comments).
- **`schedules` is a MAP keyed by schedule name** — the chart names each
  rendered Schedule `velero-<key>`. The three optional template booleans
  (`includeClusterResources`, `snapshotVolumes`,
  `defaultVolumesToFsBackup`) render presence-aware: unset means "Velero
  decides", which is different from false.
- **Chart-default-true flags render only on opt-out**: `upgradeCRDs`,
  `snapshotsEnabled`, `metrics.enabled`.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types; arm-dependent scalars are picked
  with `one([... if x != null])` / provider-keyed lookup maps, never
  chained per-arm ternaries.
- **Secrets stay out of the visible plan** — the credential `cloud`
  content is the only place secret material appears, isolated in its own
  values document and wrapped with `sensitive()`. `local.typed_values`
  itself never references a secret field.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.velero` | `spec.create_namespace` |
| `helm_release.velero` | always |

## Usage

```bash
planton tofu apply --manifest velero.yaml
```

## Local Development

```bash
terraform init
terraform validate
terraform plan -var-file=terraform.tfvars.json
terraform apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the full variable specification (generated from
the spec proto).

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Kubernetes namespace Velero was installed into |
| `release_name` | Helm release name (fixed `velero` — one installation per cluster) |
| `service_account_name` | Chart-derived `velero-server` ServiceAccount (`velero.serverServiceAccount` helper: `<fullname>-server`, fullname `velero`) — the subject cloud-side keyless bindings are written against |
| `backup_storage_location_name` | Name of the default BackupStorageLocation (`default`) — what Backup/Schedule resources reference through `storageLocation` |

## Parity

Kept in lockstep with the Pulumi module (`iac/pulumi/module/`): same
chart identity and pinned default version (12.1.0 = Velero 1.18.1), same
values rendering (backend arms, credential `cloud` formats, the
configuration collector, the schedules map), same fixed release name,
same outputs.
