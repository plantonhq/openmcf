# KubernetesPlantonOperator Terraform Module

Installs the Planton operator from the official `planton-operator` Helm
chart (OCI, `ghcr.io/plantonhq/charts` — the provider takes the repo as
`repository` plus the bare chart name, unlike Pulumi's joined string) as
ONE real Helm release, byte-identical to a hand-installed one; the chart
owns its two definitions as release resources. The typed spec renders into chart
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
- **The chart owns its definitions; the module maps two dials onto them** —
  the `PlantonPlatform` and `PlantonIdentityProvider` CRDs are resources of
  the release behind the chart's `crds.enabled` / `crds.keep` values, which
  `local.typed_values.crds` renders from `spec.crds.install` and
  `spec.crds.keep_on_uninstall` (both default true). A `chart_version`
  bump upgrades the operator and its schema together; a destroy keeps the
  definitions (Helm's `helm.sh/resource-policy: keep`) unless keeping is
  explicitly disabled. The module carries no copy of the schema and
  applies no CRD itself; `skip_crds` is never set, because it governs only
  Helm's install-once `crds/` directory, which this chart does not use.
- **Reinstall adopts the kept definitions** — kept definitions carry the
  release's identity in their Helm ownership metadata, and the release
  name is fixed, so every later install of the operator on the cluster
  owns them again.
- **The chart-version floor fails loudly** (a `lifecycle` precondition;
  HCL has no semver function, so the three parts compare by hand) —
  charts below 0.8.0 do not own their definitions and have no `crds`
  values, so the dials would be silently dropped; the precondition names
  the version to pin.
- **Readiness is verified at install time** — `wait` + `atomic` +
  `cleanup_on_fail` with a 600s timeout: a manager that never becomes
  ready (the startup guard refusing beside a sibling operator is THE
  classic case) fails THIS apply instead of surfacing later.
- **The module (not Helm) owns namespace creation** — `create_namespace`
  drives a `kubernetes_namespace_v1` resource carrying the standard
  governance labels; `helm_release.create_namespace` is always false.
