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
2. **The OpenSearch CRDs** (module-owned) — derived from the pinned
   chart, one resource per CRD keyed by the CRD's own `metadata.name`,
   with `retainOnDelete`
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
The module therefore owns the CRD lifecycle itself through the catalog's
shared package for charts that carry CRDs, `keptcrds`:

- The pinned chart is rendered in-process (Helm SDK, client-only, against
  the cluster's own Kubernetes version) with the release's own values
  plus `installCRDs: true`; the CustomResourceDefinition documents are
  kept, stamped with `planton.ai/crd-source-chart` and
  `planton.ai/crd-source-version`, and applied one `yaml.ConfigGroup`
  per CRD keyed by the CRD's own `metadata.name`, ahead of the release.
- `installCRDs: false` is re-pinned AFTER the escape-hatch merge and the
  release installs with `SkipCrds`, so Helm never touches the CRDs
  whichever way a user override points. `crds.install: false` is the
  bring-your-own-CRDs arm (the CRDs are owned elsewhere).
- Every applied CRD carries `retainOnDelete` (delivered through a
  resource transformation, the one option the yaml SDK propagates to a
  ConfigGroup's children): on `pulumi destroy` the CRDs are dropped from
  state WITHOUT being deleted from the cluster, so destroying the
  operator never touches `OpenSearchCluster` resources or their data.
  `crds.keep_on_uninstall: false` turns it off. The CRDs ride a
  dedicated upsert provider (server-side apply with force) so a
  reinstall re-adopts the kept CRDs. This is the exact semantic twin of
  the Terraform module's `kubectl_manifest` with `apply_only`.
- A `chart_version` bump re-applies the derived CRDs at the new pin (and
  adds the ones a newer chart brings); a `chart_version` below the
  version the cluster's CRDs carry is refused before anything registers,
  with what was observed, what it means, and the next step. A version
  that is not published in the repository index is refused the same way.

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
planton pulumi up --manifest e2e/manifest.yaml --module-dir <path-to-this-module>
```

## Outputs

| Output | Description |
|---|---|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`) |
| `deployment_name` | The controller-manager Deployment name, derived exactly as the chart's own fullname helper does it: `<release>` when the release name already contains `opensearch-operator`, otherwise `<release>-opensearch-operator`, truncated to 63 chars with a trailing `-` trimmed, plus `-controller-manager`. Valid as long as `helm_values` does not set `nameOverride`/`fullnameOverride` |

## Module Structure

- `main.go`: entrypoint that calls the module
- `module/main.go`: namespace → release values → derived CRDs
  (`keptcrds.Apply`) → operator release → output exports
- `module/values.go`: typed-spec → chart values rendering (the pinned
  `installCRDs: false`, the manager block, kube-rbac-proxy, RBAC scope,
  scheduling, image) and the escape-hatch merge
- `module/locals.go`: resolved namespace, release name
  (`metadata.name`), chart version, the chart-derived deployment name —
  kept in lockstep with the Terraform module's `locals.tf`
- `module/vars.go`: chart identity (repository, `opensearch-operator`),
  the pinned default version (2.8.0), the CRD render override, the
  600s timeout
- `module/helpers.go`: shared shape renderers (resources, tolerations,
  the Helm `-f` merge)
