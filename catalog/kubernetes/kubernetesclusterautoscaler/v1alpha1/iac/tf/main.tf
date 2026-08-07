# KubernetesClusterAutoscaler Terraform module.
#
# Installs the Kubernetes Cluster Autoscaler from the official Helm chart
# as a real Helm release. The typed spec renders into chart values
# (locals.typed_values); the helm_values escape hatch is passed as a SECOND
# values document, which the provider merges over the first with Helm -f
# semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.
#
# The release name is FIXED ("cluster-autoscaler"): the autoscaler
# leader-elects and owns the cluster-wide scaling decision — a second
# installation would fight the first over every scale-up, so one
# installation per cluster is the operating model.

# The optional installation namespace. Created before the release; deleted
# with the resource. kube-system installs (the upstream convention) leave
# create_namespace false — the namespace always pre-exists.
resource "kubernetes_namespace_v1" "cluster_autoscaler" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "cluster_autoscaler" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the Deployment to become Available — an autoscaler that never
  # becomes ready (bad cloud credentials crash-looping the pod is THE
  # classic failure; a ServiceMonitor rendered without the Prometheus
  # operator CRDs the other) should fail THIS apply with a readiness
  # timeout, not surface later as node groups that mysteriously never
  # scale. 600s because the image pull plus leader election on a busy
  # kube-system can exceed the usual 300.
  wait            = true
  wait_for_jobs   = true
  atomic          = true
  cleanup_on_fail = true
  timeout         = 600

  # Two documents, merged in order by the provider (helm -f semantics):
  # the typed rendering first, the user's escape hatch last. The typed
  # document carries the arm's credentials (AWS secret key, Azure client
  # secret, Civo API key) on their way into the chart's own Secret — never
  # log or output it.
  values = concat(
    [yamlencode(local.typed_values)],
    try(var.spec.helm_values, "") != "" ? [var.spec.helm_values] : []
  )

  depends_on = [kubernetes_namespace_v1.cluster_autoscaler]
}
