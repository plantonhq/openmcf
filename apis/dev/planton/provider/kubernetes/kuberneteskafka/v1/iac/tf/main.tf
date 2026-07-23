# KubernetesKafka Terraform module.
#
# Deploys one Strimzi-operator-managed KRaft Kafka cluster:
#
#   1. the namespace (optional, create_namespace),
#   2. the JMX exporter rules ConfigMap (optional, metrics.enabled) —
#      module-owned, referenced by the Kafka CR's metricsConfig,
#   3. one kafka.strimzi.io/v1 KafkaNodePool per spec.node_pools entry,
#   4. the kafka.strimzi.io/v1 Kafka CR itself.
#
# The CRs apply through kubectl_manifest (alekc/kubectl): unlike the
# hashicorp provider's kubernetes_manifest resource it needs no cluster
# connection at plan time — the cluster can be PLANNED before the Strimzi
# operator's CRDs exist, which is what lets an infra chart deploy the
# operator and its Kafka clusters in one run (and lets offline plan proofs
# work).
#
# No wait_for block, deliberately: cluster readiness depends on the
# operator (image pulls, KRaft quorum formation, listener provisioning)
# that is not part of applying the resources — the same
# never-block-on-a-controller posture as the sibling database kinds.
# Pulumi equivalent: CustomResource without await annotations.

resource "kubernetes_namespace_v1" "namespace" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.labels
  }
}

resource "kubernetes_config_map_v1" "kafka_metrics" {
  count = try(var.spec.metrics.enabled, false) ? 1 : 0

  metadata {
    name      = local.metrics_config_map_name
    namespace = local.namespace
    labels    = local.labels
  }

  data = {
    "kafka-metrics-config.yml" = local.kafka_metrics_rules
  }

  depends_on = [kubernetes_namespace_v1.namespace]
}

# Node pools are keyed by their OWN NAME (never a positional index) so
# state addresses survive pool-list reorderings and the import recipes can
# derive the CR name blind from the address key.
#
# Pools apply BEFORE the Kafka CR: Strimzi tolerates either order, but a
# Kafka CR with no matching pools reports a transient warning state the
# lanes would otherwise race.
resource "kubectl_manifest" "node_pools" {
  for_each = local.node_pool_manifests

  yaml_body = yamlencode(each.value)

  server_side_apply = true

  depends_on = [
    kubernetes_namespace_v1.namespace,
  ]
}

resource "kubectl_manifest" "kafka" {
  yaml_body = yamlencode(local.kafka_manifest)

  server_side_apply = true

  depends_on = [
    kubernetes_namespace_v1.namespace,
    kubernetes_config_map_v1.kafka_metrics,
    kubectl_manifest.node_pools,
  ]
}
