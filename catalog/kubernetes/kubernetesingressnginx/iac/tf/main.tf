# KubernetesIngressNginx Terraform module.
#
# Installs the ingress-nginx controller from the official Helm chart as a
# real Helm release. The typed spec renders into chart values
# (locals.typed_values); the helm_values escape hatch is passed as a SECOND
# values document, which the provider merges over the first with Helm -f
# semantics — the exact semantic twin of the Pulumi module's
# buildHelmValues + mergeMaps.
#
# The chart owns ALL of the controller's Kubernetes objects (Deployment/
# DaemonSet, Services, IngressClass, RBAC, admission webhook). The module
# itself creates only the optional anchor namespace and reads the
# controller Service back for the load-balancer address outputs.

# The optional installation namespace. Created before the release; deleted
# with the resource.
resource "kubernetes_namespace_v1" "ingress_nginx" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "ingress_nginx" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the release's resources to become ready — a controller that
  # never starts (bad image, unschedulable pod, webhook certgen failure)
  # should fail THIS deploy, not the first Ingress. Helm's readiness check
  # on a LoadBalancer-type Service also waits for the cloud LB address, so
  # on clusters WITHOUT a cloud LB controller (kind, bare metal) a
  # load_balancer service type times out loudly here — deliberate: use
  # node_port/host access on such clusters, and the failure names the real
  # problem instead of leaving a silently Pending entry point.
  # wait_for_jobs covers the admission-webhook certgen hook Jobs.
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

  depends_on = [kubernetes_namespace_v1.ingress_nginx]
}

# Read the controller Service back for the load-balancer address outputs.
# Gated on the load_balancer service type: for node_port/cluster_ip there is
# no LB status to read and the outputs stay empty by design. For
# load_balancer the helm wait above guarantees the address exists by the
# time this read runs.
data "kubernetes_service_v1" "controller" {
  count = local.is_load_balancer ? 1 : 0

  metadata {
    name      = local.controller_service_name
    namespace = local.namespace
  }

  depends_on = [helm_release.ingress_nginx]
}
