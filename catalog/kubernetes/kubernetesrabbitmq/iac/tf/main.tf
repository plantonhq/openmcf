# KubernetesRabbitMq Terraform module.
#
# Deploys one operator-managed RabbitMQ cluster:
#
#   1. the namespace (optional, create_namespace),
#   2. the rabbitmq.com/v1beta1 RabbitmqCluster CR itself — the
#      StatefulSet, the client and headless Services, the generated
#      credentials Secret (`<name>-default-user`) and the erlang-cookie
#      Secret are all operator-created from it. No ingress resources —
#      exposure composes from the client Service's type/annotations or
#      first-class kinds referencing the exported handles.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — a cluster can be PLANNED before the RabbitMQ
# Cluster Operator's CRD exists, which is what lets an infra chart deploy
# the operator and its clusters in one run (and lets offline plan proofs
# work).
#
# No wait_for block, deliberately: cluster readiness depends on the
# operator (image pulls, peer discovery, feature-flag sync) that is not
# part of applying the resource — the same never-block-on-a-controller
# posture as the sibling operator-CR modules.
#
# NAMING CONTRACT (operator source at the pinned release): the client
# Service is `<name>`, the headless Service `<name>-nodes`, the generated
# credentials Secret `<name>-default-user`, the StatefulSet
# `<name>-server` and each pod's PVC `persistence-<name>-server-<i>`.

# The optional namespace. Created before everything else (the CR is
# namespaced); deleted with the resource. Pre-existing-namespace
# deployments leave create_namespace false.
resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

# ---- the RabbitmqCluster ----------------------------------------------------
resource "kubectl_manifest" "rabbitmq_cluster" {
  yaml_body = yamlencode(local.rabbitmq_cluster_manifest)

  server_side_apply = true

  # BACKGROUND deletion, explicitly: the operator's own deletion
  # finalizer is the cascade, and foreground propagation deadlocks
  # against operators that keep reconciling children during deletion
  # (verified live on sibling operator-owned CRs). The provider would
  # default to Foreground whenever wait is set — never rely on that
  # default here. Pulumi twin: the pulumi.com/deletionPropagationPolicy
  # annotation.
  delete_cascade = "Background"

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}
