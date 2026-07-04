# Function-scoped settings AWS models as standalone resources but that
# are honestly part of the function's own configuration -- each is
# keyed by the function (replace-on-change), owned by exactly one
# function, and referenced by nothing else. Event source mappings are
# deliberately NOT here: a mapping has independent AWS identity and
# wires OTHER resources (queues, streams) into the function, so it is
# its own first-class kind.

# Named pointers to published versions -- the stable invocation targets
# clients reference. Materialized per-name so list edits update in
# place; repointing an alias ships or rolls back without touching
# callers.
resource "aws_lambda_alias" "this" {
  for_each = local.aliases

  name             = each.value.name
  description      = each.value.description
  function_name    = aws_lambda_function.this.function_name
  function_version = each.value.function_version

  # Canary routing: at most one additional version and its traffic
  # fraction (AWS's own constraint).
  dynamic "routing_config" {
    for_each = length(coalesce(each.value.routing_additional_version_weights, {})) > 0 ? [1] : []
    content {
      additional_version_weights = each.value.routing_additional_version_weights
    }
  }
}

# Pre-warmed execution environments, keyed by the alias qualifier --
# eliminates cold starts on the alias at the cost of paying for idle
# warmth.
resource "aws_lambda_provisioned_concurrency_config" "this" {
  for_each = local.provisioned_aliases

  function_name                     = aws_lambda_function.this.function_name
  qualifier                         = aws_lambda_alias.this[each.key].name
  provisioned_concurrent_executions = each.value.provisioned_concurrent_executions
}

# The built-in HTTPS endpoint. One per function ($LATEST-qualified);
# with authorization_type NONE, AWS additionally requires a public
# invoke permission, which the provider manages on this resource.
resource "aws_lambda_function_url" "this" {
  count = var.spec.function_url != null ? 1 : 0

  function_name      = aws_lambda_function.this.function_name
  authorization_type = var.spec.function_url.authorization_type
  invoke_mode        = var.spec.function_url.invoke_mode != "" ? var.spec.function_url.invoke_mode : null

  dynamic "cors" {
    for_each = var.spec.function_url.cors != null ? [1] : []
    content {
      allow_credentials = var.spec.function_url.cors.allow_credentials
      allow_origins     = length(coalesce(var.spec.function_url.cors.allow_origins, [])) > 0 ? var.spec.function_url.cors.allow_origins : null
      allow_methods     = length(coalesce(var.spec.function_url.cors.allow_methods, [])) > 0 ? var.spec.function_url.cors.allow_methods : null
      allow_headers     = length(coalesce(var.spec.function_url.cors.allow_headers, [])) > 0 ? var.spec.function_url.cors.allow_headers : null
      expose_headers    = length(coalesce(var.spec.function_url.cors.expose_headers, [])) > 0 ? var.spec.function_url.cors.expose_headers : null
      max_age           = var.spec.function_url.cors.max_age_seconds != 0 ? var.spec.function_url.cors.max_age_seconds : null
    }
  }
}

# Resource-policy statements authorizing external principals and AWS
# services to invoke. Statements are create/delete-only in AWS -- any
# field change replaces the statement (harmless: a statement carries no
# state). Materialized per statement_id so list edits update in place.
resource "aws_lambda_permission" "this" {
  for_each = local.invoke_permissions

  statement_id  = each.value.statement_id
  function_name = aws_lambda_function.this.function_name
  principal     = each.value.principal

  # Empty keeps the sensible default: plain function invocation.
  action = each.value.action != "" ? each.value.action : "lambda:InvokeFunction"

  # source_arn/source_account scope service-principal grants to one
  # resource/account -- the confused-deputy guard.
  source_arn             = each.value.source_arn != "" ? each.value.source_arn : null
  source_account         = each.value.source_account != "" ? each.value.source_account : null
  principal_org_id       = each.value.principal_org_id != "" ? each.value.principal_org_id : null
  function_url_auth_type = each.value.function_url_auth_type != "" ? each.value.function_url_auth_type : null
}

# Asynchronous-invocation shaping: retries, event age, and on-success /
# on-failure destinations. One config per function ($LATEST-qualified).
resource "aws_lambda_function_event_invoke_config" "this" {
  count = var.spec.async_invoke_config != null ? 1 : 0

  function_name = aws_lambda_function.this.function_name

  maximum_retry_attempts       = var.spec.async_invoke_config.maximum_retry_attempts
  maximum_event_age_in_seconds = var.spec.async_invoke_config.maximum_event_age_seconds != 0 ? var.spec.async_invoke_config.maximum_event_age_seconds : null

  dynamic "destination_config" {
    for_each = var.spec.async_invoke_config.on_success_destination_arn != "" || var.spec.async_invoke_config.on_failure_destination_arn != "" ? [1] : []
    content {
      dynamic "on_success" {
        for_each = var.spec.async_invoke_config.on_success_destination_arn != "" ? [1] : []
        content {
          destination = var.spec.async_invoke_config.on_success_destination_arn
        }
      }
      dynamic "on_failure" {
        for_each = var.spec.async_invoke_config.on_failure_destination_arn != "" ? [1] : []
        content {
          destination = var.spec.async_invoke_config.on_failure_destination_arn
        }
      }
    }
  }
}

# Recursive-loop detection. Only materialized when opting OUT of the
# AWS default (Terminate) -- deleting the resource restores the
# default, so rendering the default would be a no-op resource.
resource "aws_lambda_function_recursion_config" "this" {
  count = var.spec.recursive_loop == "Allow" ? 1 : 0

  function_name  = aws_lambda_function.this.function_name
  recursive_loop = var.spec.recursive_loop
}

# Runtime-update management. Only materialized when configured --
# deleting the resource reverts to Auto (the AWS default).
resource "aws_lambda_runtime_management_config" "this" {
  count = var.spec.runtime_management != null ? 1 : 0

  function_name       = aws_lambda_function.this.function_name
  update_runtime_on   = var.spec.runtime_management.update_runtime_on
  runtime_version_arn = var.spec.runtime_management.runtime_version_arn != "" ? var.spec.runtime_management.runtime_version_arn : null
}
