locals {
  # The workgroup's cloud name is the resource's metadata.name. Athena
  # workgroup names are create-time-immutable (changing the name replaces the
  # workgroup), which is why the name is not spec surface. Same basis as the
  # Pulumi module.
  workgroup_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  tags = {
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsAthenaWorkgroup"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Block-presence switches. The spec enables each optional feature by the
  # presence of its message, so the module's only job is translating presence
  # into the provider's nested blocks (and, where AWS requires an explicit
  # `enabled` flag inside the block, asserting it -- an absent block IS the
  # disabled state, per the provider's own diff-suppression semantics).
  has_result_config      = var.spec.result_configuration != null
  has_encryption         = local.has_result_config && var.spec.result_configuration.encryption_option != ""
  has_acl                = local.has_result_config && var.spec.result_configuration.s3_acl_option != ""
  has_engine_version     = var.spec.selected_engine_version != ""
  has_managed_results    = var.spec.managed_query_results != null
  has_content_encryption = var.spec.customer_content_encryption_kms_key != ""
  has_identity_center    = var.spec.identity_center != null
  has_access_grants      = var.spec.s3_access_grants != null
  has_monitoring         = var.spec.monitoring != null
  has_cw_logging         = local.has_monitoring && var.spec.monitoring.cloud_watch_logging != null
  has_managed_logging    = local.has_monitoring && var.spec.monitoring.managed_logging != null
  has_s3_logging         = local.has_monitoring && var.spec.monitoring.s3_logging != null

  # The spec defaults state to ENABLED; an omitted value must deploy the same
  # workgroup an explicit ENABLED would.
  workgroup_state = var.spec.state != null && var.spec.state != "" ? var.spec.state : "ENABLED"
}
