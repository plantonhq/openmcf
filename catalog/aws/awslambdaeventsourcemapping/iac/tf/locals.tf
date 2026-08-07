locals {
  # metadata.name is the Planton identity for this mapping node -- AWS assigns
  # the runtime UUID separately (exported as uuid).
  mapping_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsLambdaEventSourceMapping"
    "planton.ai/resource-id"   = var.metadata.id
  }

  is_self_managed_kafka = var.spec.self_managed_kafka != null
  is_msk_source         = !local.is_self_managed_kafka && length(var.spec.topics) > 0

  # Kafka consumer-group / schema-registry settings land in different provider
  # blocks depending on the source family.
  has_kafka_source_config = var.spec.kafka_consumer_group_id != "" || var.spec.schema_registry != null
}
