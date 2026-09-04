# KubernetesPlantonOperator Pulumi Module

Installs the Planton operator from the official `planton-operator` Helm
chart (OCI, `ghcr.io/plantonhq/charts`) as ONE real Helm release,
byte-identical to a hand-installed one; the chart owns its two definitions
as release resources. The typed spec renders
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
- **The chart owns its definitions; the module maps two dials onto them** —
  the `PlantonPlatform` and `PlantonIdentityProvider` CRDs are resources of
  the release behind the chart's `crds.enabled` / `crds.keep` values, which
  `buildHelmValues` renders from `spec.crds.install` and
  `spec.crds.keep_on_uninstall` (both default true). A `chart_version`
  bump upgrades the operator and its schema together; a destroy keeps the
  definitions (Helm's `helm.sh/resource-policy: keep`) unless keeping is
  explicitly disabled. The module carries no copy of the schema and
  applies no CRD itself; `SkipCrds` is never set, because it governs only
  Helm's install-once `crds/` directory, which this chart does not use.
- **Reinstall adopts the kept definitions** — kept definitions carry the
  release's identity in their Helm ownership metadata, and the release
  name is fixed, so every later install of the operator on the cluster
  owns them again.
- **The chart-version floor fails loudly** (`chartVersionAtLeast`) —
  charts below 0.8.0 do not own their definitions and have no `crds`
  values, so the dials would be silently dropped; the refusal names the
  version to pin. Twin: the Terraform module's lifecycle precondition.
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
