# KubernetesPlantonOperator Terraform Module

Installs the Planton operator from the official `planton-operator` Helm
chart (OCI, `ghcr.io/plantonhq/charts` — the provider takes the repo as
`repository` plus the bare chart name, unlike Pulumi's joined string) as
ONE real Helm release plus the module-owned
`plantonplatforms.planton.ai` CRD. The typed spec renders into chart
values in `locals.tf` (`local.typed_values`); the `helm_values` escape
hatch is passed as a SECOND values document the provider merges over the
first with Helm `-f` semantics — the exact semantic twin of the Pulumi
module's `buildHelmValues` + `mergeMaps`.

## Module Behavior

- **The release name is FIXED to `planton-operator`** — the operator
  enforces one installation per cluster itself at startup (a
  label-matched Deployment scan that refuses to start beside a sibling),
  so the name never derives from `metadata.name`. With the release name
  equal to the chart name, the chart's fullname helper renders the
  Deployment as plain `planton-operator`.
- **The CRD is module-owned** (`kubectl_manifest.crds`) — applied from
  the copy staged at `../crds` (extracted from the published chart at the
  pinned default version) BEFORE the release, which installs with
  `skip_crds`. The chart ships its CRD in Helm's install-once `crds/`
  directory (never upgraded, never removed); module ownership is what
  makes `chart_version` bumps carry the CRD and keep-on-uninstall a
  guarantee.
- **Keep-on-uninstall via `apply_only = true`** — "When true, Delete is
  a no-op" (provider source): destroying the operator never
  cascade-deletes platform declarations. Twin of the Pulumi module's
  retainOnDelete transformation.
- **Reinstall ADOPTS the retained CRD** — `server_side_apply` +
  `force_conflicts` adopt a CRD a previous install's destroy retained
  (the field manager differs), natively. Twin of the Pulumi module's
  upsert provider.
- **The chart-version floor fails loudly** (a `lifecycle` precondition;
  HCL has no semver function, so the three parts compare by hand) —
  charts below 0.7.0 ship operators whose reconcilers predate the
  PlantonPlatform schema the staged CRD advertises.
- **A fail-loud staged-file count** — `fileset()` over a missing
  `../crds` returns EMPTY and `for_each` would silently plan ZERO
  resources; a precondition requires exactly the staged count (1 at
  chart 0.7.1; re-stage and update together with
  `default_chart_version`).
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout: a manager that never becomes
  ready (the startup guard refusing beside a sibling operator is THE
  classic case) fails THIS apply instead of surfacing later.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.

## Provider Choice

`alekc/kubectl` for the CRD because `kubectl_manifest` needs no cluster
connection at plan time — a composed infra chart can plan this module
before the cluster (or anything on it) exists.
