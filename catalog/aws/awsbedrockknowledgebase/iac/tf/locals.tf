locals {
  # The knowledge base name is metadata.name -- the naming basis both
  # engines share so a manifest deploys identically on either.
  knowledge_base_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockKnowledgeBase"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # AWS's type discriminators are derived from which spec arm is set --
  # exactly one per the spec's CEL guards.
  kb_type = var.spec.vector != null ? "VECTOR" : (var.spec.managed != null ? "MANAGED" : (var.spec.kendra != null ? "KENDRA" : "SQL"))

  storage_type = var.spec.storage == null ? null : (
    var.spec.storage.opensearch_serverless != null ? "OPENSEARCH_SERVERLESS" : (
      var.spec.storage.opensearch_managed != null ? "OPENSEARCH_MANAGED_CLUSTER" : (
        var.spec.storage.s3_vectors != null ? "S3_VECTORS" : (
          var.spec.storage.rds != null ? "RDS" : (
            var.spec.storage.pinecone != null ? "PINECONE" : (
              var.spec.storage.mongodb_atlas != null ? "MONGO_DB_ATLAS" : (
                var.spec.storage.neptune_analytics != null ? "NEPTUNE_ANALYTICS" : "REDIS_ENTERPRISE_CLOUD"
  )))))))

  # Data sources keyed by their stable entry names (the for_each keys both
  # engines share and the data_source_ids output map keys).
  data_sources = { for d in var.spec.data_sources : d.name => d }
}
