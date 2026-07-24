# KubernetesOpenSearchOperator Pulumi Module

Installs the OpenSearch Kubernetes Operator from the official
`opensearch-operator` Helm chart
(`https://opensearch-project.github.io/opensearch-k8s-operator/`) as a
single Helm release named after `metadata.name`. The typed spec renders
into chart values in `module/values.go`; the `helm_values` escape hatch
merges LAST over them with Helm `-f` semantics (maps deep-merge, later
document wins, lists replace) — the exact semantic twin of the
Terraform module's `helm_release` with `values = [typed, helm_values]`.

## What the Module Creates

1. **Namespace** (optional) — created with the standard governance
   labels when `create_namespace` is true; otherwise the namespace must
   already exist
2. **The ten OpenSearch CRDs** (module-owned) — applied from the files
   staged at `../crds/`, one resource per CRD keyed by the CRD's own
   `metadata.name`, with `retainOnDelete`
3. **Helm Release `<metadata.name>`** — the `opensearch-operator`
   chart (pinned default 2.8.0 — the newest SERVED chart whose default
   manager image is a STABLE operator release; the 2.8.3+/3.0.x served
   charts default to a prerelease image and the 3.x line migrates the
   CRDs to the `opensearch.org` API group), depending on both of the
   above

## CRD Keep-on-Uninstall (the load-bearing design)

The chart templates its CRDs as release-owned resources with NO
keep-on-uninstall knob upstream — a Helm-owned install would
cascade-delete every `OpenSearchCluster` (and its data) on uninstall.
The module therefore owns the CRD lifecycle itself:

- `installCRDs: false` is pinned UNCONDITIONALLY in the rendered
  values (never overridable — it renders before the escape-hatch merge
  but is not a knob the spec models)
- `module/crds.go` applies each staged CRD file as its own resource
  with the `retainOnDelete` option: on `pulumi destroy` the CRDs are
  dropped from state WITHOUT being deleted from the cluster, so
  destroying the operator never touches `OpenSearchCluster` resources
  or their data. This is the exact semantic twin of the Terraform
  module's `kubectl_manifest` with `apply_only = true`.
- A `chart_version` upgrade re-applies the staged CRD files — upgrade
  the staged set together with the pinned default version.

## Rendering Notes

- **Chart-default-matching values render only on divergence** — every
  `manager.*` key, the `kubeRbacProxy.enable` flag, `useRoleBindings`,
  and the scheduling fields render only when the spec carries a value,
  so an empty spec installs the chart exactly as upstream ships it
  (except the always-pinned `installCRDs: false`).
- **The `manager.resources` key renders only when the spec sets it** —
  the chart SHIPS default requests/limits (requests 100m/350Mi, limits
  200m/500Mi), and an empty spec must leave them intact. Helm
  deep-merges per key, so a partial spec block overrides only the
  halves it carries.
- **The image override deep-merges per half** — `manager.image.repository`
  and `manager.image.tag` render independently, leaving the chart's
  `pullPolicy` (and the appVersion-derived default tag) intact.
- **Pull secrets render as the raw Kubernetes object list**
  (`manager.imagePullSecrets: [{name: ...}]`) — the chart pipes the
  list straight into the pod spec.
- **The presence-tracked toggles render on presence** —
  `parallelRecoveryEnabled` and `kubeRbacProxy.enable` (both
  true-defaulted in the chart) render whenever the spec carries a
  value; the plain bools (`pprofEndpointsEnabled`, `useRoleBindings`,
  both false-defaulted) render only when true.

## Values Mapping

| Spec field | Chart value |
|---|---|
| `watch_namespace` | `manager.watchNamespace` |
| `use_role_bindings` | `useRoleBindings` |
| `log_level` | `manager.loglevel` |
| `dns_base` | `manager.dnsBase` |
| `parallel_recovery_enabled` | `manager.parallelRecoveryEnabled` |
| `pprof_endpoints_enabled` | `manager.pprofEndpointsEnabled` |
| `kube_rbac_proxy_enabled` | `kubeRbacProxy.enable` |
| `resources` | `manager.resources` |
| `node_selector` | `nodeSelector` |
| `tolerations` | `tolerations` |
| `image_pull_secrets` | `manager.imagePullSecrets` (list of `{name}`) |
| `image.repository` | `manager.image.repository` |
| `image.tag` | `manager.image.tag` |
| (always) | `installCRDs: false` |
| `helm_values` | merged LAST, Helm `-f` semantics |

## Wait / Atomic Posture

The release installs with `Atomic` + `CleanupOnFail` and a 600s
timeout, waiting for readiness. An operator that never becomes ready
(an unpullable image from a private mirror is the classic case) fails
THIS deploy with a readiness timeout instead of surfacing later as
OpenSearch clusters that mysteriously never reconcile.

## Usage

```shell
planton pulumi up --manifest hack/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`) |
| `deployment_name` | The controller-manager Deployment name, derived exactly as the chart's own fullname helper does it: `<release>` when the release name already contains `opensearch-operator`, otherwise `<release>-opensearch-operator`, truncated to 63 chars with a trailing `-` trimmed, plus `-controller-manager`. Valid as long as `helm_values` does not set `nameOverride`/`fullnameOverride` |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → CRDs → operator release → output exports
- `module/crds.go`: the module-owned CRD apply (per-CRD resources keyed
  by CRD name, `retainOnDelete`)
- `module/values.go`: typed-spec → chart values rendering (the pinned
  `installCRDs: false`, the manager block, kube-rbac-proxy, RBAC scope,
  scheduling, image) and the escape-hatch merge
- `module/locals.go`: resolved namespace, release name
  (`metadata.name`), chart version, the chart-derived deployment name —
  kept in lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity (repository, `opensearch-operator`),
  the pinned default version (2.8.0), the staged-CRD directory, the
  600s timeout
- `module/helpers.go`: shared shape renderers (resources, tolerations,
  the Helm `-f` merge)
