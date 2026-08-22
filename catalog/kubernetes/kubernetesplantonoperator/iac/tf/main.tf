# KubernetesPlantonOperator Terraform module.
#
# Installs the Planton operator — the lifecycle manager that reconciles
# PlantonPlatform declarations (KubernetesPlantonPlatform resources) into
# running self-hosted Planton platforms — from the official
# planton-operator Helm chart (OCI, ghcr.io/plantonhq/charts) as ONE real
# Helm release plus the module-owned PlantonPlatform CRD.
#
# CRD LIFECYCLE IS MODULE-OWNED: the chart ships its CRD in the crds/
# directory — Helm's install-once posture (created on first install, never
# upgraded, never removed). This module applies the CRD itself from the
# copy staged at ../crds (extracted from the published chart at the pinned
# default version) and installs the release with skip_crds, so:
#   - a chart_version upgrade carries the matching CRD update (re-staged
#     with the pin) instead of silently running new operator code against
#     an old schema;
#   - keep-on-uninstall is a guarantee (apply_only below), never an
#     accident — destroying the operator never cascade-deletes
#     PlantonPlatform declarations or the platforms behind them.
#
# ONE OPERATOR PER CLUSTER: the operator enforces this itself at startup
# (a label-matched Deployment scan that refuses to start beside a
# sibling), so the release name is FIXED — the collision is impossible to
# express instead of merely refused. More PLATFORMS need no second
# operator: one operator watches all namespaces.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave
# create_namespace false).
resource "kubernetes_namespace_v1" "planton_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The module-owned plantonplatforms.planton.ai CRD. Keyed by the CRD's OWN
# metadata.name (never a positional index) so state addresses stay stable.
#
#   - server_side_apply + force_conflicts: adopts a CRD retained by a
#     previous install's destroy (the field manager differs) — the exact
#     semantic twin of the Pulumi module's upsert provider.
#   - apply_only: "When true, Delete is a no-op" (provider source) — the
#     keep-on-uninstall mechanism; twin of the Pulumi module's
#     retainOnDelete transformation.
resource "kubectl_manifest" "crds" {
  for_each = local.crd_manifests

  yaml_body         = each.value
  server_side_apply = true
  force_conflicts   = true
  apply_only        = true
}

# The operator release.
resource "helm_release" "planton_operator" {
  name       = local.release_name
  repository = local.helm_oci_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag) and the CRD
  # lifecycle (kubectl_manifest.crds above) — the chart's own crds/
  # install never runs.
  create_namespace = false
  skip_crds        = true

  # Wait for the operator to become Available — a manager that never
  # becomes ready (the one-per-cluster startup guard refusing beside a
  # sibling operator is THE classic case) should fail THIS apply with a
  # readiness timeout, not surface later as PlantonPlatform resources that
  # mysteriously never reconcile.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Two documents, merged in order by the provider (helm -f semantics):
  # the typed rendering first, the user's escape hatch last.
  values = concat(
    [yamlencode(local.typed_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : []
  )

  depends_on = [kubernetes_namespace_v1.planton_operator, kubectl_manifest.crds]

  lifecycle {
    # FAIL LOUDLY below the schema-contract floor: charts below
    # min_chart_version ship operators whose reconcilers predate the
    # PlantonPlatform schema the staged CRD advertises — the API server
    # would ACCEPT fields the running operator silently ignores. Twin:
    # the Pulumi module's chartVersionAtLeast guard. HCL has no semver
    # function, so the three parts compare by hand.
    precondition {
      condition = (
        tonumber(split(".", local.chart_version)[0]) > tonumber(split(".", local.min_chart_version)[0])
        ) || (
        tonumber(split(".", local.chart_version)[0]) == tonumber(split(".", local.min_chart_version)[0])
        && tonumber(split(".", local.chart_version)[1]) > tonumber(split(".", local.min_chart_version)[1])
        ) || (
        tonumber(split(".", local.chart_version)[0]) == tonumber(split(".", local.min_chart_version)[0])
        && tonumber(split(".", local.chart_version)[1]) == tonumber(split(".", local.min_chart_version)[1])
        && tonumber(split(".", local.chart_version)[2]) >= tonumber(split(".", local.min_chart_version)[2])
      )
      error_message = "Chart version ${local.chart_version} predates the PlantonPlatform schema this catalog models — use ${local.min_chart_version} or newer (older operators silently ignore fields the staged CRD accepts)."
    }

    # FAIL LOUDLY when the staged CRD file did not travel with the module:
    # fileset() over a missing ../crds returns EMPTY and for_each would
    # silently plan ZERO resources — the install would then "work" against
    # whatever CRD happens to exist. One is the staged count at chart
    # 0.7.1 — re-stage ../crds and update this count together with
    # default_chart_version. Twin of the Pulumi module's count check.
    precondition {
      condition     = try(var.spec.skip_crds, false) || length(local.crd_manifests) == 1
      error_message = "The staged CRD directory carries ${length(local.crd_manifests)} CRDs, expected 1 — the module owns the CRD lifecycle and cannot install without its full staged set."
    }
  }
}
