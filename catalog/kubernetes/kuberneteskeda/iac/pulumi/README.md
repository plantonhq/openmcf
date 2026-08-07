# KubernetesKeda Pulumi Module

Installs KEDA from the official Helm chart (`keda` at
`https://kedacore.github.io/charts`) as a real Helm release. The typed spec
renders into chart values in `module/values.go`; the `helm_values` escape
hatch merges LAST over them with Helm `-f` semantics (maps deep-merge,
later document wins, lists replace) — the exact semantic twin of the
Terraform module's `helm_release` with `values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true; otherwise the namespace must already
   exist
2. **Helm Release** — named `keda`, FIXED: the component registers the
   cluster-wide `v1beta1.external.metrics.k8s.io` APIService, a singleton,
   so one installation per cluster is an upstream constraint and the
   release name never derives from `metadata.name`. No `fullnameOverride`
   either: the chart names its components `keda-operator` /
   `keda-operator-metrics-apiserver` / `keda-admission-webhooks`
   independent of the release name.

## Rendering Quirks

- **The chart's sizing layout is ASYMMETRIC**: replica counts live under
  each component (`operator.replicaCount`, `metricsServer.replicaCount`,
  `webhooks.replicaCount`) while container resources are grouped under ONE
  shared top-level `resources` block keyed per component — and the metrics
  server's key there is `metricServer` (SINGULAR), unlike the
  `metricsServer` component block. `values.go` renders both halves and
  keeps the trap contained.
- **CRD keep mechanism**: the chart has no native keep knob, so
  `keep_on_uninstall` (default true) rides the standard
  `helm.sh/resource-policy: keep` annotation onto the CRDs through the
  chart's `crds.additionalAnnotations` passthrough — otherwise a plain
  uninstall cascade-deletes every ScaledObject/ScaledJob/
  TriggerAuthentication in the cluster. Rendered only when
  `install && keep`.
- **`certificates` renders nil for the default** — type `operator` is the
  chart's own default (the operator self-generates certificates and
  patches the APIService caBundle) and needs no values; `cert_manager`
  renders the `certManager` block, with an explicit issuer as
  `generate: false` + name/kind + group `cert-manager.io`.
- **One `spec.prometheus` flag fans out per component** — the chart
  mirrors its per-component telemetry layout (`prometheus.operator`,
  `prometheus.metricServer`).

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and waits (300s
timeout) for the components to become Available. A KEDA that never becomes
ready — a ServiceMonitor rendered without the Prometheus operator CRDs is
THE classic install failure; broken internal TLS wiring the other — fails
THIS deploy with a readiness timeout instead of surfacing later as
ScaledObjects that mysteriously never scale.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Kubernetes namespace KEDA was installed into |
| `release_name` | Helm release name (fixed `keda` — one installation per cluster) |
| `operator_service_account_name` | The chart's fixed `keda-operator` service account — the subject cloud-side keyless bindings (IRSA trust policies, GCP WI bindings, Entra federated credentials) are written against |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → Helm release → output exports
- `module/values.go`: typed-spec → chart values rendering (CRD lifecycle
  and keep annotation, watch scope, the asymmetric replicas/resources
  layout, webhooks, pod-identity arms, internal TLS, HTTP timeout,
  scheduling, per-component telemetry), escape-hatch merge
- `module/locals.go`: resolved namespace and chart version — kept in
  lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity, pinned default chart version (2.20.1),
  the fixed release name, and the fixed operator service-account name
