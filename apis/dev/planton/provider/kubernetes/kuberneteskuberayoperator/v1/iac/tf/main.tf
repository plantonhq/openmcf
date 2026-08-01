# KubernetesKubeRayOperator Terraform module.
#
# Installs the KubeRay operator from the official kuberay-operator Helm
# chart as a single Helm release named after metadata.name. The operator
# reconciles RayCluster (declared with KubernetesRayCluster), RayJob and
# RayService custom resources into running Ray clusters.
#
# CRD LIFECYCLE: the chart ships its three ray.io CRDs (rayclusters,
# rayjobs, rayservices) from its crds/ DIRECTORY — Helm installs them
# once, never upgrades them on chart upgrades, and LEAVES them (and every
# Ray declaration) on uninstall (no release ownership metadata). That
# upstream posture is exactly the keep-on-uninstall this catalog wants for
# workload-bearing CRDs, so the module neither re-owns nor templates them
# — chart-version bumps that change CRDs are applied manually per the
# upstream release notes. NOTE the CRDs are large (~1MB each) and install
# server-side.
#
# NO WEBHOOK: the operator validates in its reconcile loop. There is no
# admission webhook, no certificate machinery, and no cert-manager
# dependency — a bad RayCluster surfaces on the CR's status conditions,
# not as an admission rejection.
#
# MULTI-INSTANCE SAFETY (the name pins): the chart hardcodes nameOverride
# and fullnameOverride to "kuberay-operator" in its values — a second
# install anywhere on the cluster would collapse onto the same child
# names by construction. locals.typed_values pins both to metadata.name
# instead, so instances stay distinguishable.
#
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.

# The optional installation namespace. Created before the release; deleted
# with the resource (pre-existing-namespace installs leave create_namespace
# false). Watch namespaces (spec.watch_namespaces) must already exist —
# they are a watch scope, deliberately not this module's resources.
resource "kubernetes_namespace_v1" "kuberay_operator" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# The operator release.
resource "helm_release" "kuberay_operator" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the operator Deployment to become Available — an unpullable
  # image, a missing ServiceMonitor CRD, or a broken config should fail
  # THIS apply with a readiness timeout, not surface later as RayClusters
  # that mysteriously never reconcile.
  wait            = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the user's escape hatch second. This chart
  # needs no re-pin document: it has no release-owned CRDs and no
  # webhook machinery whose keys an escape-hatch value could weaponize —
  # an operator deliberately overriding the name pins owns the collision
  # consciously.
  values = concat(
    [yamlencode(local.typed_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : []
  )

  depends_on = [
    kubernetes_namespace_v1.kuberay_operator,
  ]

  lifecycle {
    precondition {
      # 63-char Kubernetes name limit minus the chart's longest derived
      # suffix, "-leader-election" (16 chars, the leader-election
      # Role/RoleBinding). The module pins fullnameOverride to this
      # identity, so the budget is exact. Twin of the Pulumi module's
      # fail-loud guard.
      condition     = length(var.metadata.name) <= 47
      error_message = "metadata.name must be 47 characters or fewer: the chart derives \"<name>-leader-election\" (16-char suffix) and Kubernetes caps names at 63."
    }
  }
}
