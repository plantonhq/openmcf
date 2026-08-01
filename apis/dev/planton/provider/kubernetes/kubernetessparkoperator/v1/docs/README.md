# KubernetesSparkOperator: Research and Design

## Introduction

KubernetesSparkOperator installs the Apache Spark Kubernetes Operator
from the official ASF `spark-kubernetes-operator` Helm chart
(https://apache.github.io/spark-kubernetes-operator, pinned 1.8.0) as
a single Helm release named after `metadata.name`. The operator
reconciles `SparkApplication` (one batch/streaming job, run to
completion) and `SparkCluster` (a long-lived standalone cluster)
custom resources into running Spark workloads.

## The Deployment Landscape

One operator serves the whole cluster under the default posture: it
watches every namespace. The catalog deliberately models only the
engine — `SparkApplication` objects are per-JOB, run-to-completion
declarations, submitted per pipeline run (typically from an
orchestrator such as KubernetesAirflow) or declared via
KubernetesManifest. A per-job declaration has the lifecycle of a
pipeline run, not of infrastructure, which is why it is not a catalog
kind.

### The chart/operator pinning truth

Chart 1.8.0 pairs with operator 1.0.0 (the chart's appVersion) — the
newest SERVED chart, verified against the repository index. Versions
must exist as served charts in the repository index.

## The CRD Lifecycle: crds/ Directory, Keep by Construction

The chart ships its two spark.apache.org CRDs (`sparkapplications`,
`sparkclusters`) from its `crds/` DIRECTORY. Helm's contract for that
directory: install the CRDs once, never upgrade them on chart
upgrades, and leave them on uninstall (no release ownership metadata).

That upstream posture is exactly the keep-on-uninstall this catalog
wants for workload-bearing CRDs — removing the operator never deletes
a Spark workload declaration — so the modules neither re-own nor
template the CRDs. The one consequence to know: chart-version bumps
never touch the CRDs, so apply the new release's CRD files manually
when a bump changes them, per the upstream release notes.

## No Admission Webhook

This operator validates in its reconcile loop. There is no admission
webhook, no certificate machinery, and no cert-manager dependency —
one less lifecycle to manage and nothing that can fail-close the
cluster's write path. A bad SparkApplication surfaces on the CR's
status, not as an admission rejection.

## The Workload Surface: Fence = Watch Scope = RBAC

Spark driver pods create and delete executor pods at runtime, so
every namespace that runs Spark needs a service account with real
pod-management permissions. The chart owns that surface
(`workloadResources`), and the spec models it as ONE decision
(`spec.workload`):

- **Empty `namespaces` (cluster-wide)**: the operator watches every
  namespace, the chart creates a workload ClusterRole, and the
  workload service account lands in the release namespace.
- **Non-empty `namespaces` (fenced)**: the chart CREATES each listed
  namespace (know this before pointing at names you manage
  elsewhere), plants the service account and a namespace-scoped
  Role/RoleBinding in each, drops the workload ClusterRole, and the
  chart's `overrideWatchedNamespaces` (default true) wires the
  operator's `spark.kubernetes.operator.watched.namespaces` property
  from the same list — one value, one truth. SparkApplications
  anywhere else are ignored without an error, so a missing namespace
  in the list looks like a job that never starts.

The chart marks workload resources `helm.sh/resource-policy: keep`:
the fenced namespaces and their RBAC SURVIVE uninstall by upstream
design — running jobs must not lose their identity mid-flight.

The workload service account name stays the upstream convention
(`spark` unless overridden): SparkApplications reference it by that
conventional name.

## Multi-Instance Safety: The RBAC Name Re-Pins

The chart hardcodes all its cluster-scoped RBAC names as plain values
(`spark-operator-clusterrole`, …) — a second install anywhere on the
cluster would collide by construction. The modules derive every RBAC
name from the release identity instead
(`<name>-clusterrole`, `<name>-clusterrolebinding`,
`<name>-config-monitor`, `<name>-config-monitor-binding`,
`<name>-workload-clusterrole`, `<name>-workload-role`,
`<name>-workload-rolebinding`), so instances coexist — the same
defense as the fullname pin, applied to the chart's values-borne
names. The chart derives the workload binding's roleRef ITSELF from
`clusterRole.create` (ClusterRole when cluster-wide, Role when
fenced) — only the name is the module's to pin (template-verified).

The 40-character name budget follows: the longest derived suffix is
`-config-monitor-binding` (23 chars) against the Kubernetes
63-character name limit, and both modules fail loudly over the
budget.

## Properties-Based Configuration

The operator is properties-file configured
(`spark-operator.properties`, keys like
`spark.kubernetes.operator.reconciler.intervalSeconds`). The chart
APPENDS the module-rendered document over its built-in defaults
(`operatorConfiguration.append`, chart default true — kept). The full
key catalog ships with the operator's docs at the pinned version.

Two related surfaces:

- **Leader election is module-owned, never a spec knob.** Any replica
  count beyond 1 REQUIRES leader election (the chart's own contract —
  it refuses multi-replica installs without it), so the module renders
  `spark.kubernetes.operator.leaderElection.enabled=true` exactly when
  `replicas > 1`. A knob could drift from the replica count; a
  derivation cannot.
- **`dynamic_config` is the hot-reload arm**: the chart creates a
  ConfigMap the operator re-reads at runtime, plus the RBAC that lets
  it watch that ConfigMap (`operatorRbac.configManagement`, chart
  default true — kept). Changes to the hot properties apply WITHOUT an
  operator restart. Off by default, the upstream default — most
  installs prefer restart-on-change semantics.

## Image Form

The operator image is combined-form: the chart's default is
`apache/spark-kubernetes-operator` (Docker Hub implied).
`image_registry` replaces ONLY the registry part — the air-gap mirror
seam for the one image THIS component's pods pull. Spark WORKLOAD
images ride each SparkApplication's own image field; this never
rewrites those.

## Design Decisions

- **The install is blocking.** The Helm release waits for the operator
  Deployment to become Available (atomic, 600s timeout, cleanup on
  fail) — a JVM with a 30s-initial-delay startup probe. An unpullable
  image or a broken config fails THIS apply with a readiness timeout
  instead of surfacing later as SparkApplications that mysteriously
  never reconcile.
- **The module owns namespace creation** (`create_namespace`), never
  the Helm release — and only for the INSTALLATION namespace. Workload
  namespaces are chart-created and chart-kept, deliberately not module
  resources.
- **Chart-default-matching values render only on divergence** — with
  one deliberate always-rendered block: the RBAC name re-pins.
- **No post-merge re-pin document exists**: this chart has no
  release-owned CRDs and no webhook machinery whose keys an
  escape-hatch value could weaponize. The RBAC name re-pins live in
  the first (typed) values document; an operator deliberately
  overriding them via `helm_values` owns the collision consciously.
- **The chart's tuned JVM defaults stand** (parallel GC, 80% RAM
  percentage, crash on OOM; 1 CPU / 2Gi with requests = limits) unless
  the spec diverges — the upstream-tested posture.

## Version Pins and Naming Contracts

| What | Value | Notes |
|---|---|---|
| Chart | `spark-kubernetes-operator` at https://apache.github.io/spark-kubernetes-operator | Pinned 1.8.0 (spec default) |
| Operator image | `apache/spark-kubernetes-operator` (Docker Hub implied) | Chart 1.8.0 = operator 1.0.0 |
| CRD API group | `spark.apache.org` | Two CRDs (`sparkapplications`, `sparkclusters`), chart `crds/` directory — installed once, kept on uninstall, never upgraded by chart bumps |
| Workload service account | `spark` (upstream convention) | Exported as `workload_service_account`; SparkApplications reference it |
| RBAC names | `<name>-clusterrole`, `<name>-config-monitor-binding`, `<name>-workload-*`, … | Release-derived (the chart hardcodes them); the 23-char suffix behind the 40-char name budget |

## IaC Twins

Pulumi (`module/values.go` + `module/locals.go`) and Terraform
(`main.tf` + `locals.tf`) render identical chart values: the same
RBAC name re-pins, the same fence/watch-scope wiring, the same
module-owned leader-election property, and the same blocking-install
posture. Keep the typed-value rendering and the RBAC derivations in
lockstep.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
