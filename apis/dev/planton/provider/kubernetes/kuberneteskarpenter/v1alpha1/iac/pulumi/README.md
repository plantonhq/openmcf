# KubernetesKarpenter Pulumi Module

Installs Karpenter from the official OCI-served Helm charts
(`oci://public.ecr.aws/karpenter/karpenter` and the companion
`karpenter-crd` chart) as TWO real Helm releases in the same namespace. The
typed spec renders into controller-chart values in `module/values.go`; the
`helm_values` escape hatch merges LAST over them with Helm `-f` semantics
(maps deep-merge, later document wins, lists replace) — the exact semantic
twin of the Terraform module's `helm_release` with
`values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true; otherwise the namespace must already
   exist (always the case for `kube-system`, upstream's recommended home)
2. **`karpenter-crd` Helm Release** (when `crds.install`, default true) —
   the companion CRD chart, upstream's supported mechanism for keeping CRDs
   upgradable: Helm installs the copies bundled inside the main chart once
   and NEVER upgrades them. When `keep_on_uninstall` (default true) the
   release stamps `helm.sh/resource-policy: keep` onto every CRD through
   the chart's `additionalAnnotations` — without it a plain uninstall
   cascade-deletes every NodePool/EC2NodeClass/NodeClaim in the cluster.
3. **`karpenter` Helm Release** — the controller chart, installed with
   `SkipCrds` UNCONDITIONALLY (the CRD release owns the CRDs; skipping
   keeps this release's shape deterministic whether or not `crds.install`
   is on) and `depends_on` the CRD release.

Both release names are FIXED: Karpenter owns the cluster-wide
`karpenter.sh` label domain, its CRDs, and node lifecycle — one
installation per cluster is an upstream constraint, so the names never
derive from `metadata.name`.

## OCI Chart Resolution (engine asymmetry)

Pulumi's `helm.v3.Release` resolves `oci://` registries through the CHART
REFERENCE — the joined `oci://public.ecr.aws/karpenter/<chart>` string with
NO `RepositoryOpts`. The Terraform provider instead takes
`repository = "oci://public.ecr.aws/karpenter"` plus the bare chart name
and joins them internally. Same chart bytes, different wiring — both
modules carry the split with explanatory comments.

## Rendering Notes

- **Defaults are applied module-side**: fields whose spec default mirrors
  the chart default (replicas, logLevel, batching windows, scheduling
  policies, featureGates, priorityClassName, reservedENIs,
  vmMemoryOverheadPercent) render with the default APPLIED, so both
  engines produce byte-identical values whether or not the platform's
  defaulting middleware ran.
- **Type fidelity**: the chart's `settings.reservedENIs` default is the
  STRING `"0"` — the spec's int32 renders as a string. The chart's
  `settings.vmMemoryOverheadPercent` default is the NUMBER `0.075` — the
  spec's string parses and renders as a number. Rendering the other type
  would change the values document the chart sees.
- **The whole `featureGates` map renders** with defaults applied — the
  deployment template composes the `FEATURE_GATES` env var from all six
  keys unconditionally, and `reservedCapacity` defaults TRUE unlike the
  other five gates.
- **`nodeSelector` merges, `tolerations` replace**: Helm deep-merges maps,
  so spec entries narrow the chart's `kubernetes.io/os: linux` default —
  but lists replace wholesale, so tolerations render only when the spec
  provides them (an empty list would DROP the chart's CriticalAddonsOnly
  default).
- **IRSA rides the service-account annotation**
  (`eks.amazonaws.com/role-arn`), not a settings entry; empty means EKS
  Pod Identity, configured cloud-side with no annotation needed.

## Wait / Atomic Posture

Both releases install with `Atomic` + `CleanupOnFail` and wait (600s
timeout) for readiness. A Karpenter that never becomes ready (a
ServiceMonitor rendered without the Prometheus operator CRDs, a bad IRSA
trust policy) fails THIS deploy with a readiness timeout instead of
surfacing later as pods that stay Pending forever because no nodes ever
appear.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Kubernetes namespace Karpenter was installed into |
| `release_name` | Controller Helm release name (fixed `karpenter`) |
| `crd_release_name` | CRD Helm release name (fixed `karpenter-crd`; empty when `crds.install` is false) |
| `service_account_name` | The controller's service account (fixed `karpenter` — SA name defaults to the chart fullname, which resolves to the release name because the release name contains the chart name); the subject IRSA trust policies and EKS Pod Identity associations are written against |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → CRD release → controller release → output
  exports, carrying the OCI join and the unconditional SkipCrds
- `module/values.go`: typed-spec → chart values rendering (cluster
  identity, the AWS arm with its two type-fidelity conversions, controller
  sizing, batching, scheduling posture, the full featureGates map,
  controller-pod scheduling, ServiceMonitor), escape-hatch merge
- `module/locals.go`: resolved namespace, chart version, and CRD lifecycle
  flags — kept in lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity (OCI repo + both chart names), pinned
  default version (1.14.0 — both charts version with the controller), the
  fixed release/service-account names, and the chart-default table
