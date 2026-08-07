# KubernetesCilium Pulumi Module

Installs Cilium from the official Helm chart (`cilium` at
`https://helm.cilium.io`) as a real Helm release. The typed spec renders
into chart values in `module/values.go`; the `helm_values` escape hatch
merges LAST over them with Helm `-f` semantics (maps deep-merge, later
document wins, lists replace) — the exact semantic twin of the Terraform
module's `helm_release` with `values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true; otherwise the namespace must already
   exist (kube-system installs leave the flag false — kube-system always
   exists)
2. **Helm Release** — named `cilium`, FIXED: Cilium is the node dataplane
   (the agent DaemonSet, operator, and generated CNI configuration are
   cluster singletons), so one dataplane per cluster is an upstream
   constraint and the release name never derives from `metadata.name`. No
   `fullnameOverride` either: the chart names its workloads with fixed
   names (DaemonSet `cilium`, Deployment `cilium-operator`) regardless of
   the release name, so there is nothing to pin.

## Rendering Quirks

- **`kubeProxyReplacement` is a STRING in the chart's values** (declared
  default `"false"`; historically `"strict"`/`"partial"`) — rendered as
  `"true"` only when the spec flag is set, keeping the values document
  byte-identical with what the chart declares and with the Terraform
  module. **`k8sServicePort` is also a string** in values.yaml (default
  `""`), so the number renders via `strconv.Itoa` — the Terraform twin
  uses `tostring()`.
- **`hubble.metrics.enabled` is upstream's LIST of metric families** (null
  disables) — not a boolean despite the name.
- **`cni.exclusive` renders on presence, not truth** — an explicit false
  is exactly what chaining setups must send (CEL-enforced).
- **One `operator` map, two spec arms** — operator sizing and operator
  telemetry (fanned out from the single `spec.prometheus` toggle together
  with the agent's `prometheus` block) merge into one map instead of the
  later overwriting the earlier.

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and waits with a
600-second timeout — not the usual 300 — because the install path is
heavier than an ordinary workload chart: the agent DaemonSet must roll out
on EVERY node plus the operator, and on a fresh cluster nodes transition
NotReady→Ready only as Cilium wires each one — the rollout itself unblocks
scheduling. A dataplane that never converges fails THIS deploy instead of
surfacing later as pods stuck in ContainerCreating.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Kubernetes namespace Cilium was installed into |
| `release_name` | Helm release name (fixed `cilium` — one dataplane per cluster) |
| `cluster_name` | Cluster identity Cilium runs under (resolved `spec.cluster_name`; `default` when unset) — the name in Hubble flows and any future Cluster Mesh |
| `hubble_relay_service_name` | `hubble-relay` (fixed by the chart) when `hubble.relay` is enabled; empty otherwise |
| `hubble_ui_service_name` | `hubble-ui` (fixed by the chart) when `hubble.ui` is enabled; empty otherwise |
| `gateway_class_name` | `cilium` (fixed by the chart) when `gateway_api` is enabled; empty otherwise |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → Helm release → output exports
- `module/values.go`: typed-spec → chart values rendering (IPAM, routing,
  kube-proxy replacement strings, chaining, cloud arms, Hubble,
  encryption, Gateway API, bandwidth manager, operator/telemetry merge),
  escape-hatch merge
- `module/locals.go`: resolved namespace, chart version, cluster name, and
  the fixed chart-template output names — kept in lockstep with the
  Terraform module's `locals.tf`
- `module/vars.go`: chart identity, pinned default chart version (1.19.6),
  and the fixed release name
