# KubernetesClusterAutoscaler Pulumi Module

Installs the Kubernetes Cluster Autoscaler from the official Helm chart
(`cluster-autoscaler` at `https://kubernetes.github.io/autoscaler`) as a
real Helm release. The typed spec renders into chart values in
`module/values.go`; the `helm_values` escape hatch merges LAST over them
with Helm `-f` semantics (maps deep-merge, later document wins, lists
replace) — the exact semantic twin of the Terraform module's `helm_release`
with `values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance labels
   when `create_namespace` is true; otherwise the namespace must already
   exist (always the case for the upstream-conventional `kube-system`
   target)
2. **Helm Release** — named `cluster-autoscaler`, FIXED: the autoscaler
   leader-elects and owns the cluster-wide scaling decision, so one
   installation per cluster is the operating model and the release name
   never derives from `metadata.name`

## Rendering Quirks

- **The Deployment has a render gate** (`templates/deployment.yaml` line
  1): it only renders when `autoDiscovery.clusterName` /
  `autoDiscovery.namespace` / `autoDiscovery.labels` or
  `autoscalingGroups` is set. `autoscalingGroupsnamePrefix` — the GCE
  contract — does NOT satisfy the gate, and kwok sets no gate key at all.
  The gce and kwok arms therefore render
  `autoDiscovery.clusterName = metadata.name` (a benign, deterministic
  value no gce/kwok template consumes beyond the gate; the chart README
  documents "any-name" as required for GCE). Without it the release would
  "succeed" while installing NO autoscaler pod.
- **Civo node pools ride `helm_values`** — the chart requires
  `autoscalingGroups` for civo (its README), but the typed spec does not
  model civo pools; supplying them through `helm_values` also satisfies
  the render gate.
- **Cluster API discovery needs at least one dimension** — the typed
  `namespace` field satisfies the gate when set; with it empty, ship
  `autoDiscovery.clusterName` or `autoDiscovery.labels` via `helm_values`
  per upstream's Cluster API guidance.
- **Credentials become the chart's own Secret** — AWS access keys, all
  Azure identity fields, and the Civo credentials land in chart values
  which the chart materializes into `templates/secret.yaml` and wires to
  the pod as env vars via `secretKeyRef`. The module never logs them.
- **`extraArgs` precedence** — the typed `scaling` block renders first,
  `extra_args` merges over it (user wins per key), and the chart's own
  extraArgs defaults (`logtostderr`/`stderrthreshold`/`v`) survive
  chart-side because Helm coalesces the provided map over the chart
  default per key. Values render as strings on both engines — they are
  CLI flag text.
- **Azure workload identity sets only the extension flag** — VERIFIED: no
  chart template adds the `azure.workload.identity/use` pod label;
  clusters relying on the azure-workload-identity webhook add `podLabels`
  via `helm_values`.

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and waits (600s
timeout) for the Deployment to become Available. An autoscaler that never
becomes ready — bad cloud credentials crash-looping the pod is THE classic
failure; a ServiceMonitor rendered without the Prometheus operator CRDs the
other — fails THIS deploy with a readiness timeout instead of surfacing
later as node groups that mysteriously never scale.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Kubernetes namespace the autoscaler was installed into |
| `release_name` | Helm release name (fixed `cluster-autoscaler` — one installation per cluster) |
| `service_account_name` | The autoscaler's service account, derived from the chart's fullname template: `cluster-autoscaler-<cloudProvider>-cluster-autoscaler` (verified in `_helpers.tpl` — the default name is `<cloudProvider>-<chartName>`, which never equals the release name, so fullname prefixes the release). The subject IRSA trust policies / GCP WI bindings / Entra federated credentials are written against |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → Helm release → output exports
- `module/values.go`: typed-spec → chart values rendering (one provider
  arm per install, the deployment render gate, scaling flags → extraArgs
  with user-wins precedence, deployment sizing, ServiceMonitor),
  escape-hatch merge
- `module/locals.go`: resolved namespace, chart version, cloudProvider,
  and the derived service-account name — kept in lockstep with the
  Terraform module's `locals.tf`
- `module/vars.go`: chart identity, pinned default chart version (9.59.0
  — ships autoscaler 1.35.0; chart and app versions move separately), and
  the fixed release name
