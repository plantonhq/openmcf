locals {
  # The function name is metadata.name -- create-time immutable in AWS,
  # and the basis both engines share so a manifest deploys identically
  # on either.
  function_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsLambda"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Zip vs container image decides package_type and which code
  # arguments are sent (CEL guarantees exactly one source is set).
  is_image = var.spec.image_uri != ""

  # The log group the function writes to: the custom group from
  # logging_config, or the AWS-default "/aws/lambda/<name>" that Lambda
  # creates on first invocation. Exported so log-consuming resources
  # (subscription filters, dashboards) can compose either way.
  custom_log_group = var.spec.logging_config != null ? var.spec.logging_config.log_group : ""
  log_group_name   = local.custom_log_group != "" && local.custom_log_group != null ? local.custom_log_group : "/aws/lambda/${local.function_name}"

  # Aliases keyed by name so each materializes as its own resource and
  # list edits update in place. CEL enforces name uniqueness.
  aliases = { for a in coalesce(var.spec.aliases, []) : a.name => a }

  # Aliases that also pin provisioned concurrency (pre-warmed
  # execution environments keyed by the alias qualifier).
  provisioned_aliases = {
    for name, a in local.aliases : name => a
    if a.provisioned_concurrent_executions != null
  }

  # Invoke permissions keyed by statement_id (CEL enforces uniqueness).
  invoke_permissions = { for p in coalesce(var.spec.invoke_permissions, []) : p.statement_id => p }

  # Per-qualifier scaling configs keyed by qualifier (CEL enforces
  # uniqueness). Keys like "$LATEST.PUBLISHED" are valid map keys.
  scaling_configs = { for s in coalesce(var.spec.scaling_configs, []) : s.qualifier => s }
}
