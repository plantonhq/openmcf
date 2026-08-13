locals {
  # The state machine's cloud name is the resource's metadata.name -- the same
  # basis the Pulumi module uses. AWS restricts names to 1-80 chars of
  # [0-9A-Za-z_-], which the platform's metadata.name validation already
  # satisfies.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsStepFunction"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # State machine type -- STANDARD when unset. Changing the type replaces the
  # state machine (AWS ForceNew), which the spec documents.
  sm_type = var.spec.type != "" ? var.spec.type : "STANDARD"

  # The ASL definition arrives as a nested object (the tfvars layer passes
  # protobuf Structs through as real objects, not strings); serialize it to
  # the JSON the AWS API expects. ASL key casing survives jsonencode.
  definition = jsonencode(var.spec.definition)

  # The logging block renders for ANY configured level, including an
  # explicit OFF (the disable send -- see main.tf). Only a fully omitted
  # level sends nothing. HCL's && does NOT short-circuit, so attribute
  # access on the nullable block must ride try() -- a bare
  # var.spec.logging.level errors at apply time whenever the block is
  # omitted.
  logging_level      = try(var.spec.logging.level, "")
  logging_configured = local.logging_level != ""

  # AWS requires the CloudWatch log group ARN in log_destination to carry a
  # ":*" suffix (the provider rejects it at plan time otherwise). Referenced
  # log-group ARNs arrive without the suffix, so append it when missing --
  # users should never have to know about this quirk.
  raw_log_destination = try(var.spec.logging.log_destination, "")
  log_destination     = local.raw_log_destination != "" && !endswith(local.raw_log_destination, ":*") ? "${local.raw_log_destination}:*" : local.raw_log_destination

  # Encryption: the presence of the block selects CUSTOMER_MANAGED_KMS_KEY
  # (kms_key_id is required inside it); absence leaves AWS-owned keys, which
  # is AWS's default and carries no extra cost.
  has_encryption = try(var.spec.encryption.kms_key_id, "") != ""

  # A zero reuse period means "let AWS default" (300s); AWS accepts 60-900.
  kms_data_key_reuse_period = try(var.spec.encryption.kms_data_key_reuse_period_seconds, 0) > 0 ? var.spec.encryption.kms_data_key_reuse_period_seconds : null
}
