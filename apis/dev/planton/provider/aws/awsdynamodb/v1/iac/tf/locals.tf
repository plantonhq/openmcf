locals {
  # The table name is metadata.name -- create-only in AWS, and the
  # basis both engines share so a manifest deploys identically on
  # either.
  table_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsDynamodb"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # The spec's key_schema (one HASH, optional RANGE -- CEL-enforced) is
  # lowered to the provider's hash_key/range_key scalars. Restore-created
  # tables inherit their key schema from the source, so the list may be
  # empty.
  table_hash_key  = try([for k in var.spec.key_schema : k.attribute_name if k.key_type == "HASH"][0], null)
  table_range_key = try([for k in var.spec.key_schema : k.attribute_name if k.key_type == "RANGE"][0], null)

  # GSI names that also get contributor insights, alongside the table.
  # Empty unless insights is enabled with opted-in indexes.
  insights_gsi_names = var.spec.contributor_insights == null ? [] : (
    var.spec.contributor_insights.enabled ? var.spec.contributor_insights.gsi_index_names : []
  )
}
