# KubernetesGatekeeper Terraform module.
#
# Installs OPA Gatekeeper from the official chart as a real Helm release.
# The typed spec renders into chart values (locals.typed_helm_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.
#
# WEBHOOK LIFECYCLE: the chart OWNS the webhook configurations as release
# objects (unlike engines that register webhooks at runtime) — uninstall
# removes them with everything else. The policy webhook is fail-open by
# default; the namespace-label check webhook is fail-closed (both typed).
#
# CRD LIFECYCLE: the engine CRDs ship in the chart's crds/ directory —
# Helm installs them once and NEVER upgrades or deletes them; the chart's
# own pre-install/pre-upgrade Job (upgradeCRDs) keeps them current on
# upgrades, and uninstall KEEPS them by design. Constraint-template CRDs
# Gatekeeper creates at runtime also survive uninstall until their
# templates are deleted.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "gatekeeper" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "gatekeeper" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the engine to become Ready — the chart's own post-install
  # probe job curls the webhook endpoint, so a green install means the
  # webhook actually serves.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second. No
  # fullnameOverride re-pin here — gatekeeper's chart hardcodes its
  # resource names (no fullname derivation exists to protect).
  values = concat(
    [yamlencode(local.typed_helm_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : []
  )

  depends_on = [
    kubernetes_namespace_v1.gatekeeper,
  ]
}
