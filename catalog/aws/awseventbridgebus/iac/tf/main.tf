resource "aws_cloudwatch_event_bus" "this" {
  name = local.resource_name

  # Partner event source (SaaS integrations) — when set, AWS requires the bus
  # name to match the source name exactly; both are create-time immutable.
  event_source_name = local.event_source_name

  description = var.spec.description != "" ? var.spec.description : null

  # Events are encrypted with an AWS-owned key unless a customer-managed key
  # is supplied.
  kms_key_identifier = var.spec.kms_key_identifier != "" ? var.spec.kms_key_identifier : null

  # Bus-level DLQ: catches events that cannot be delivered to ANY target on
  # any rule attached to this bus (rules carry their own per-target DLQs).
  dynamic "dead_letter_config" {
    for_each = var.spec.dead_letter_config != null ? [var.spec.dead_letter_config] : []
    content {
      arn = dead_letter_config.value.arn
    }
  }

  # Event delivery logging to CloudWatch Logs.
  dynamic "log_config" {
    for_each = var.spec.log_config != null ? [var.spec.log_config] : []
    content {
      level          = log_config.value.level
      include_detail = log_config.value.include_detail != "" ? log_config.value.include_detail : null
    }
  }

  tags = local.aws_tags

  lifecycle {
    # "default" names the account's built-in bus; AWS rejects creating
    # another. Failing at plan time beats an opaque CreateEventBus error.
    precondition {
      condition     = local.resource_name != "default"
      error_message = "The bus name (metadata.name) must not be \"default\" — every AWS account already has a default event bus; rules can target it via event_bus_name without creating a bus resource."
    }
  }
}

# The bus's resource-based policy is a single-per-bus setting (cross-account /
# cross-org PutEvents grants) keyed by the bus name. AWS models it as its own
# PutPermission API call, so it materializes as its own resource; deleting it
# removes all cross-account grants.
resource "aws_cloudwatch_event_bus_policy" "this" {
  count = local.resource_policy != null ? 1 : 0

  event_bus_name = aws_cloudwatch_event_bus.this.name
  policy         = local.resource_policy
}

# Event archives recorded from THIS bus — one archive per spec entry, keyed by
# the archive's own name (names are identity: reordering entries is a no-op
# and a rename is honestly a new archive). Replay is an on-demand EventBridge
# operation (StartReplay), never declarative configuration.
#
# The entries arrive untyped (the spec's event_pattern Struct keeps the list
# out of Terraform's unifiable object types), and tfvars omits zero-value
# fields entirely, so every optional field reads through try() with the
# omit-to-AWS-default fallback.
resource "aws_cloudwatch_event_archive" "this" {
  for_each = { for archive in var.spec.archives : archive.name => archive }

  name             = each.value.name
  event_source_arn = aws_cloudwatch_event_bus.this.arn

  description = try(each.value.description, "") != "" ? each.value.description : null

  # 0 / absent = retain indefinitely (the AWS default).
  retention_days = try(each.value.retention_days, 0) != 0 ? each.value.retention_days : null

  # Absent = archive every event delivered to the bus.
  event_pattern = try(each.value.event_pattern, null) != null ? jsonencode(each.value.event_pattern) : null

  # Absent = AWS-owned key. This is the ARCHIVE's key, independent of the
  # bus's own kms_key_identifier.
  kms_key_identifier = try(each.value.kms_key_identifier, "") != "" ? each.value.kms_key_identifier : null
}
