# Kubernetes Persistent Volume Claim: Research Documentation

## Introduction

Pods are ephemeral; their data usually shouldn't be. Kubernetes separates the request for durable storage from its fulfillment: a **PersistentVolumeClaim** is a namespaced `core/v1` resource declaring how much storage a consumer needs, how it will be accessed, and (optionally) which StorageClass provisions it. The cluster satisfies the claim by binding it to a **PersistentVolume** — either dynamically provisioned by a CSI driver on demand, or a pre-provisioned volume matched by capacity, access modes, class, and (optionally) selector.

The claim's mental model has two load-bearing subtleties worth stating precisely:

- **The empty-vs-absent class distinction.** A claim whose `storageClassName` is ABSENT gets the cluster's default StorageClass. A claim whose `storageClassName` is the EMPTY STRING opts out of dynamic provisioning entirely and binds only to pre-provisioned volumes. These are different requests, and a tool that flattens them into one string silently converts one into the other.
- **Pending is a designed state, not a failure.** Under a `WaitForFirstConsumer` StorageClass (the norm for zonal cloud disks), a claim stays Pending until a pod uses it — the volume then provisions in the topology that pod scheduled into. Anything that waits for Bound before a consumer exists waits forever.

Planton's **KubernetesPersistentVolumeClaim** component brings the full `core/v1` surface to the platform with schema-level validation, typed class selection, composition through references, and dual-IaC support.

## Evolution and Historical Context

### Origins (Kubernetes 1.0–1.6)

PersistentVolumes and PersistentVolumeClaims shipped with Kubernetes 1.0 (2015) as the two-sided storage contract: administrators pre-provisioned volumes, users claimed them. Dynamic provisioning arrived as an annotation-driven alpha in 1.2 and became a first-class mechanism with StorageClass's graduation in 1.6 (2017) — the claim's `storageClassName` field replacing the `volume.beta.kubernetes.io/storage-class` annotation. The empty-vs-absent distinction dates from this transition: an empty class name preserved the pre-dynamic behavior (static binding only), while absence deferred to the new default-class machinery.

### The CSI era

