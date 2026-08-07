# KubernetesOpenSearchOperator Terraform Module

Installs the OpenSearch Kubernetes Operator from the official
`opensearch-operator` Helm chart
(`https://opensearch-project.github.io/opensearch-k8s-operator/`) as a
single Helm release named after `metadata.name`. The typed spec renders
into chart values in `locals.tf` (`local.typed_values`); the
`helm_values` escape hatch is passed as a SECOND values document the
provider merges over the first with Helm `-f` semantics — the exact
semantic twin of the Pulumi module's `buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The module (not the chart) owns the CRDs** — the chart templates
  its ten CRDs as release-owned resources with NO keep-on-uninstall
  knob upstream, so a Helm-owned install would cascade-delete every
  `OpenSearchCluster` (and its data) on uninstall. The module pins
  `installCRDs: false` UNCONDITIONALLY in the rendered values and
  applies the CRD files staged at `../crds/` itself, one
  `kubectl_manifest` per CRD keyed by the CRD's own `metadata.name`.
- **CRD keep-on-uninstall: `apply_only = true`** — the provider's
  Delete becomes a NO-OP (verified in the provider source: "When true,
  Delete is a no-op"), so destroying this module removes the CRDs from
  state WITHOUT deleting them from the cluster; OpenSearch clusters
  and their data survive an operator uninstall. Server-side apply
  keeps the megabyte-scale CRDs co-ownable without client-side
  annotation bloat. The exact semantic twin of the Pulumi module's
  `retainOnDelete` on each CRD.
- **The release depends on the CRDs** (and the optional namespace), so
  the operator never starts against an unregistered API group.
- **The pinned default chart version is 2.8.0** — the newest SERVED
  chart whose default manager image is a STABLE operator release; the
  2.8.3+/3.0.x served charts default to a prerelease image and the
  3.x line migrates the CRDs to the `opensearch.org` API group. A
  `chart_version` bump must upgrade the staged `../crds/` files
  together with it.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout. An operator that never
  becomes ready (an unpullable image from a private mirror is the
  classic case) fails THIS apply with a readiness timeout instead of
  surfacing later as OpenSearch clusters that never reconcile.
- **The module (not Helm) owns namespace creation** —
  `create_namespace` drives a `kubernetes_namespace_v1` resource
  carrying the standard governance labels;
  `helm_release.create_namespace` is always false.

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

## Rendering Quirks

- **Chart-default-matching values render only on divergence** — every
  `manager.*` key, `kubeRbacProxy.enable`, `useRoleBindings`, and the
  scheduling fields render only when the spec carries a value, so an
  empty spec installs the chart exactly as upstream ships it (except
  the always-pinned `installCRDs: false`).
- **The `manager.resources` key renders only when the spec sets it** —
  the chart SHIPS default requests/limits (requests 100m/350Mi, limits
  200m/500Mi), and an empty spec must leave them intact. Helm
  deep-merges per key, so a partial spec block overrides only the
  halves it carries.
- **The image override deep-merges per half** —
  `manager.image.repository` and `manager.image.tag` render
  independently, leaving the chart's `pullPolicy` (and the
  appVersion-derived default tag) intact.
- **Pull secrets render as the raw Kubernetes object list**
  (`manager.imagePullSecrets: [{name: ...}]`) — the chart pipes the
  list straight into the pod spec.
- **The presence-tracked toggles render on presence** —
  `parallelRecoveryEnabled` and `kubeRbacProxy.enable` (both
  true-defaulted in the chart) render whenever the spec carries a
  value; the plain bools (`pprofEndpointsEnabled`, `useRoleBindings`,
  both false-defaulted) render only when true.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered YAML.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.opensearch_operator` | `spec.create_namespace` |
| `kubectl_manifest.crds` (one per staged CRD, keyed by CRD name) | always |
| `helm_release.opensearch_operator` | always |

## Usage

```bash
planton tofu apply --manifest opensearch-operator.yaml
```

## Local Development

```bash
tofu init -backend=false
tofu validate
tofu plan -var-file=terraform.tfvars.json
tofu apply -var-file=terraform.tfvars.json
```

## Inputs

See `variables.tf` for the full variable specification (generated from
the spec proto). The spec arrives from the proto→tfvars converter in
snake_case with the `namespace` foreign key (KubernetesNamespace)
resolved to a literal string before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`) |
| `deployment_name` | The controller-manager Deployment name, derived exactly as the chart's own fullname helper does it (verified against `templates/_helpers.tpl`): `<release>` when the release name already contains `opensearch-operator`, otherwise `<release>-opensearch-operator`, truncated to 63 chars with a trailing `-` trimmed, plus `-controller-manager`. Valid as long as `helm_values` does not set `nameOverride`/`fullnameOverride` |

## Parity

Kept in lockstep with the Pulumi module (`../pulumi/module/`): same
chart identity and pinned default version (2.8.0), same
`metadata.name` release name, same values rendering (the pinned
`installCRDs: false`, the divergence-only manager block, the
render-only-when-set resources key, the per-half image override), same
module-owned CRD posture (`apply_only` here, `retainOnDelete` there),
same atomic/wait posture, same outputs.
