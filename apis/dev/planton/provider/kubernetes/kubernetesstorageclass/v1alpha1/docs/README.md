# Kubernetes Storage Class: Research Documentation

## Introduction

Persistent storage in Kubernetes is a two-sided contract: workloads request storage through PersistentVolumeClaims, and the cluster satisfies those requests by provisioning PersistentVolumes. The StorageClass is the piece in between — the cluster's storage menu. A `storage.k8s.io/v1` StorageClass names a provisioner (the CSI driver that creates volumes), the provisioner-specific parameters (disk type, IOPS, encryption), and the lifecycle policies for the volumes it provisions: what happens on claim deletion (reclaim), when binding happens (binding mode), and whether volumes can grow (expansion).

Three structural facts define the resource:

- **StorageClasses are cluster-scoped.** One class serves claims from every namespace; there is no namespace field.
- **The parameters vocabulary belongs to the provisioner, not to Kubernetes.** The API stores `parameters` as an opaque string map and never validates it — the CSI driver consumes it at provision time. A typo surfaces as a provisioning failure on the first claim, never at class creation.
- **The identity fields are immutable.** `provisioner` and `parameters` (and upstream, `reclaimPolicy` and `volumeBindingMode`) cannot be changed in place; a change means replacing the class. Existing volumes keep the settings they were provisioned with.

Planton's **KubernetesStorageClass** component brings the full `storage.k8s.io/v1` surface to the platform with schema-level validation, a typed default-class field, and dual-IaC support.

## Evolution and Historical Context

### Origins (Kubernetes 1.2–1.6)

Dynamic provisioning entered Kubernetes as an alpha annotation-driven mechanism in 1.2 (2016); StorageClass graduated to `storage.k8s.io/v1` in 1.6 (2017). The original model paired each class with an in-tree plugin (`kubernetes.io/aws-ebs`, `kubernetes.io/gce-pd`) — provisioning code compiled into Kubernetes itself.

### The CSI migration

The Container Storage Interface moved provisioning out of the Kubernetes core and into vendor-owned drivers. CSI reached GA in 1.13 (2018), and the in-tree cloud plugins were progressively migrated and then removed — the in-tree AWS EBS, GCE PD, and Azure Disk plugins are gone from current Kubernetes, and their StorageClass provisioner strings now name CSI drivers (`ebs.csi.aws.com`, `pd.csi.storage.gke.io`, `disk.csi.azure.com`). The practical consequence: the provisioner string is a driver-installation dependency, and each driver publishes its own parameters vocabulary.

### WaitForFirstConsumer (1.12)

`volumeBindingMode: WaitForFirstConsumer` reached GA in 1.12, solving the zonal-storage scheduling problem: with the original `Immediate` binding, a volume provisions the moment the claim is created — before any pod exists — and on a multi-zone cluster it can land in a zone where the eventual pod cannot schedule. WaitForFirstConsumer inverts the order: the scheduler places the pod first, then the volume provisions in that pod's topology. The visible consequence is that **a claim under such a class stays Pending until a pod consumes it — correct behavior that is routinely misread as a failure.** `allowedTopologies` arrived with the same machinery, and is only honored under WaitForFirstConsumer binding: immediate binding cannot see pod topology, so the combination is meaningless (and rejected for provisioners that don't support it).

### Expansion (1.24) and the default-class annotation

Volume expansion (`allowVolumeExpansion`) reached GA in 1.24 — grow-only, driver-dependent, and off by default. The default StorageClass mechanism, by contrast, never became a field at all: a cluster's default class is marked by the `storageclass.kubernetes.io/is-default-class: "true"` annotation, a convention the PVC admission controller reads when a claim names no class. Kubernetes tolerates multiple defaults by picking the newest, but that behavior has shifted across versions and the documented contract remains: keep exactly one.

### What StorageClass never became

Upstream deliberately kept the class a thin, opaque pointer: no typed per-driver parameters, no validation of mount options (an invalid option fails at pod mount, not at class creation), no migration of existing volumes when a class changes. Per-volume attribute mutation went to a separate resource entirely (VolumeAttributesClass, still graduating). The class is a naming and defaulting mechanism, not a controller.

## The Semantics in Detail

### Reclaim policy

`reclaimPolicy` is stamped onto volumes at provision time. `Delete` (the default) removes the backing volume with the claim — right for disposable, recreatable data. `Retain` keeps the volume (and its cost) after claim deletion; it must be manually reclaimed. The policy on the class only affects volumes provisioned *after* the class exists — it is a template, not retroactive.

### Binding mode and the Pending state

`Immediate` binds and provisions on claim creation. `WaitForFirstConsumer` defers until a pod uses the claim. Under WaitForFirstConsumer, `kubectl get pvc` shows `Pending` with the event `waiting for first consumer to be created before binding` — the designed steady state for an unconsumed claim, not an error condition. Anything that waits for a claim to reach Bound before a consumer exists will wait forever.

### Expansion

