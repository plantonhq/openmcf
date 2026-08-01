# KubernetesSparkOperator Terraform module.
#
# Installs the Apache Spark Kubernetes Operator from the official ASF
# spark-kubernetes-operator Helm chart as a single Helm release named
# after metadata.name. The operator reconciles SparkApplication (per-job)
# and SparkCluster (long-lived) custom resources into running Spark
# workloads.
#
# CRD LIFECYCLE: the chart ships its two spark.apache.org CRDs from its
# crds/ DIRECTORY — Helm installs them once, never upgrades them on chart
# upgrades, and LEAVES them on uninstall (no release ownership metadata).
# That upstream posture is exactly the keep-on-uninstall this catalog
# wants for workload-bearing CRDs, so the module neither re-owns nor
# templates them — chart-version bumps that change CRDs are applied
# manually per the upstream release notes.
#
# NO WEBHOOK: this operator validates in its reconcile loop. There is no
# admission webhook, no certificate machinery, and no cert-manager
# dependency — one less lifecycle to manage and nothing that can
# fail-close the cluster's write path.
#
# MULTI-INSTANCE SAFETY (the RBAC name re-pins): the chart hardcodes all
# its cluster-scoped RBAC names as plain values ("spark-operator-
# clusterrole", …) — a second install anywhere on the cluster would
# collide by construction. locals.typed_values derives every RBAC name
# from the release identity instead, so instances coexist.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave create_namespace
# false). Workload namespaces (spec.workload.namespaces) are CHART-created
# and chart-kept (helm.sh/resource-policy: keep) — deliberately not this
# module's resources.
resource "kubernetes_namespace_v1" "spark_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The operator release.
resource "helm_release" "spark_operator" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the operator Deployment to become Available — a JVM with a
  # 30s-initial-delay startup probe; an unpullable image or a broken
  # config should fail THIS apply with a readiness timeout, not surface
  # later as SparkApplications that mysteriously never reconcile.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second. This chart
  # needs no re-pin document: it has no release-owned CRDs and no
  # webhook machinery whose keys an escape-hatch value could weaponize —
  # the RBAC name re-pins live in the FIRST document and an operator
  # deliberately overriding them owns the collision consciously.
  values = concat(
    [yamlencode(local.typed_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : []
  )

  depends_on = [
    kubernetes_namespace_v1.spark_operator,
  ]

  lifecycle {
    precondition {
      # 63-char Kubernetes name limit minus the module's longest derived
      # suffix, "-workload-clusterrole"/"-config-monitor-binding"
      # (23 chars). The module pins fullnameOverride AND every RBAC name
      # to this identity, so the budget is exact. Twin of the Pulumi
      # module's fail-loud guard.
      condition     = length(var.metadata.name) <= 40
      error_message = "metadata.name must be 40 characters or fewer: the module derives \"<name>-config-monitor-binding\" (23-char suffix) and Kubernetes caps names at 63."
    }
  }
}
