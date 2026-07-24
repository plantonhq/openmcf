# Computed values for the KubernetesKafkaConnector module. Every
# resolution here has an exact twin in the Pulumi module — keep them in
# lockstep.
#
# HCL DISCIPLINE: conditional entries are `key = cond ? value : null`
# inside one object literal, pruned with a for-expression (preserves value
# types); optional nested reads go through try().

locals {
  namespace = var.spec.namespace

  # Resource-identity labels plus the Strimzi binding label — without
  # strimzi.io/cluster the cluster operator never picks the connector up.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKafkaConnector"
      "strimzi.io/cluster"       = var.spec.connect_cluster
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  connector_manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaConnector"
    metadata = {
      name      = var.metadata.name
      namespace = local.namespace
      labels    = local.labels
    }
    spec = {
      for k, v in {
        class    = var.spec.connector_class
        tasksMax = try(var.spec.tasks_max, null)
        version  = try(coalesce(var.spec.version), "") != "" ? var.spec.version : null
        config   = length(try(var.spec.config, {})) > 0 ? var.spec.config : null
        state    = try(coalesce(var.spec.state), "") != "" ? var.spec.state : null

        autoRestart = try(var.spec.auto_restart, null) == null ? null : {
          for ak, av in {
            enabled     = try(var.spec.auto_restart.enabled, false)
            maxRestarts = try(var.spec.auto_restart.max_restarts, null)
          } : ak => av if av != null
        }

        # The offset ConfigMap targets are declarations only — the
        # list/alter actions run when the resource carries the
        # strimzi.io/connector-offsets annotation (an operational verb
        # outside this module's scope).
        listOffsets = try(var.spec.list_offsets, null) == null ? null : {
          toConfigMap = {
            name = var.spec.list_offsets.to_config_map
          }
        }
        alterOffsets = try(var.spec.alter_offsets, null) == null ? null : {
          fromConfigMap = {
            name = var.spec.alter_offsets.from_config_map
          }
        }
      } : k => v if v != null
    }
  }
}