`allowVolumeExpansion: true` lets a claim's `storage_request` be raised after creation; Kubernetes never shrinks a volume. The CSI driver must implement expansion (the major cloud drivers do). Enabling it on a class costs nothing on capable drivers and converts "database outgrew its disk" from a migration into an edit.

### The immutability boundary

`provisioner` and `parameters` are immutable upstream — the API rejects in-place changes. Since StorageClass names are cluster-unique, replacing a class must delete the old object before creating the new one. Volumes already provisioned keep their original parameters; only future provisioning changes.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

There is no `kubectl create storageclass` generator. Manual means `kubectl apply -f` of hand-written YAML.

**Verdict:** No shortcut exists; even ad-hoc use is YAML authoring.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
  annotations:
    storageclass.kubernetes.io/is-default-class: "false"
provisioner: ebs.csi.aws.com
parameters:
  type: gp3
  encrypted: "true"
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
allowVolumeExpansion: true
```

**Pros:**
- Declarative, version-controllable, full surface

**Cons:**
- Parameters are unvalidated opaque strings; typos surface as provisioning failures on the first claim
- The default-class marker is a magic annotation string, easy to typo and easy to double-set
- Immutability surprises: editing `provisioner` or `parameters` is rejected at apply with no offer to replace

**Verdict:** The baseline; thin guardrails around an opaque contract.

### Level 2: Terraform

```hcl
resource "kubernetes_storage_class_v1" "fast_ssd" {
  metadata {
    name = "fast-ssd"
  }
  storage_provisioner    = "ebs.csi.aws.com"
  parameters             = { type = "gp3", encrypted = "true" }
  reclaim_policy         = "Delete"
  volume_binding_mode    = "WaitForFirstConsumer"
  allow_volume_expansion = true
}
```

**Pros:**
- Full IaC lifecycle (plan, apply, destroy, import) with drift detection; the provider forces replacement on immutable-field changes

**Cons:**
- The provider models `allowed_topologies` as a single selector term — the API's multiple OR'd terms are inexpressible in HCL
- Parameters remain unvalidated pass-through

**Verdict:** Production-grade lifecycle; one expressiveness gap on topology.

### Level 3: Pulumi

```go
storageClass, err := storagev1.NewStorageClass(ctx, "fast-ssd", &storagev1.StorageClassArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name: pulumi.String("fast-ssd"),
    },
    Provisioner:          pulumi.String("ebs.csi.aws.com"),
    Parameters:           pulumi.StringMap{"type": pulumi.String("gp3")},
    VolumeBindingMode:    pulumi.String("WaitForFirstConsumer"),
    AllowVolumeExpansion: pulumi.Bool(true),
})
```

**Pros:**
- Full programming language, preview before apply, complete surface including multiple topology terms

**Cons:**
- Types describe the wire shape, not the semantics; a dueling default class or a typo'd parameter passes the compiler

**Verdict:** Excellent IaC choice; validation gap same as Terraform.

### Other Methods

**Helm:** classes templated inside infrastructure charts — common for platform add-ons that ship their own storage tier, but a cluster-scoped singleton templated per-release invites collisions.

**Managed-cluster defaults:** every managed Kubernetes offering ships classes out of the box (EKS: gp2/gp3; GKE: standard-rwo/premium-rwo; AKS: default/managed-csi). Fine until the default's reclaim policy, binding mode, or expansion setting doesn't match requirements — which is exactly when a custom class enters.

## Comparative Analysis

| Aspect | YAML | Terraform | Pulumi | Planton |
|--------|------|-----------|--------|---------|
| Validation | API server | Plan time (shape only) | Preview time (shape only) | Schema + CEL |
| Name/enum contracts checked early | No | Partial | No | Yes |
| Default-class as typed field | No (annotation string) | No (annotation string) | No (annotation string) | Yes, `is_default_class` |
| Deterministic policies across engines | N/A | Provider-dependent | Provider-dependent | Always explicit |
| Multiple topology terms | Yes | No (provider limit) | Yes | Pulumi yes; TF fails plan loudly |
| Dual IaC | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### Full surface, validated early

The spec models the entire `storage.k8s.io/v1` StorageClass — provisioner, parameters, reclaim policy, binding mode, expansion, mount options, topology terms, default-class marker — and moves the checkable contracts to validation time:

- **Name contract**: DNS-subdomain shape (lowercase alphanumeric, hyphens, dots, 253-char bound) enforced by CEL before anything reaches a cluster
- **Enum-typed policies**: `reclaim_policy` and `volume_binding_mode` are closed enums, not free strings — `Recycle` (deprecated upstream) and typos are unrepresentable
- **Topology term shape**: every term requires at least one requirement; every requirement requires a key and at least one value — the empty-term mistake is rejected at the schema

### Parameters stay upstream-shaped

`parameters` is deliberately an opaque string map, exactly as upstream: each CSI driver defines its own keys, and typed modeling would couple the class to one driver. Values that name secrets (e.g. `csi.storage.k8s.io/provisioner-secret-name`) are references to Secret objects, never secret material itself.

### The default class is a field, not a magic string

`is_default_class: true` renders the `storageclass.kubernetes.io/is-default-class: "true"` annotation — the upstream mechanism, surfaced as a typed boolean. The spec's annotation map is documented as NOT the place for the marker, so the two paths cannot silently disagree.

### Deterministic policies

Both IaC modules resolve `reclaim_policy` and `volume_binding_mode` to their API strings with the API server's own defaults (`Delete`, `Immediate`) applied module-side, and always submit them explicitly. `allow_volume_expansion` is likewise always sent. Both engines submit byte-identical objects for the same manifest.

### Replacement, not mutation

`provisioner` and `parameters` are immutable upstream. The Pulumi module sets delete-before-replace so a forced replacement never collides with the cluster-unique class name; the Terraform provider forces replacement on the same fields. Neither engine attempts an in-place update the API would reject.

### One parity exception, failed loudly

The Terraform Kubernetes provider models `allowed_topologies` as a **single** selector term (`max_items = 1`), while the API — and the Pulumi module — accept multiple OR'd terms. The Terraform module passes one term through intact and **fails the plan with a precondition** when the spec lists several, with an error message that names the two escapes: deploy with the Pulumi provisioner, or combine the zone values into a single term (values within one requirement already OR together). Failing loudly at plan beats silently dropping terms.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, resource creation, and output export
- **`locals.go`**: Computes merged labels, annotations (including the default-class marker), and the resolved reclaim policy and binding mode
- **`storageclass.go`**: Creates the `storage.k8s.io/v1` StorageClass with delete-before-replace semantics
- **`outputs.go`**: Exports `storage_class_name`, `provisioner`, and `is_default_class`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes merged labels, the default-class annotation, and the same resolved policy strings
- **`main.tf`**: Creates the `kubernetes_storage_class_v1` resource, with the single-topology-term precondition
- **`outputs.tf`**: Exports the same three outputs

### Resource Count

This is a lean component — it creates exactly **one Kubernetes resource**: the StorageClass itself. The complexity is in the policy semantics and the immutability handling, not in resource orchestration.

## Production Best Practices

### Class design

1. **`wait_for_first_consumer` for anything zonal**: EBS, GCE PD, and Azure Disk volumes are zone-bound; immediate binding provisions before the scheduler runs and can strand the volume where no pod fits
2. **Enable expansion up front**: grow-only resizing costs nothing to allow on capable drivers and converts a disk-full incident into a claim edit
3. **Name classes for the tier, not the driver**: `fast-ssd`, `archive-hdd` — provisioner and parameters are immutable, and a class named `gp3-v2` ages badly

### Default-class discipline

1. **Exactly one default per cluster**: demote the existing default (including a managed cluster's built-in one) before promoting a replacement — multiple defaults make newest-wins behavior version-dependent
2. **Pin classes explicitly in claims that matter**: the default is a convenience for ad-hoc claims; production claims should name their class

### Operational awareness

1. **A Pending claim under WaitForFirstConsumer is healthy**: the event `waiting for first consumer to be created before binding` is the designed steady state, not a failure to investigate
2. **Parameters fail late by design**: the class object creates cleanly on any cluster; the first claim exercising a typo'd parameter is where the error surfaces. Test a class with a real claim and consumer before declaring it done
3. **Retained volumes accrue cost**: `retain` moves cleanup from the platform to a human; audit released volumes periodically

## Conclusion

KubernetesStorageClass is a deliberately complete, deliberately lean component: the full upstream surface, with the checkable contracts (names, enums, topology shapes) enforced at the schema, the default-class annotation lifted to a typed field, both engines submitting identical objects with the API's own defaults made explicit, and the one provider expressiveness gap — multiple topology terms on Terraform — converted from silent truncation into a loud plan failure. Combined with claim references through `storage_class_name`, it makes a platform-defined storage tier something claims compose with, not something they hope the cluster default happens to be.

## References

- [Kubernetes Storage Classes Documentation](https://kubernetes.io/docs/concepts/storage/storage-classes/)
- [StorageClass API Reference](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/storage-class-v1/)
- [Dynamic Volume Provisioning](https://kubernetes.io/docs/concepts/storage/dynamic-provisioning/)
- [Volume Binding Mode & Topology](https://kubernetes.io/docs/concepts/storage/storage-classes/#volume-binding-mode)
- [Default StorageClass Behavior](https://kubernetes.io/docs/tasks/administer-cluster/change-default-storage-class/)
- [AWS EBS CSI Driver Parameters](https://github.com/kubernetes-sigs/aws-ebs-csi-driver/blob/master/docs/parameters.md)
- [GCE PD CSI Driver](https://github.com/kubernetes-sigs/gcp-compute-persistent-disk-csi-driver)
- [Azure Disk CSI Driver Parameters](https://github.com/kubernetes-sigs/azuredisk-csi-driver/blob/master/docs/driver-parameters.md)
- [Pulumi Kubernetes StorageClass](https://www.pulumi.com/registry/packages/kubernetes/api-docs/storage/v1/storageclass/)
- [Terraform kubernetes_storage_class_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/storage_class_v1)
