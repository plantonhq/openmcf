# KubernetesCertManager Terraform module.
#
# Installs cert-manager from the official Helm chart as a real Helm release.
# The typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which the
# provider merges over the first with Helm -f semantics — the exact semantic
# twin of the Pulumi module's buildHelmValues + mergeMaps.
#
# The chart owns ALL of cert-manager's Kubernetes objects, including the
# controller ServiceAccount (serviceAccount.create stays true; the workload
# identity annotation rides serviceAccount.annotations). The module itself
# creates only the optional anchor namespace.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "cert_manager" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "cert_manager" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the whole install to be ready — including the startupapicheck
  # hook Job that proves the webhook actually serves. A cert-manager whose
  # webhook is not ready rejects every Issuer/Certificate apply, so a
  # premature "success" would just move the failure downstream.
  wait            = true
  wait_for_jobs   = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Two documents, merged in order by the provider (helm -f semantics):
  # the typed rendering first, the user's escape hatch last.
  values = concat(
    [yamlencode(local.typed_values)],
    var.spec.helm_values != "" ? [var.spec.helm_values] : []
  )

  depends_on = [kubernetes_namespace_v1.cert_manager]
}
