---
title: "Storage Class"
description: "Storage Class deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesstorageclass"
---

# Kubernetes Storage Class

Deploys a Kubernetes StorageClass — the cluster's storage menu — to a target cluster through a single declarative manifest, covering the complete `storage.k8s.io/v1` surface: provisioner, provisioner-specific parameters, reclaim policy, volume binding mode, expansion, mount options, topology restrictions, and the default-class marker. The IaC module handles label merging, policy defaults, and the default-class annotation automatically.

## What Gets Created

When you deploy a KubernetesStorageClass resource, Planton provisions:

- **StorageClass** — a cluster-scoped `storage.k8s.io/v1` StorageClass naming the provisioner (CSI driver), its parameters, and the lifecycle policies for volumes it provisions
- **Labels** — standard Planton tracking labels merged with any user-provided labels
- **Annotations** — user-provided annotations, plus the `storageclass.kubernetes.io/is-default-class: "true"` annotation when `is_default_class` is set

**The parameters vocabulary belongs to the provisioner, not to Kubernetes.** Parameters pass through verbatim to the CSI driver and are consumed at provision time — a typo surfaces as a provisioning failure on the first claim, not at apply.

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **The named CSI driver installed** in the cluster — the class object itself is creatable regardless, but claims of the class only provision when the driver exists

## Quick Start

Create a file `storage-class.yaml` — an expandable gp3 class for EKS:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesStorageClass
metadata:
  name: fast-ssd
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesStorageClass.fast-ssd
spec:
  name: fast-ssd
  provisioner: ebs.csi.aws.com
  parameters:
    type: gp3
    encrypted: "true"
  volume_binding_mode: wait_for_first_consumer
  allow_volume_expansion: true
```

Deploy:

```shell
planton apply -f storage-class.yaml
```

Claims request the class by name (`storage_class_name: {value: fast-ssd}`). With `wait_for_first_consumer` binding, a claim stays Pending until a pod uses it — the volume then provisions in the zone that pod scheduled into. That Pending state is correct behavior, not an error.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.name` | `string` | Name of the StorageClass (`metadata.name` in the cluster) — the value claims put in `storage_class_name`. | 1–253 characters, matches `^[a-z0-9]([a-z0-9.-]{0,251}[a-z0-9])?$` |
| `spec.provisioner` | `string` | The CSI driver (or in-tree plugin) provisioning volumes — e.g. `ebs.csi.aws.com`, `pd.csi.storage.gke.io`, `disk.csi.azure.com`. **IMMUTABLE**: changing it replaces the class. | required |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.parameters` | `map<string, string>` | `{}` | Provisioner-specific parameters, passed through verbatim (e.g. EBS CSI: `{type: gp3, iops: "6000", encrypted: "true"}`). Each driver defines its own keys. **IMMUTABLE**. Values naming secrets are references, never secret material. |
| `spec.reclaim_policy` | `delete \| retain` | `delete` | What happens to a provisioned volume when its claim is deleted. `retain` keeps the volume for manual reclamation. |
| `spec.volume_binding_mode` | `immediate \| wait_for_first_consumer` | `immediate` | When claims bind and volumes provision. `wait_for_first_consumer` waits for a consuming pod and provisions in its topology — the right choice for zonal storage (EBS, GCE PD, Azure Disk). |
| `spec.allow_volume_expansion` | `bool` | `false` | Allows claims of this class to grow after creation (never shrink). Requires driver support. |
| `spec.mount_options` | `list(string)` | `[]` | Mount options for volumes of this class (e.g. `["noatime"]`). NOT validated by Kubernetes — invalid options fail at pod start. |
| `spec.allowed_topologies` | `list(term)` | `[]` | Restricts where volumes may provision (e.g. pin to zones). Requirements within a term AND; terms OR. Only honored with `wait_for_first_consumer` binding. |
| `spec.is_default_class` | `bool` | `false` | Marks this class as the cluster's default — rendered as the `storageclass.kubernetes.io/is-default-class: "true"` annotation. A cluster should have exactly ONE default; demote the old one first. |
| `spec.labels` / `spec.annotations` | `map<string, string>` | `{}` | Merged with standard Planton labels / applied to the object. Set `is_default_class` instead of adding the default-class annotation here. |

## Examples

### GCP SSD Class With Retain

Volumes survive claim deletion — the right shape for data that must outlive its workload:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesStorageClass
metadata:
  name: pd-ssd-retain
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesStorageClass.pd-ssd-retain
spec:
  name: pd-ssd-retain
  provisioner: pd.csi.storage.gke.io
  parameters:
    type: pd-ssd
  reclaim_policy: retain
  volume_binding_mode: wait_for_first_consumer
  allow_volume_expansion: true
```

### Zone-Pinned Class

Volumes may only provision in the listed zones. Topology restrictions require `wait_for_first_consumer`:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesStorageClass
metadata:
  name: gp3-us-east-1a
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesStorageClass.gp3-us-east-1a
spec:
  name: gp3-us-east-1a
  provisioner: ebs.csi.aws.com
  parameters:
    type: gp3
  volume_binding_mode: wait_for_first_consumer
  allow_volume_expansion: true
  allowed_topologies:
    - match_label_expressions:
        - key: topology.kubernetes.io/zone
          values:
            - us-east-1a
            - us-east-1b
```

### Promote to Cluster Default

Marks the class as the one claims receive when they name none. Demote the existing default first — two defaults make claim behavior undefined:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesStorageClass
metadata:
  name: standard-gp3
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesStorageClass.standard-gp3
spec:
  name: standard-gp3
  provisioner: ebs.csi.aws.com
  parameters:
    type: gp3
    encrypted: "true"
  volume_binding_mode: wait_for_first_consumer
  allow_volume_expansion: true
  is_default_class: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `storageClassName` | `string` | Name of the StorageClass object as created in the cluster — the value claims put in their `storage_class_name` |
| `provisioner` | `string` | The provisioner (CSI driver) backing this class |
| `isDefaultClass` | `bool` | Whether this class is annotated as the cluster's default StorageClass |

## Related Components

- [KubernetesPersistentVolumeClaim](/docs/catalog/kubernetes/persistent-volume-claim) — requests storage from this class; reference the class from `spec.storage_class_name` so both deploy in one run
- [KubernetesStatefulSet](/docs/catalog/kubernetes/statefulset) — per-replica storage via volume-claim templates that name this class
- [KubernetesDeployment](/docs/catalog/kubernetes/deployment) — mounts claims provisioned from this class
