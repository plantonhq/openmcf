# KubernetesVelero Terraform module.
#
# Installs Velero from the official Helm chart as a real Helm release. The
# typed spec renders into chart values (locals.typed_values); credential
# material rides a SECOND, sensitive()-wrapped values document; the
# helm_values escape hatch is passed as the LAST document. The provider
# merges the documents in order with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps (which
# renders credentials inline and masks the whole map when secrets are
# present).
#
# The chart owns ALL of Velero's Kubernetes objects, including the server
# ServiceAccount (serviceAccount.server.create stays true; keyless
# identity annotations ride serviceAccount.server.annotations) and the
# credentials Secret (credentials.secretContents). The module itself
# creates only the optional anchor namespace.

# The optional installation namespace. Created before the release;
# deleted with the resource.
resource "kubernetes_namespace_v1" "velero" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "velero" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the server — and the node-agent DaemonSet / upgradeCRDs job
  # when enabled — to become ready. A Velero that never comes up (a
  # ServiceMonitor rendered without the Prometheus operator CRDs, an
  # unpullable plugin image) should fail THIS apply with a readiness
  # timeout, not surface later as backups that silently never run. 600s
  # covers the CRD-upgrade job plus plugin pulls on slow registries.
  wait            = true
  wait_for_jobs   = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Documents merged in order by the provider (helm -f semantics): the
  # typed rendering first, the isolated credential document next (the
  # ONLY place secret material appears — sensitive() masks it in plans
  # and state), the user's escape hatch last.
  values = concat(
    [yamlencode(local.typed_values)],
    local.credentials_values_doc != null ? [sensitive(local.credentials_values_doc)] : [],
    var.spec.helm_values != "" ? [var.spec.helm_values] : []
  )

  depends_on = [kubernetes_namespace_v1.velero]
}
