# KubernetesMetricsServer Terraform module.
#
# Installs metrics-server from the official Helm chart as a real Helm
# release. The typed spec renders into chart values (locals.typed_values);
# the helm_values escape hatch is passed as a SECOND values document, which
# the provider merges over the first with Helm -f semantics — the exact
# semantic twin of the Pulumi module's buildHelmValues + mergeMaps.
#
# The release name is FIXED ("metrics-server"): the component registers the
# cluster-wide v1beta1.metrics.k8s.io APIService, a singleton — one
# installation per cluster is an upstream constraint.

# The optional installation namespace. Created before the release; deleted
# with the resource (kube-system installs leave create_namespace false).
resource "kubernetes_namespace_v1" "metrics_server" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "metrics_server" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the Deployment to become Available — a metrics-server that
  # never becomes ready (kubelet TLS rejection is THE classic cause:
  # self-signed kubelets without kubelet_insecure_tls) should fail THIS
  # deploy with a readiness timeout, not surface later as HPAs that
  # mysteriously never scale. The readinessProbe (/readyz) only passes once
  # the first kubelet scrape succeeds, so a green deploy means metrics are
  # actually flowing.
  wait            = true
  wait_for_jobs   = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 300

  # Two documents, merged in order by the provider (helm -f semantics):
  # the typed rendering first, the user's escape hatch last.
  values = concat(
    [yamlencode(local.typed_values)],
    var.spec.helm_values != "" ? [var.spec.helm_values] : []
  )

  depends_on = [kubernetes_namespace_v1.metrics_server]
}
