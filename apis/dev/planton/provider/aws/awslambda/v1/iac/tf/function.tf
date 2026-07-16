# The Lambda function itself. Create-time immutable in AWS: the
# function name and the package type (zip vs container image).
# Everything else -- code, memory, timeout, VPC attachment, layers,
# logging -- edits in place (code changes roll the function without
# replacing it).
resource "aws_lambda_function" "this" {
  function_name = local.function_name
  role          = var.spec.role_arn
  description   = var.spec.description != "" ? var.spec.description : null
  tags          = local.aws_tags

  # Zip vs container image is a create-time fork (CEL guarantees
  # exactly one source): package_type cannot change on a live function.
  package_type = local.is_image ? "Image" : "Zip"

  image_uri         = local.is_image ? var.spec.image_uri : null
  s3_bucket         = local.is_image ? null : var.spec.s3.bucket
  s3_key            = local.is_image ? null : var.spec.s3.key
  s3_object_version = local.is_image || var.spec.s3 == null ? null : (var.spec.s3.object_version != "" ? var.spec.s3.object_version : null)

  # The hash is what makes code rolls declarative: a new value rolls
  # the function even when the S3 key is rewritten in place; an
  # unchanged value is a no-op.
  source_code_hash = var.spec.source_code_hash != "" ? var.spec.source_code_hash : null

  # BYOK for the deployment package itself -- distinct from
  # kms_key_arn, which encrypts environment variables.
  source_kms_key_arn = var.spec.source_kms_key_arn != "" ? var.spec.source_kms_key_arn : null

  # Runtime/handler drive zip execution; images carry their own (CEL
  # keeps them empty in image mode).
  runtime = var.spec.runtime != "" ? var.spec.runtime : null
  handler = var.spec.handler != "" ? var.spec.handler : null

  # The provider models architecture as a single-element list. Empty
  # keeps the AWS default (x86_64).
  architectures = var.spec.architecture != "" ? [var.spec.architecture] : null

  # 0 keeps the AWS defaults (128 MB / 3 s) -- the provider then sends
  # its own defaults, which match.
  memory_size = var.spec.memory_size_mb != 0 ? var.spec.memory_size_mb : null
  timeout     = var.spec.timeout_seconds != 0 ? var.spec.timeout_seconds : null

  # null in the spec means "draw from the unreserved account pool",
  # which the provider expresses as -1; 0 is the explicit kill switch.
  reserved_concurrent_executions = var.spec.reserved_concurrent_executions != null ? var.spec.reserved_concurrent_executions : -1

  publish = var.spec.publish

  kms_key_arn = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null

  code_signing_config_arn = var.spec.code_signing_config_arn != "" ? var.spec.code_signing_config_arn : null

  layers = length(coalesce(var.spec.layer_arns, [])) > 0 ? var.spec.layer_arns : null

  dynamic "ephemeral_storage" {
    for_each = var.spec.ephemeral_storage_mb != 0 ? [1] : []
    content {
      size = var.spec.ephemeral_storage_mb
    }
  }

  dynamic "environment" {
    for_each = length(coalesce(var.spec.environment, {})) > 0 ? [1] : []
    content {
      variables = var.spec.environment
    }
  }

  # VPC attachment travels as a set (CEL): subnets + security groups
  # together, IPv6 egress only on top of them.
  dynamic "vpc_config" {
    for_each = length(coalesce(var.spec.subnet_ids, [])) > 0 ? [1] : []
    content {
      subnet_ids                  = var.spec.subnet_ids
      security_group_ids          = var.spec.security_group_ids
      ipv6_allowed_for_dual_stack = var.spec.ipv6_allowed_for_dual_stack
    }
  }

  dynamic "dead_letter_config" {
    for_each = var.spec.dead_letter_target_arn != "" ? [1] : []
    content {
      target_arn = var.spec.dead_letter_target_arn
    }
  }

  dynamic "tracing_config" {
    for_each = var.spec.tracing_mode != "" ? [1] : []
    content {
      mode = var.spec.tracing_mode
    }
  }

  # EFS mounts require the VPC attachment to reach the file system's
  # mount targets (CEL enforces the coupling).
  dynamic "file_system_config" {
    for_each = var.spec.file_system_config != null ? [1] : []
    content {
      arn              = var.spec.file_system_config.access_point_arn
      local_mount_path = var.spec.file_system_config.local_mount_path
    }
  }

  dynamic "image_config" {
    for_each = var.spec.image_config != null ? [1] : []
    content {
      entry_point       = length(coalesce(var.spec.image_config.entry_point, [])) > 0 ? var.spec.image_config.entry_point : null
      command           = length(coalesce(var.spec.image_config.command, [])) > 0 ? var.spec.image_config.command : null
      working_directory = var.spec.image_config.working_directory != "" ? var.spec.image_config.working_directory : null
    }
  }

  # SnapStart snapshots published versions only (CEL couples it to
  # publish) -- "None" is the provider's off state, so the block is
  # only rendered when enabled.
  dynamic "snap_start" {
    for_each = var.spec.snap_start ? [1] : []
    content {
      apply_on = "PublishedVersions"
    }
  }

  dynamic "logging_config" {
    for_each = var.spec.logging_config != null ? [1] : []
    content {
      # The provider requires log_format whenever the block is present;
      # "Text" is the AWS default.
      log_format            = var.spec.logging_config.log_format != "" ? var.spec.logging_config.log_format : "Text"
      application_log_level = var.spec.logging_config.application_log_level != "" ? var.spec.logging_config.application_log_level : null
      system_log_level      = var.spec.logging_config.system_log_level != "" ? var.spec.logging_config.system_log_level : null
      log_group             = var.spec.logging_config.log_group != "" ? var.spec.logging_config.log_group : null
    }
  }
}
