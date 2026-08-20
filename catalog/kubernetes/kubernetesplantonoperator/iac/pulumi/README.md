# KubernetesPlantonOperator Pulumi Module

Installs the Planton operator from the official `planton-operator` Helm
chart (OCI, `ghcr.io/plantonhq/charts`) as ONE real Helm release plus the
module-owned `plantonplatforms.planton.ai` CRD. The typed spec renders
into chart values in `values.go` (`buildHelmValues`); the `helm_values`
escape hatch merges LAST with Helm `-f` semantics (`mergeMaps`) — the
exact semantic twin of the Terraform module's two-document `values` list.

## Module Behavior

- **The release name is FIXED to `planton-operator`** — the operator
  enforces one installation per cluster itself at startup (a
  label-matched Deployment scan that refuses to start beside a sibling),
  so the name never derives from `metadata.name`: the collision is
  impossible to express instead of merely refused. With the release name
  equal to the chart name, the chart's fullname helper renders the
  Deployment as plain `planton-operator` — the name the startup guard's
  labels ride.
- **The CRD is module-owned** (`crds.go`) — applied from the copy staged
  at `../crds` (extracted from the published chart at the pinned default
  version) BEFORE the release, which installs with `SkipCrds`. The chart
  ships its CRD in Helm's install-once `crds/` directory (never upgraded,
  never removed); module ownership is what makes `chart_version` bumps
  carry the CRD and keep-on-uninstall a guarantee.
- **Keep-on-uninstall via a retainOnDelete TRANSFORMATION** — the option
  must reach the ConfigGroup's children, and only a transformation
  propagates there (ordinary resource options do not). The exact semantic
  twin of the Terraform module's `apply_only = true`.
- **Reinstall ADOPTS the retained CRD** — the CRD applies through the
  UPSERT provider (`GetWithKubernetesProviderConfigUpsert`): a destroy
  leaves the CRD by design, so the next install's server-side apply must
  adopt it instead of failing AlreadyExists. Only the CRD rides that
  provider; every other resource keeps create-conflict semantics.
- **The chart-version floor fails loudly** (`chartVersionAtLeast`) —
  charts below 0.7.0 ship operators whose reconcilers predate the
  PlantonPlatform schema the staged CRD advertises; the API server would
  ACCEPT fields the running operator silently ignores, so the module
  refuses instead. Twin: the Terraform module's lifecycle precondition.
- **A fail-loud staged-file count** — an empty `../crds` would silently
  apply ZERO CRDs; the module requires exactly the staged count (1 at
  chart 0.7.1; re-stage and update together with `DefaultChartVersion`).
- **Readiness is verified at install time** — `Atomic` +
  `CleanupOnFail` with a 600s timeout: a manager that never becomes
  ready (the startup guard refusing beside a sibling operator is THE
  classic case) fails THIS deploy with a readiness timeout instead of
  surfacing later as declarations that never reconcile.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a Namespace resource carrying the standard governance labels;
  the release's `CreateNamespace` is always false.

## OCI Quirk

Pulumi's `helm.v3.Release` does not resolve `oci://` through
`RepositoryOpts` the way the Terraform provider does — the chart
reference is the JOINED `oci://ghcr.io/plantonhq/charts/planton-operator`
string (see `main.go`).
