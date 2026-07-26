# Computed values for the KubernetesKafkaUser module. Every resolution
# here has an exact twin in the Pulumi module — keep them in lockstep.
#
# HCL DISCIPLINE: conditional entries are `key = cond ? value : null`
# inside one object literal, pruned with a for-expression (preserves value
# types); optional nested reads go through try().

locals {
  namespace = var.spec.namespace
  username  = var.metadata.name

  # tls-external users authenticate with certificates issued OUTSIDE the
  # cluster — the user operator generates no Secret for them (nor for
  # authentication-less ACL-only principals), so the handle is honestly
  # empty. Twin of the Pulumi module's conditional.
  secret_name = try(var.spec.authentication.type, "") != "" && try(var.spec.authentication.type, "") != "tls-external" ? var.metadata.name : ""

  # Resource-identity labels plus the Strimzi binding label — without
  # strimzi.io/cluster the user operator never picks the resource up.
  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesKafkaUser"
      "strimzi.io/cluster"       = var.spec.kafka_cluster
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  quotas_body = try(var.spec.quotas, null) == null ? null : (
    length({
      for k, v in {
        producerByteRate       = try(var.spec.quotas.producer_byte_rate, null)
        consumerByteRate       = try(var.spec.quotas.consumer_byte_rate, null)
        requestPercentage      = try(var.spec.quotas.request_percentage, null)
        controllerMutationRate = try(var.spec.quotas.controller_mutation_rate, null)
      } : k => v if v != null
    }) > 0 ? {
      for k, v in {
        producerByteRate       = try(var.spec.quotas.producer_byte_rate, null)
        consumerByteRate       = try(var.spec.quotas.consumer_byte_rate, null)
        requestPercentage      = try(var.spec.quotas.request_percentage, null)
        controllerMutationRate = try(var.spec.quotas.controller_mutation_rate, null)
      } : k => v if v != null
    } : null
  )

  user_manifest = {
    apiVersion = "kafka.strimzi.io/v1"
    kind       = "KafkaUser"
    metadata = {
      name      = local.username
      namespace = local.namespace
      labels    = local.labels
    }
    spec = {
      for k, v in {
        authentication = try(var.spec.authentication, null) == null ? null : {
          type = var.spec.authentication.type
        }

        authorization = try(var.spec.authorization, null) == null ? null : {
          type = coalesce(try(var.spec.authorization.type, null), "simple")
          acls = [
            for acl in var.spec.authorization.acls : {
              for ak, av in {
                resource = {
                  for rk, rv in {
                    type        = acl.resource.type
                    name        = try(coalesce(acl.resource.name), "") != "" ? acl.resource.name : null
                    patternType = try(coalesce(acl.resource.pattern_type), "") != "" ? acl.resource.pattern_type : null
                  } : rk => rv if rv != null
                }
                operations = acl.operations
                host       = try(coalesce(acl.host), "") != "" ? acl.host : null
              } : ak => av if av != null
            }
          ]
        }

        quotas = local.quotas_body
      } : k => v if v != null
    }
  }
}
