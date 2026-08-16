# CloudWatch Logs delivery: the vended pipeline (source -> deliveries
# -> destinations) and/or the legacy cross-account Kinesis destination.
#
# Lifecycle facts the render below depends on:
#   - the vended source is per (resource, log_type): name, log_type, and
#     resource_arn all replace on change;
#   - a delivery's source and destination both replace on change; only
#     the wire-format settings update in place;
#   - AWS allows ONE delivery per (source, destination-type) pair -
#     e.g. one to S3 plus one to Firehose, never two to S3;
#   - for CloudFront sources AWS prepends
#     "AWSLogs/{account-id}/CloudFront/" to the S3 suffix path; the
#     provider strips that prefix on reads, so configure only your own
#     path segment;
#   - the cross-account destination's access policy PERSISTS at AWS
#     when only the policy resource is destroyed (a provider no-op
#     delete); destroying the destination itself is real;
#   - the cross-account destination's first create retries for up to
#     two minutes while the logs.amazonaws.com trust on the role
#     propagates.

resource "aws_cloudwatch_log_delivery_source" "this" {
  count = var.spec.vended != null && var.spec.vended.source != null ? 1 : 0

  name         = var.spec.vended.source.name
  log_type     = var.spec.vended.source.log_type
  resource_arn = var.spec.vended.source.resource_arn

  tags = local.aws_tags
}

resource "aws_cloudwatch_log_delivery_destination" "this" {
  for_each = local.destinations

  name = each.value.name

  # XRAY destinations carry no configuration block (spec-guaranteed);
  # every other type requires one.
  dynamic "delivery_destination_configuration" {
    for_each = each.value.destination_resource_arn != "" ? [1] : []
    content {
      destination_resource_arn = each.value.destination_resource_arn
    }
  }

  delivery_destination_type = each.value.delivery_destination_type != "" ? each.value.delivery_destination_type : null
  output_format             = each.value.output_format != "" ? each.value.output_format : null

  tags = local.aws_tags
}

resource "aws_cloudwatch_log_delivery_destination_policy" "this" {
  for_each = local.destination_policies

  delivery_destination_name   = aws_cloudwatch_log_delivery_destination.this[each.key].name
  delivery_destination_policy = each.value
}

resource "aws_cloudwatch_log_delivery" "this" {
  for_each = local.deliveries

  delivery_source_name = aws_cloudwatch_log_delivery_source.this[0].name
  delivery_destination_arn = (
    each.value.destination_name != ""
    ? aws_cloudwatch_log_delivery_destination.this[each.value.destination_name].arn
    : each.value.destination_arn
  )

  record_fields   = length(each.value.record_fields) > 0 ? each.value.record_fields : null
  field_delimiter = each.value.field_delimiter != "" ? each.value.field_delimiter : null

  dynamic "s3_delivery_configuration" {
    for_each = each.value.s3_configuration != null ? [each.value.s3_configuration] : []
    content {
      enable_hive_compatible_path = s3_delivery_configuration.value.enable_hive_compatible_path
      suffix_path                 = s3_delivery_configuration.value.suffix_path != "" ? s3_delivery_configuration.value.suffix_path : null
    }
  }

  tags = local.aws_tags
}

resource "aws_cloudwatch_log_destination" "this" {
  count = local.cross_account != null ? 1 : 0

  name       = local.cross_account.name
  role_arn   = local.cross_account.role_arn
  target_arn = local.cross_account.target_arn

  tags = local.aws_tags
}

resource "aws_cloudwatch_log_destination_policy" "this" {
  count = local.cross_account != null ? 1 : 0

  destination_name = aws_cloudwatch_log_destination.this[0].name
  access_policy    = jsonencode(local.cross_account.access_policy)

  force_update = local.cross_account.force_update != null ? local.cross_account.force_update : null
}
