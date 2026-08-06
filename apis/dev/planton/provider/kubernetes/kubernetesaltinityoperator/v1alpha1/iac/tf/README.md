# KubernetesAltinityOperator Terraform Module

Installs the Altinity ClickHouse operator from the official
`altinity-clickhouse-operator` Helm chart
(`https://docs.altinity.com/clickhouse-operator/`) as a single Helm
release named after `metadata.name`. The typed spec renders into chart
values in `locals.tf` (`local.typed_values`); the `helm_values` escape
hatch is passed as a SECOND values document the provider merges over
the first with Helm `-f` semantics — the exact semantic twin of the
Pulumi module's `buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The CHART (not the module) owns the CRDs** — the four CRDs ship in
  the chart's `crds/` directory, so Helm installs them on first install
  and NEVER deletes them on uninstall (keep-on-uninstall is inherent:
  removing the operator never cascade-deletes `ClickHouseInstallation`
  resources or their data). CRD UPGRADES are carried by the chart's
  pre-install/pre-upgrade hook job (`crdHook`, enabled by default),
  which server-side-applies the CRDs on every install and upgrade. The
  module vendors nothing and leaves `skip_crds` false — unlike sibling
  operator modules whose charts template CRDs release-owned.
- **`fullnameOverride` is pinned to `metadata.name`** and re-pinned
  AFTER the `helm_values` merge — the one deliberate exception to the
  escape hatch's last-word contract. The Deployment, the credentials
  Secret and the metrics Service names (and the exported outputs built
  from them) all derive from the fullname. Keep the resource name at 39
  characters or fewer: the longest generated child name
  (`<fullname>-keeper-templatesd-files`) adds 24 characters against the
  Kubernetes 63-character cap.
- **The pinned default chart version is 0.27.2** — chart versions track
  operator releases one-to-one (chart 0.27.2 = operator image 0.27.2);
  pick versions from the SERVED repository index
  (`https://docs.altinity.com/clickhouse-operator/index.yaml`).
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout. An operator that never becomes
  ready (an unpullable image from a private mirror is the classic case)
  fails THIS apply with a readiness timeout instead of surfacing later
  as ClickHouse clusters that never reconcile.
- **The module (not Helm) owns namespace creation** —
  `create_namespace` drives a `kubernetes_namespace_v1` resource
  carrying the standard governance labels;
  `helm_release.create_namespace` is always false.

## Values Mapping

| Spec field | Chart value |
|---|---|
| `watch_namespaces` | `watchNamespaces` (entries are regexps) |
| `namespace_scoped_rbac` | `rbac.namespaceScoped` |
| `operator_credentials.username` | `secret.username` |
| `operator_credentials.password` | `secret.password` |
| `metrics.enabled` | `metrics.enabled` |
| `metrics.resources` | `metrics.resources` |
| `crd_hook.enabled` | `crdHook.enabled` |
| `crd_hook.image.repo` / `.tag` | `crdHook.image.repository` / `.tag` |
| `resources` | `operator.resources` |
| `image.repo` / `.tag` | `operator.image.repository` / `.tag` |
| `service_monitor_enabled` | `serviceMonitor.enabled` |
| `node_selector` | `nodeSelector` |
| `tolerations` | `tolerations` |
| `image_pull_secrets` | `imagePullSecrets` (list of `{name}`) |
| (always) | `fullnameOverride: <metadata.name>` |
| `helm_values` | merged LAST, Helm `-f` semantics (`fullnameOverride` re-pinned after it) |

## Rendering Quirks

- **Chart-default-matching values render only on divergence** — every
  key renders only when the spec carries a value, so an empty spec
  installs the chart exactly as upstream ships it (except the
  always-pinned `fullnameOverride`), and the rendered values carry no
  empty sections.
- **The `secret` block is omitted ENTIRELY when `operator_credentials`
  is absent** — the chart then provisions its publicly documented
  default credentials (`clickhouse_operator` /
  `clickhouse_operator_password`), which is why the spec calls unset
  UNSAFE FOR PRODUCTION. When present, the proto's username default
  (`clickhouse_operator`) resolves in `locals.tf`, so a password-only
  spec still names the standard operator user.
- **The image overrides deep-merge per half** —
  `operator.image.repository` / `operator.image.tag` (and the crdHook
  pair) render independently, leaving the chart's `registry`,
  `pullPolicy` and the appVersion-derived default tag intact.
- **The `resources` keys render only when the spec sets them** — the
  chart ships no default requests/limits, and Helm deep-merges per key,
  so a partial spec block overrides only the halves it carries.
- **The presence-tracked toggles render on presence** —
  `metrics.enabled` and `crdHook.enabled` (both true-defaulted in the
  chart) render whenever the spec carries a value; the plain bools
  (`rbac.namespaceScoped`, `serviceMonitor.enabled`, both
  false-defaulted) render only when true.
- **Pull secrets render as the raw Kubernetes object list**
  (`imagePullSecrets: [{name: ...}]`) — the chart pipes the list
  straight into the operator pod spec.
- **Null-prune idiom throughout** — conditional entries are written as
  `key = cond ? value : null` inside one object literal and pruned, so
  numbers and booleans keep their types in the rendered YAML.

## Resources

| Resource | Condition |
|---|---|
| `kubernetes_namespace_v1.altinity_operator` | `spec.create_namespace` |
| `helm_release.altinity_operator` | always |

## Usage

```bash
planton tofu apply --manifest altinity-operator.yaml
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
snake_case with the StringValueOrRef foreign keys — `namespace`
(KubernetesNamespace) and `operator_credentials.password` — resolved to
literal strings before Terraform runs.

## Outputs

| Output | Description |
|--------|-------------|
| `namespace` | Namespace the operator runs in |
| `release_name` | Helm release name of the operator (`metadata.name`) |
| `deployment_name` | The operator Deployment name — the chart names it exactly the fullname, which the module pins to `metadata.name` |
| `credentials_secret_name` | The chart-managed credentials Secret (keys `username`/`password`) — named exactly the fullname; created unconditionally by the chart (with the documented defaults when `operator_credentials` is unset) |
| `metrics_endpoint` | `http://<name>-metrics.<namespace>.svc.cluster.local:8888/metrics` — the chart's `<fullname>-metrics` Service carries the 8888 `ch-metrics` port only while `metrics.enabled`; empty when the exporter is disabled |

## Parity

Kept in lockstep with the Pulumi module (`../pulumi/module/`): same
chart identity and pinned default version (0.27.2), same
`metadata.name` release name, same values rendering (the pinned and
re-pinned `fullnameOverride`, the divergence-only blocks, the
omit-when-absent secret, the per-half image overrides), same
chart-owned CRD posture (`skip_crds` false on both engines), same
atomic/wait posture, same outputs.
