# KubernetesAltinityOperator Terraform module.
#
# Installs the Altinity ClickHouse operator from the official
# altinity-clickhouse-operator Helm chart as a single Helm release named
# after metadata.name. The operator reconciles ClickHouseInstallation /
# ClickHouseKeeperInstallation custom resources (declared through
# KubernetesClickHouse) into running ClickHouse clusters.
#
# CRD LIFECYCLE: the CHART owns the CRDs — they ship in its crds/
# directory (Helm-native handling: installed on first install, NEVER
# deleted on uninstall, so removing the operator never cascade-deletes
# ClickHouseInstallation resources or their data), and the chart's
# pre-install/pre-upgrade hook job (crdHook, enabled by default)
# server-side-applies them on every install and upgrade so chart upgrades
# carry CRD schema changes. The module therefore vendors NOTHING and
# leaves skip_crds false — unlike sibling operator modules whose charts
# template CRDs release-owned.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave create_namespace
# false).
resource "kubernetes_namespace_v1" "altinity_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The operator release.
resource "helm_release" "altinity_operator" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Chart-native CRD handling stays on: skip_crds true would leave a
  # first install without the four CRDs registered (the crdHook job only
  # covers install/upgrade hooks, not the crds/ directory).
  skip_crds = false

  # Wait for the operator to become Available — an operator that never
  # becomes ready (an unpullable image from a private mirror is the
  # classic case) should fail THIS apply with a readiness timeout, not
  # surface later as ClickHouse clusters that mysteriously never
  # reconcile.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second — and
  # fullnameOverride re-pinned LAST, the one deliberate exception to the
  # escape hatch's last-word contract (twin of the Pulumi module). The
  # Deployment / credentials Secret / metrics Service names — and the
  # exported outputs built from them — all derive from the fullname;
  # letting an override move it would break the naming budget and every
  # output.
  values = concat(
    [yamlencode(local.typed_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : [],
    [yamlencode({ fullnameOverride = local.release_name })]
  )

  depends_on = [
    kubernetes_namespace_v1.altinity_operator,
  ]
}
