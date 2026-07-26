# Computed values for the KubernetesKafkaTopic module. Every resolution
# here has an exact twin in the Pulumi module — keep them in lockstep.
#
# HCL DISCIPLINE: conditional entries are `key = cond ? value : null`
# inside one object literal, pruned with a for-expression (preserves value
# types); optional nested reads go through try().

locals {
  namespace = var.spec.namespace

  # The actual Kafka topic name: the override when set (dots, underscores,
  # uppercase — names Kubernetes metadata cannot carry), metadata.name
  # otherwise. Twin of the Pulumi module's fallback.
  topic_name = coalesce(try(var.spec.topic_name, null) != null && var.spec.topic_name != "" ? var.spec.topic_name : null, var.metadata.name)

  # Resource-identity labels plus the Strimzi binding label — without
  # strimzi.io/cluster the topic operator never picks the resource up.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKafkaTopic"
      "strimzi.io/cluster"       = var.spec.kafka_cluster
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  topic_manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaTopic"
    metadata = {
      name      = var.metadata.name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = {
      for k, v in {
        topicName  = try(coalesce(var.spec.topic_name), "") != "" ? var.spec.topic_name : null
        partitions = try(var.spec.partitions, null)
        replicas   = try(var.spec.replicas, null)
        config     = length(try(var.spec.config, {})) > 0 ? var.spec.config : null
      } : k => v if v != null
    }
  }
}
