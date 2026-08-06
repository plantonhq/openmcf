# KubernetesCilium Terraform module.
#
# Installs Cilium from the official Helm chart as a real Helm release. The
# typed spec renders into chart values (locals.typed_values); the
# helm_values escape hatch is passed as a SECOND values document, which the
# provider merges over the first with Helm -f semantics — the exact semantic
# twin of the Pulumi module's buildHelmValues + mergeMaps.
#
# The release name is FIXED ("cilium"): Cilium is the node dataplane — the
# agent DaemonSet, operator, and generated CNI configuration are cluster
# singletons, so one dataplane per cluster is an upstream constraint.

# The optional installation namespace. Created before the release; deleted
# with the resource (kube-system installs leave create_namespace false).
resource "kubernetes_namespace_v1" "cilium" {
  count = var.spec.create_namespace ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "helm_release" "cilium" {
  name       = local.release_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.chart_version
  namespace  = local.namespace

  # The module owns namespace creation (create_namespace flag).
  create_namespace = false

  # Wait for the whole dataplane to come up. 600s (not the default 300)
  # because the install path is heavier than an ordinary workload chart:
  # the agent DaemonSet must roll out on EVERY node plus the operator, and
  # on a fresh cluster nodes transition NotReady->Ready only as Cilium
  # wires each one — the rollout itself unblocks scheduling. A dataplane
  # that never converges should fail THIS deploy, not surface later as
  # pods stuck in ContainerCreating.
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

  depends_on = [kubernetes_namespace_v1.cilium]
}
