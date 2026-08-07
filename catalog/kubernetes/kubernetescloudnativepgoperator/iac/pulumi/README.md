# KubernetesCloudNativePgOperator Pulumi Module

Installs CloudNativePG from the official Helm charts
(`https://cloudnative-pg.github.io/charts` — one repository serves both
charts) as up to TWO real Helm releases in the same namespace. The typed
spec renders into operator-chart values in `module/values.go`; the
`helm_values` escape hatch merges LAST over them with Helm `-f`
semantics (maps deep-merge, later document wins, lists replace) — the
exact semantic twin of the Terraform module's `helm_release` with
`values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true; otherwise the namespace must
   already exist
2. **Helm Release `cnpg`** — the operator chart (`cloudnative-pg`,
   pinned default 0.29.0 = operator 1.30.0). The release name is FIXED:
   the operator registers cluster-scoped CRDs and webhooks whose service
   name is baked into the chart (`cnpg-webhook-service` — embedded in
   the webhook certificate and not configurable), so one installation
   per cluster is an upstream constraint and the name never derives from
   `metadata.name`
3. **Helm Release `plugin-barman-cloud`** (when
   `barman_cloud_plugin.enabled`) — the Barman Cloud CNPG-I plugin chart
   (pinned default 0.7.0 = plugin v0.13.0), a SEPARATE release in the
   SAME namespace: upstream forbids folding the plugin into the
   operator's release (Helm ownership of shared resources would
   conflict). Installed AFTER the operator so the plugin's CNPG-I
   registration always lands on a running operator; uninstall unwinds in
   reverse. Its release name is fixed for the same singleton reason (the
   plugin's gRPC service name `barman-cloud` is baked into its TLS
   certificate)

## Rendering Notes

- **Chart-default-matching values render only on divergence** —
  `crds.create` (only on explicit opt-out), `config.clusterWide` (only
  when fencing), monitoring flags (only when on) — the rendered values
  stay minimal on both engines.
- **The typed watch field OWNS `WATCH_NAMESPACE`** — a user entry under
  that key in `operator_config` is always stripped; the key renders only
  from `watch.namespaces` (comma-joined) when `cluster_wide` is false.
  Spec CEL rules guarantee namespaces are present exactly when the fence
  is on.
- **The chart folds three typed concerns into one `config` block** —
  `clusterWide`, `data` (the operator's ConfigMap entries), and
  `maxConcurrentReconciles` — rendered together so the block appears at
  most once.
- **No `fullnameOverride`** — the chart hard-codes the names that matter
  (the webhook service is `cnpg-webhook-service` regardless of release
  name); there is nothing for an override to pin.
- **`helm_values` scopes to the OPERATOR chart only** — the plugin
  release renders from its own typed fields (container resources only;
  everything else rides the plugin chart's defaults). The two charts
  share value keys like `resources` and `image`, so forwarding one
  document to both would misconfigure the plugin.
- **CRD keep is unconditional** — the chart stamps
  `helm.sh/resource-policy: keep` on every CRD, so uninstalling never
  cascade-deletes the Cluster resources (and the databases behind them);
  no keep knob is needed or modeled.

## Wait / Atomic Posture

Both releases install with `Atomic` + `CleanupOnFail` and a 600s
timeout, waiting for readiness. An operator that never becomes ready (a
PodMonitor rendered without the Prometheus operator CRDs is THE classic
install failure) fails THIS deploy with a readiness timeout instead of
surfacing later as Cluster resources that mysteriously never reconcile.
The plugin release is where the cert-manager dependency surfaces: the
plugin chart renders cert-manager Issuer/Certificate resources
unconditionally, and without cert-manager on the cluster its
Certificates never become ready — the release rolls back with a clear
timeout.

## Usage

```shell
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator (and the plugin, when enabled) runs in |
| `release_name` | Helm release name of the operator (fixed `cnpg` — one installation per cluster) |
| `barman_plugin_release_name` | Helm release name of the plugin when enabled (`plugin-barman-cloud`); empty otherwise — the handle KubernetesPostgres backup blocks key off |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → operator release → plugin release →
  output exports
- `module/values.go`: typed-spec → operator-chart values rendering (CRD
  lifecycle, sizing, the config block with the WATCH_NAMESPACE
  precedence, telemetry, scheduling, image), the escape-hatch merge, and
  the plugin chart's minimal values
- `module/locals.go`: resolved namespace, chart versions, and the plugin
  arm — kept in lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identities, pinned default versions (0.29.0 =
  operator 1.30.0; plugin 0.7.0 = v0.13.0), the fixed release names, the
  600s timeout
- `module/helpers.go`: shared shape renderers (resources, tolerations,
  the Helm `-f` merge)