The Container Storage Interface (GA in 1.13) moved provisioning into vendor-owned drivers, and with it came the claim-side features that depend on driver capability: volume expansion (GA 1.24 — grow-only, edit the claim's request), raw block volumes (`volumeMode: Block`, GA 1.18), volume cloning (`dataSource` naming a PVC, GA 1.18), and snapshot restore (`dataSource` naming a VolumeSnapshot, snapshots GA in 1.20). Each is a claim-spec field whose actual behavior is delegated to the CSI driver — the API accepts the request; the driver decides whether it can be honored.

### WaitForFirstConsumer changes what "created" means

StorageClass's `WaitForFirstConsumer` binding mode (GA 1.12) solved zonal topology mismatches by deferring provisioning until a pod consumes the claim. It also quietly changed the claim's lifecycle contract: an unconsumed claim under such a class is *permanently, correctly Pending*. Tooling written for the Immediate-binding era — including IaC providers that default to waiting for Bound — hangs on every such claim.

### What remains gated

`volumeAttributesClassName` (mutating provisioned-volume attributes like IOPS via a VolumeAttributesClass) and cross-namespace data sources (`dataSourceRef` with a namespace) are still feature-gated upstream and are deliberately unmodeled in this spec until they graduate.

## The Semantics in Detail

### Access modes are requests, not guarantees

`accessModes` declares how the volume may be mounted: `ReadWriteOnce` (one node at a time — the mode every block driver supports), `ReadOnlyMany`, `ReadWriteMany` (requires a shared-filesystem driver: EFS, Filestore, Azure Files, CephFS, NFS), and `ReadWriteOncePod` (exactly one pod, the strictest isolation). A claim requesting `ReadWriteMany` against a block-storage class does not error at creation — it simply never provisions. The API requires at least one mode; there is no server default.

### The class-selection triangle

Three shapes, mutually exclusive by meaning:

1. **Absent class** → the cluster's default StorageClass (the `storageclass.kubernetes.io/is-default-class` annotation) applies at admission
2. **Named class** → dynamic provisioning through that class (or static binding to volumes carrying that class)
3. **Empty-string class** → no dynamic provisioning, ever; bind only to a matching pre-provisioned volume

### Static binding

`volumeName` binds the claim to one specific PersistentVolume by name; a `selector` narrows candidate volumes by label. Both are static-binding tools — dynamically provisioned volumes are created *for* the claim, not selected. The volume's capacity, access modes, and class must still satisfy the claim.

### Data sources

`dataSource` populates the new volume from an existing object instead of provisioning it empty: a same-namespace PersistentVolumeClaim (clone) or VolumeSnapshot (restore, `snapshot.storage.k8s.io` group). Both require a CSI driver implementing the operation, and the clone/restore size must be at least the source's.

### Expansion, never contraction

Raising a bound claim's `storage_request` triggers expansion when the class allows it and the driver supports it. Kubernetes never shrinks a volume; a lower request is rejected.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

There is no `kubectl create pvc` generator. Manual means `kubectl apply -f` of hand-written YAML.

**Verdict:** No shortcut exists; even ad-hoc use is YAML authoring.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: app-data
  namespace: backend
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: fast-ssd
```

**Pros:**
- Declarative, version-controllable, full surface

**Cons:**
- The empty-vs-absent class distinction is invisible in review (`storageClassName: ""` vs the line missing)
- Quantity typos (`10gi`, `10G B`) and access-mode typos surface only at admission
- No plan/preview, no state management

**Verdict:** The baseline; the sharp edges are quiet ones.

### Level 2: Terraform

```hcl
resource "kubernetes_persistent_volume_claim_v1" "app_data" {
  metadata {
    name      = "app-data"
    namespace = "backend"
  }
  spec {
    access_modes = ["ReadWriteOnce"]
    resources {
      requests = { storage = "10Gi" }
    }
    storage_class_name = "fast-ssd"
  }
  wait_until_bound = false
}
```

**Pros:**
- Full IaC lifecycle (plan, apply, destroy, import) with drift detection

**Cons:**
- `wait_until_bound` defaults to TRUE — every claim under a WaitForFirstConsumer class hangs the apply until timeout unless the author knows to disable it
- The provider's PVC resource cannot express `dataSource` at all — clones and snapshot restores are inexpressible in HCL

**Verdict:** Production-grade lifecycle with a hang-by-default footgun and a genuine expressiveness gap.

### Level 3: Pulumi

```go
pvc, err := corev1.NewPersistentVolumeClaim(ctx, "app-data", &corev1.PersistentVolumeClaimArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name:      pulumi.String("app-data"),
        Namespace: pulumi.String("backend"),
    },
    Spec: &corev1.PersistentVolumeClaimSpecArgs{
        AccessModes: pulumi.ToStringArray([]string{"ReadWriteOnce"}),
        Resources: &corev1.VolumeResourceRequirementsArgs{
            Requests: pulumi.ToStringMap(map[string]string{"storage": "10Gi"}),
        },
        StorageClassName: pulumi.String("fast-ssd"),
    },
})
```

**Pros:**
- Full programming language, preview before apply, complete surface including data sources

**Cons:**
- Pulumi's readiness logic awaits PVC Bound by default — the same WaitForFirstConsumer hang, opted out via the `pulumi.com/skipAwait` annotation
- Types describe the wire shape, not the semantics

**Verdict:** Excellent IaC choice; the await default needs the same deliberate opt-out.

### Other Methods

**Helm:** claims templated inside application charts — common, but chart-owned claims couple data lifecycle to release lifecycle (helm uninstall can delete the claim and, with a `Delete` reclaim class, the data).

**StatefulSet volumeClaimTemplates:** the right mechanism for per-replica storage — the controller stamps one claim per replica and manages their identity. Standalone claims are for storage whose lifecycle is independent of any one workload.

## Comparative Analysis

| Aspect | YAML | Terraform | Pulumi | Planton |
|--------|------|-----------|--------|---------|
| Validation | API server | Plan time (shape only) | Preview time (shape only) | Schema + CEL |
| Quantity/access-mode contracts checked early | No | No | No | Yes |
| Empty-vs-absent class distinction | Invisible | Flattened to one string | Expressible, easy to miss | Typed field, conflict rejected |
| WaitForFirstConsumer deploy hang | N/A | Hangs by default | Hangs by default | Never waits, both engines |
| Data source (clone/snapshot) | Yes | **Inexpressible** | Yes | Pulumi yes; TF fails plan loudly |
| Namespace/class as reference | No | Manual wiring | Manual wiring | First-class |
| Dual IaC | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### Full surface, validated early

The spec models the entire `core/v1` PersistentVolumeClaimSpec — access modes, requests and limits, class selection, volume mode, static binding, selectors, data sources — and moves the checkable contracts to validation time:

- **Quantity contracts**: `storage_request` (required) and `storage_limit` must be Kubernetes quantities (`10Gi`, `500Mi`) — checked by CEL before anything reaches a cluster
- **Access-mode vocabulary**: only the four upstream modes are representable
- **Selector operator contracts**: `In`/`NotIn` require values, `Exists`/`DoesNotExist` forbid them — the exact admission rule, surfaced before deployment
- **Data-source kind**: a closed enum (`persistent_volume_claim`, `volume_snapshot`) with the unspecified value rejected

### The empty-vs-absent distinction is a typed field

`disable_dynamic_provisioning` pins `storageClassName` to the empty string (bind only pre-provisioned volumes); omitting `storage_class_name` leaves the field absent (cluster default applies). The two cannot be confused, and a CEL rule rejects setting `disable_dynamic_provisioning` alongside a named class — the contradiction is unrepresentable.

### Deploys never wait for Bound

Both engines opt out of their providers' bind-await behavior — the Pulumi module sets the `pulumi.com/skipAwait` annotation; the Terraform module sets `wait_until_bound = false`. A claim under a WaitForFirstConsumer class is correctly Pending until a pod consumes it, and awaiting would hang every such deploy. For the same reason, the stack outputs deliberately avoid bind-time status (bound volume name, phase): the claim's existence and its requested size are the deploy-time truths; binding is the consumer's event.

### Defaults applied module-side, sent explicitly

The Kubernetes-default access mode (`ReadWriteOnce`) is applied in the modules when the spec omits `access_modes` — the API itself REQUIRES the field, so there is no server default to defer to. `volume_mode` resolves to its API string (`Filesystem` default) and is always sent. Both engines submit identical claims for the same manifest.

### Namespace and class by value or reference

Both `spec.namespace` and `spec.storage_class_name` are `StringValueOrRef` fields: literal names, or references to `KubernetesNamespace` / `KubernetesStorageClass` resources — the class reference resolving through the class's exported `storage_class_name` output. An infra chart can create the namespace, the class, and the claim in one run with ordering handled by the resource graph.

### One parity exception, failed loudly

The Terraform Kubernetes provider's PVC resource cannot express `spec.dataSource`/`dataSourceRef` — clones and snapshot restores are inexpressible in its schema. The Pulumi module sends the data source natively; the Terraform module **fails the plan with a precondition** when `data_source` is set, with an error message naming the escapes: deploy the claim with the Pulumi provisioner, or drop the field. Failing loudly beats silently provisioning an EMPTY volume where the user asked for restored data.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates provider init, resource creation, and output export
- **`locals.go`**: Computes merged labels, annotations (including skipAwait), the resolved namespace, access modes, volume mode, and the three-valued storageClassName
- **`persistentvolumeclaim.go`**: Creates the `core/v1` PersistentVolumeClaim, mapping the data source onto the API's TypedLocalObjectReference (PVC clones in the core group; snapshot restores naming `snapshot.storage.k8s.io`)
- **`outputs.go`**: Exports `pvc_name`, `namespace`, and `storage_request`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes merged labels, resolved namespace, access modes, volume mode, and the same three-valued storageClassName
- **`main.tf`**: Creates the `kubernetes_persistent_volume_claim_v1` resource with `wait_until_bound = false` and the data-source precondition
- **`outputs.tf`**: Exports the same three outputs

### Resource Count

This is a lean component — it creates exactly **one Kubernetes resource**: the PersistentVolumeClaim itself. The complexity is in the class-selection semantics and the binding-timing discipline, not in resource orchestration.

## Production Best Practices

### Claim design

1. **Name the class explicitly in production**: the cluster default is a moving target; `storage_class_name` — ideally as a reference to a platform-defined class — pins the tier
2. **Request `ReadWriteMany` only from drivers that provide it**: block-storage classes never satisfy it, and the claim pends forever with no error at creation
3. **Size on an expandable class**: expansion is grow-only and class-gated; a full disk on a non-expandable class is a migration, not an edit

### Lifecycle discipline

1. **Pending under WaitForFirstConsumer is healthy**: the event `waiting for first consumer to be created before binding` is the design working — deploy pipelines must not treat it as failure (this component's engines already don't)
2. **Mind the reclaim policy behind the claim**: deleting a claim on a `Delete`-reclaim class deletes the data; put irreplaceable data on a `retain` class
3. **Keep StatefulSet storage in the StatefulSet**: `volume_claim_templates` gives per-replica claims with controller-managed identity; standalone claims are for workload-independent lifecycles

### Static binding and sources

1. **Static binding is deliberate work**: `disable_dynamic_provisioning`, `volume_name`, and `selector` are the adopt-existing-data tools — the claim pends until a matching PersistentVolume exists, which is the contract, not a bug
2. **Clones and restores are driver features**: verify the CSI driver implements them before building recovery runbooks on `data_source` — and remember these claims deploy via the Pulumi engine only

## Conclusion

KubernetesPersistentVolumeClaim is a deliberately complete, deliberately lean component: the full upstream surface, with the resource's quiet sharp edges — the empty-vs-absent class distinction, the request-not-guarantee access modes, the permanently-Pending-by-design binding mode — lifted into typed fields, schema validation, and engine behavior that never hangs a deploy waiting for a consumer. Combined with references to platform-defined StorageClasses and namespaces, it makes durable storage something an infra chart composes in one run rather than a sequence of hand-ordered applies.

## References

- [Kubernetes Persistent Volumes Documentation](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
- [PersistentVolumeClaim API Reference](https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/)
- [Storage Classes and Default Behavior](https://kubernetes.io/docs/concepts/storage/storage-classes/)
- [CSI Volume Cloning](https://kubernetes.io/docs/concepts/storage/volume-pvc-datasource/)
- [Volume Snapshots](https://kubernetes.io/docs/concepts/storage/volume-snapshots/)
- [Volume Expansion](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#expanding-persistent-volumes-claims)
- [Pulumi Kubernetes PersistentVolumeClaim](https://www.pulumi.com/registry/packages/kubernetes/api-docs/core/v1/persistentvolumeclaim/)
- [Terraform kubernetes_persistent_volume_claim_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/persistent_volume_claim_v1)
