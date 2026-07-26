# KubernetesKafkaMirrorMaker2 Terraform module.
#
# Deploys one Strimzi-operator-managed MirrorMaker 2 replication engine:
#
#   1. the namespace (optional, create_namespace),
#   2. the JMX exporter rules ConfigMap (optional, metrics.enabled) —
#      module-owned, referenced by the CR's metricsConfig,
#   3. the kafka.strimzi.io/v1 KafkaMirrorMaker2 CR itself.
#
# The CR applies through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — the engine can be PLANNED before the Strimzi
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its Kafka workloads in one run (and lets offline plan
# proofs work).
#
# No wait_for block, deliberately: engine readiness depends on the
# operator (image pulls, worker group formation, connector startup) that
# is not part of applying the resources — the same
# never-block-on-a-controller posture as the sibling operator-CR kinds.
# Pulumi equivalent: CustomResource without await annotations.

resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "kubernetes_config_map_v1" "mirrormaker2_metrics" {
  count = try(var.spec.metrics.enabled, false) ? 1 : 0

  metadata {
    name      = local.metrics_config_map_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    "metrics-config.yml" = local.mirrormaker2_metrics_rules
  }

  depends_on = [kubernetes_namespace_v1.namespace]
}

resource "kubectl_manifest" "mirrormaker2" {
  yaml_body = yamlencode(local.mirrormaker2_manifest)

  server_side_apply = true

  depends_on = [
    kubernetes_namespace_v1.namespace,
    kubernetes_config_map_v1.mirrormaker2_metrics,
  ]
}
