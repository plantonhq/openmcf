# KubernetesMetricsServer Pulumi Module

Installs metrics-server from the official Helm chart (`metrics-server` at
`https://kubernetes-sigs.github.io/metrics-server/`) as a real Helm release.
The typed spec renders into chart values in `module/values.go`; the
`helm_values` escape hatch merges LAST over them with Helm `-f` semantics
(maps deep-merge, later document wins, lists replace) — the exact semantic
twin of the Terraform module's `helm_release` with
`values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true; otherwise the namespace must already
   exist (it usually does — kube-system installs leave the flag false)
2. **Helm Release** — named `metrics-server`, FIXED: the component registers
   the cluster-wide `v1beta1.metrics.k8s.io` APIService, a singleton, so one
   installation per cluster is an upstream constraint and the release name
   never derives from `metadata.name`. The chart's `fullnameOverride` is
   pinned to the release name so chart objects (the Service the APIService
   routes to) get deterministic names.

The module owns the chart's `defaultArgs` list: it re-renders the chart's
default argument list with the typed substitutions
(`kubelet_preferred_address_types`, `metric_resolution`) applied, keeping
the pod spec canonical instead of appending duplicate flags.

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and waits (300s
timeout) for the Deployment to become Available. The chart's `/readyz` probe
only passes once the first kubelet scrape succeeds, so a wrong scrape-side
TLS posture (self-signed kubelets without `kubelet_insecure_tls`) fails THIS
deploy with a readiness timeout instead of surfacing later as HPAs that
never scale — and a green deploy means metrics are actually flowing.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Kubernetes namespace metrics-server was installed into |
| `release_name` | Helm release name (always `metrics-server`) |
| `service_name` | The Service the APIService routes to (the chart fullname, pinned to the release name) |
| `api_service_name` | `v1beta1.metrics.k8s.io`; empty when `spec.api_service.create` is false |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → Helm release → output exports
- `module/values.go`: typed-spec → chart values rendering (defaultArgs
  ownership, APIService, TLS arms, telemetry, image), escape-hatch merge
- `module/locals.go`: resolved namespace, chart version, APIService name —
  kept in lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity, pinned default chart version, and the
  fixed release name
