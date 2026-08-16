# An API Gateway usage plan with its API keys: stage coverage, quota,
# throttle ceilings, and the keys the plan admits.
#
# Lifecycle facts the renders below depend on:
#   - the plan's api_stages PATCH away cleanly on update, and the
#     provider detaches every stage before deleting the plan -- stage
#     coverage is freely editable;
#   - product_code cannot be set at create (AWS rejects it); the
#     provider applies it via a follow-up PATCH -- no ordering concern
#     for this module, but worth knowing when reading apply logs;
#   - a usage_plan_key is pure membership (key <-> plan); every field
#     change replaces it, which is free and instant;
#   - key VALUES are secrets: AWS generates them unless the spec pins
#     one, and this module never exports them.

resource "aws_api_gateway_usage_plan" "this" {
  # metadata.name is the naming basis on both engines.
  name         = var.metadata.name
  description  = var.spec.description != "" ? var.spec.description : null
  product_code = var.spec.product_code != "" ? var.spec.product_code : null

  # Stage coverage with optional per-method throttles.
  dynamic "api_stages" {
    for_each = var.spec.api_stages
    content {
      api_id = api_stages.value.rest_api_id
      stage  = api_stages.value.stage_name

      dynamic "throttle" {
        for_each = api_stages.value.method_throttles
        content {
          path        = throttle.value.path
          burst_limit = throttle.value.burst_limit
          rate_limit  = throttle.value.rate_limit
        }
      }
    }
  }

  dynamic "quota_settings" {
    for_each = var.spec.quota != null ? [var.spec.quota] : []
    content {
      limit  = quota_settings.value.limit
      period = quota_settings.value.period
      offset = quota_settings.value.offset
    }
  }

  dynamic "throttle_settings" {
    for_each = var.spec.throttle != null ? [var.spec.throttle] : []
    content {
      burst_limit = throttle_settings.value.burst_limit
      rate_limit  = throttle_settings.value.rate_limit
    }
  }

  tags = local.aws_tags
}

resource "aws_api_gateway_api_key" "this" {
  for_each = local.api_keys

  name = each.value.name

  # Always sent explicitly (empty when unset): the two engines'
  # providers carry different branded defaults ("Managed by
  # Terraform" / "Managed by Pulumi") -- an explicit send keeps the
  # created key identical across engines.
  description = each.value.description

  # Rendered only on an explicit choice so the module never fights AWS's
  # enabled-by-default (null lets the provider default of true apply --
  # the same send path as the Pulumi module).
  enabled = each.value.enabled

  customer_id = each.value.customer_id != "" ? each.value.customer_id : null

  # Omitted = AWS generates the value (recommended); a pinned value
  # arrives resolved from the managed-secret reference.
  value = each.value.value != "" ? each.value.value : null

  tags = local.aws_tags
}

# The membership attaching each key to this plan. Pure join resource --
# every field change replaces it, which is free and instant.
resource "aws_api_gateway_usage_plan_key" "this" {
  for_each = local.api_keys

  key_id        = aws_api_gateway_api_key.this[each.key].id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.this.id
}
