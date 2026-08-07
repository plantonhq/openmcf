# KubernetesKubeRayOperator: Research and Design

## Introduction

KubernetesKubeRayOperator installs the KubeRay operator from the
official `kuberay-operator` Helm chart
(https://ray-project.github.io/kuberay-helm, pinned 1.6.2) as a
single Helm release named after `metadata.name`. The operator is the
ENGINE of the Ray story in this catalog: KubernetesRayCluster
declares `RayCluster` custom resources, and this operator reconciles
them into running Ray clusters. RayJob and RayService CRs are
unmodeled here — author them directly against the same operator when
needed.

## The Deployment Landscape

One operator serves the whole cluster under the default posture: it
watches every namespace and runs leader election out of the box. Ray
clusters themselves are the many-per-cluster workload — the catalog
splits the concern in two, this kind installing the engine once and
KubernetesRayCluster declaring each cluster.

### The chart/operator pinning truth

Chart 1.6.2 pairs with operator image `quay.io/kuberay/operator:v1.6.2`.
Versions must exist as SERVED charts in the repository index.

## The CRD Lifecycle: crds/ Directory, Keep by Construction

The chart ships its three ray.io CRDs (`rayclusters`, `rayjobs`,
`rayservices`) from its `crds/` DIRECTORY. Helm's contract for that
directory: install once, never upgrade on chart upgrades, leave on
uninstall (no release ownership metadata).

That upstream posture is exactly the keep-on-uninstall this catalog
wants for workload-bearing CRDs — removing the operator never deletes
a Ray declaration — so the modules neither re-own nor template the
CRDs. Chart-version bumps never touch them: apply the new release's
CRD files manually when a bump changes them, per the upstream release
notes. The CRDs are large (~1MB each) and install server-side.

## No Admission Webhook

The operator validates in its reconcile loop. There is no admission
webhook, no certificate machinery, and no cert-manager dependency. A
bad RayCluster surfaces on the CR's status conditions, not as an
admission rejection.

## The Feature-Gate Merge: Helm Lists Replace

The chart's `featureGates` value is a LIST, and Helm lists REPLACE,
never merge: rendering only the spec's entries would silently DROP
every chart-default gate. So when the spec flips any gate, the
modules render the FULL list — the chart's defaults at the pinned
1.6.2 (verified against its values.yaml):

| Gate | Chart default |
|---|---|
| `RayClusterStatusConditions` | on |
| `RayJobDeletionPolicy` | on |
| `RayMultiHostIndexing` | on |
| `RayServiceIncrementalUpgrade` | off |
| `RayCronJob` | off |

— overridden by name from the spec, then spec gates the defaults
don't know appended. The default list is kept in lockstep between the
two engines and re-verified on every chart bump. The consequence for
users: only list gates you are deliberately flipping; unlisted gates
keep their chart-default state.

## The Watch Scope

The chart's `watchNamespace` key (singular key, list value) feeds the
operator's `--watch-namespace` flag verbatim, and the chart scopes
its per-namespace reconcile RBAC to the same list. Empty (the chart
default) = cluster-wide, the normal one-operator-per-cluster posture —
the key stays unrendered. Non-empty is the fenced multi-tenant
posture: Ray declarations OUTSIDE the listed namespaces are ignored
without an error. Unlike the Spark operator's workload surface, the
watch namespaces here are ONLY a watch scope — they must already
exist and are deliberately not module resources.

## Two Image Seams

- **The operator image** is `quay.io/kuberay/operator` (the chart
  default). `image_registry` replaces ONLY the registry part (the
  swap drops quay.io; the path stays `kuberay/operator`; the tag
  stays the chart's appVersion-locked default) — the air-gap mirror
  seam for the one image THIS component's pods pull.
- **Ray cluster images** are a different seam entirely: each
  KubernetesRayCluster declares its own image, and mirroring the
  operator does nothing for them. Mirror the Ray image on each
  cluster declaration.

## Design Decisions

- **The name pins.** The chart hardcodes `nameOverride`,
  `fullnameOverride`, AND `serviceAccount.name` to `kuberay-operator`
  in its values — every install would collapse onto the same child
  names and the same ServiceAccount by construction. The modules pin
  all three to `metadata.name`, the catalog's Helm-kind identity
  convention, so instances stay distinguishable. The 47-character
  name budget follows: the chart's longest derived suffix is
  `-leader-election` (16 chars, the leader-election Role/RoleBinding)
  against the Kubernetes 63-character limit; both modules fail loudly
  over the budget.
- **`batchScheduler.name` is the only rendered scheduler knob.** The
  chart's `batchScheduler.enabled` is the deprecated legacy flag and
  MUTUALLY EXCLUSIVE with `name` (the chart errors when both are
  set) — the modules never render it. The named scheduler (volcano,
  yunikorn, or scheduler-plugins) must already run on the cluster;
  the operator only emits its scheduling directives.
- **The install is blocking.** The Helm release waits for the operator
  Deployment to become Available (atomic, 600s timeout, cleanup on
  fail): an unpullable image, a missing ServiceMonitor CRD, or a
  broken config fails THIS apply with a readiness timeout instead of
  surfacing later as RayClusters that mysteriously never reconcile.
- **The ServiceMonitor fails loudly.** `service_monitor_enabled`
  requires the monitoring.coreos.com CRDs on the cluster
  (KubernetesKubePrometheusStack first) — the install FAILS without
  them, by upstream design rather than a silent skip.
- **The module owns namespace creation** (`create_namespace`), never
  the Helm release — pre-existing-namespace installs leave the flag
  false.
- **Chart-default-matching values render only on divergence** — the
  chart's real defaults stand (100m CPU / 512Mi limits; upstream
  sizes ~500MB per 500 managed Ray pods), `leaderElectionEnabled` and
  `metrics.enabled` render only on an explicit false, and the name
  pins are the one always-rendered set.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `kuberay-operator` at https://ray-project.github.io/kuberay-helm | Pinned 1.6.2 (spec default) |
| Operator image | `quay.io/kuberay/operator` | Chart 1.6.2 = operator v1.6.2 |
| CRD API group | `ray.io` | Three CRDs (`rayclusters`, `rayjobs`, `rayservices`), chart `crds/` directory — installed once, kept on uninstall, never upgraded by chart bumps; ~1MB each, server-side apply |
| Name pins | `nameOverride` / `fullnameOverride` / `serviceAccount.name` = `metadata.name` | The chart hardcodes all three to `kuberay-operator`; the 16-char `-leader-election` suffix behind the 47-char name budget |
| Feature-gate defaults (1.6.2) | 3 on, 2 off (table above) | Re-verify on every chart bump; both engines carry the same list |

## IaC Twins

Pulumi (`module/values.go` + `module/locals.go`) and Terraform
(`main.tf` + `locals.tf`) render identical chart values: the same
name pins, the same feature-gate merge over the chart's default list,
the same watch-scope wiring, and the same blocking-install posture.
Keep the typed-value rendering and the feature-gate default list in
lockstep.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
