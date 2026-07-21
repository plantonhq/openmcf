# Kubernetes Persistent Volume Claim

Deploys a Kubernetes PersistentVolumeClaim — the durable-disk primitive — to a target cluster through a single declarative manifest, covering the complete `core/v1` surface: access modes, storage requests and limits, StorageClass selection (including the empty-vs-absent distinction), volume mode, static binding, volume selectors, and data sources (clone/snapshot-restore). The IaC module handles label merging, namespace resolution, and API defaults automatically — and never blocks a deploy waiting for the claim to bind.

## What Gets Created

When you deploy a KubernetesPersistentVolumeClaim resource, Planton provisions:

- **PersistentVolumeClaim** — a `core/v1` PersistentVolumeClaim requesting storage from the cluster; the cluster binds it to a PersistentVolume that satisfies the request
- **Labels** — standard Planton tracking labels merged with any user-provided labels
- **Annotations** — user-provided annotations applied to the claim metadata

**A claim under a `wait_for_first_consumer` StorageClass stays Pending until a pod uses it — correct behavior, not an error.** Neither IaC engine waits for the claim to reach Bound, so deploys never hang waiting for a consumer that arrives later.

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, or a `KubernetesNamespace` resource referenced from `spec.namespace` so both deploy in one run
- **A StorageClass** — the cluster default (when the spec names none), a named class, or a `KubernetesStorageClass` resource referenced from `spec.storage_class_name`

## Quick Start

Create a file `pvc.yaml` — the simplest useful claim, using the cluster's default StorageClass:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPersistentVolumeClaim
metadata:
  name: app-data
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPersistentVolumeClaim.app-data
spec:
  namespace:
    value: backend
  name: app-data
  storage_request: 10Gi
```

Deploy:

```shell
planton apply -f pvc.yaml
```

A workload mounts the claim by referencing `app-data` in its volume mounts. If the default class binds on first consumer, the claim shows Pending until that workload's pod schedules — expected.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the claim (`metadata.name` in the cluster) — the value workload volume mounts reference. | 1–253 characters, matches `^[a-z0-9]([a-z0-9.-]{0,251}[a-z0-9])?$` |
| `spec.storage_request` | `string` | Requested storage as a Kubernetes quantity (e.g. `"10Gi"`, `"500Mi"`). Growing later requires the class to allow expansion; Kubernetes never shrinks. | Kubernetes quantity format |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespace` | `StringValueOrRef` | `default` | Namespace the claim lives in. Accepts a literal name (`{ value: my-namespace }`) or a reference to a `KubernetesNamespace` resource. |
| `spec.access_modes` | `list(string)` | `["ReadWriteOnce"]` | How the volume may be mounted: `ReadWriteOnce` (one node — every block driver), `ReadWriteMany` (shared filesystems like EFS/Filestore/Azure Files), `ReadOnlyMany`, `ReadWriteOncePod` (exactly one pod). |
| `spec.storage_class_name` | `StringValueOrRef` | cluster default | The StorageClass provisioning this claim — a literal class name or a reference to a `KubernetesStorageClass` resource. When omitted, the cluster's DEFAULT class applies. |
| `spec.disable_dynamic_provisioning` | `bool` | `false` | Pins `storageClassName` to `""` — the claim binds only to a pre-provisioned PersistentVolume, never triggering provisioning. Mutually exclusive with `storage_class_name` (an empty class name and an absent one mean different things to Kubernetes). |
| `spec.storage_limit` | `string` | none | Upper bound on the volume's size. Rarely needed — only meaningful to drivers that honor limits. |
| `spec.volume_mode` | `filesystem \| block` | `filesystem` | `filesystem` is what nearly every workload wants; `block` exposes a raw device for workloads that manage their own on-disk format (requires driver support). |
| `spec.volume_name` | `string` | none | Binds to one specific pre-provisioned PersistentVolume by name — the adopt-existing-data path. |
| `spec.selector` | `LabelSelector` | none | Narrows which pre-provisioned volumes may bind. Only meaningful for static binding. |
| `spec.data_source` | `{kind, name}` | none | Populates the volume from a source instead of empty: `kind: persistent_volume_claim` (clone) or `kind: volume_snapshot` (restore), same-namespace only. Requires a capable CSI driver. **Pulumi engine only** — the Terraform module rejects it at plan time. |
| `spec.labels` / `spec.annotations` | `map<string, string>` | `{}` | Merged with standard Planton labels / applied to the object. |

## Examples

### Pinned to a Platform Class

The production shape — the claim names its class instead of inheriting the cluster default:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPersistentVolumeClaim
metadata:
  name: postgres-data
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPersistentVolumeClaim.postgres-data
spec:
  namespace:
    value: backend
  name: postgres-data
  access_modes:
    - ReadWriteOnce
  storage_request: 50Gi
  storage_class_name:
    value: fast-ssd
```

### Shared Volume Across Pods

ReadWriteMany needs a shared-filesystem driver (EFS, Filestore, Azure Files) — block-storage classes will never provision it:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPersistentVolumeClaim
metadata:
  name: shared-assets
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPersistentVolumeClaim.shared-assets
spec:
  namespace:
    value: backend
  name: shared-assets
  access_modes:
    - ReadWriteMany
  storage_request: 100Gi
  storage_class_name:
    value: efs-sc
```

### Adopt a Pre-Provisioned Volume

Dynamic provisioning disabled; the claim binds only to a matching pre-provisioned PersistentVolume, narrowed by a selector:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPersistentVolumeClaim
metadata:
  name: archive-data
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPersistentVolumeClaim.archive-data
spec:
  namespace:
    value: backend
  name: archive-data
  access_modes:
    - ReadWriteOnce
  storage_request: 500Gi
  disable_dynamic_provisioning: true
  selector:
    match_labels:
      tier: archive
```

### Clone an Existing Claim

The new volume starts as a copy of `postgres-data`. Data sources deploy via the Pulumi engine only:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesPersistentVolumeClaim
metadata:
  name: postgres-data-copy
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesPersistentVolumeClaim.postgres-data-copy
spec:
  namespace:
    value: backend
  name: postgres-data-copy
  access_modes:
    - ReadWriteOnce
  storage_request: 50Gi
  storage_class_name:
    value: fast-ssd
  data_source:
    kind: persistent_volume_claim
    name: postgres-data
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `pvcName` | `string` | Name of the PersistentVolumeClaim object as created in the cluster — the value workload volume mounts reference |
| `namespace` | `string` | Namespace the claim was created in |
| `storageRequest` | `string` | The requested storage size, as a Kubernetes quantity |

The outputs deliberately avoid bind-time status (bound volume name, phase): a claim under a `wait_for_first_consumer` class is correctly Pending until a pod consumes it.

## Related Components

- [KubernetesStorageClass](/docs/catalog/kubernetes/kubernetesstorageclass) — defines the storage tier; reference it from `spec.storage_class_name` so both deploy in one run
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) — provides the target namespace; reference it from `spec.namespace`
- [KubernetesDeployment](/docs/catalog/kubernetes/kubernetesdeployment) — mounts the claim by name in its volume mounts
- [KubernetesStatefulSet](/docs/catalog/kubernetes/kubernetesstatefulset) — for per-replica storage, use its `volume_claim_templates` instead of standalone claims
