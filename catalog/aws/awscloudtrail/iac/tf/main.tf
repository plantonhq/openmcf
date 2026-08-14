# A CloudTrail trail: the account's API audit log delivered to S3,
# with optional CloudWatch Logs mirroring, SNS delivery notices,
# Insights anomaly detection, and organization-wide capture.
#
# Lifecycle facts the renders below depend on:
#   - AWS validates the delivery bucket's POLICY at create ("Incorrect
#     S3 bucket policy is detected") -- the bucket policy is the
#     consumer's contract (AwsS3Bucket spec.policy), never this
#     module's;
#   - the classic and advanced selector styles are mutually exclusive
#     on a trail (the spec CEL guarantees only one arrives here);
#   - AWS expects the CloudWatch group ARN in its ":*" suffix form;
#     the module appends the suffix when the referenced value lacks
#     it, so both engines send the identical ARN;
#   - the delegated-admin registration is an ACCOUNT-GLOBAL act with
#     its own lifecycle (deregistered on destroy) -- it has no
#     structural edge to the trail resource, so no depends_on.

locals {
  # AWS expects "arn:...:log-group:<name>:*" -- normalize once so a
  # bare group ARN reference renders identically on both engines.
  cloudwatch_log_group_arn = var.spec.cloudwatch_logs != null ? (
    can(regex(":\\*$", var.spec.cloudwatch_logs.log_group_arn))
    ? var.spec.cloudwatch_logs.log_group_arn
    : "${var.spec.cloudwatch_logs.log_group_arn}:*"
  ) : null
}

resource "aws_cloudtrail" "this" {
  # metadata.name is the trail name on both engines (AWS: 3-128 chars,
  # letters/digits/._-, starts and ends alphanumeric).
  name           = var.metadata.name
  s3_bucket_name = var.spec.s3_bucket_name
  s3_key_prefix  = var.spec.s3_key_prefix != "" ? var.spec.s3_key_prefix : null

  is_multi_region_trail = var.spec.is_multi_region_trail
  is_organization_trail = var.spec.is_organization_trail

  # Rendered only on an explicit choice so the module never fights the
  # provider defaults (include_global_service_events and
  # enable_logging both default true).
  include_global_service_events = var.spec.include_global_service_events
  enable_logging                = var.spec.enable_logging

  enable_log_file_validation = var.spec.enable_log_file_validation

  kms_key_id     = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null
  sns_topic_name = var.spec.sns_topic_name != "" ? var.spec.sns_topic_name : null

  cloud_watch_logs_group_arn = local.cloudwatch_log_group_arn
  cloud_watch_logs_role_arn  = var.spec.cloudwatch_logs != null ? var.spec.cloudwatch_logs.role_arn : null

  # Classic selectors: management scope plus coarse data-event scopes.
  dynamic "event_selector" {
    for_each = var.spec.event_selectors
    content {
      read_write_type                  = event_selector.value.read_write_type != "" ? event_selector.value.read_write_type : "All"
      include_management_events        = event_selector.value.include_management_events
      exclude_management_event_sources = event_selector.value.exclude_management_event_sources

      dynamic "data_resource" {
        for_each = event_selector.value.data_resources
        content {
          type   = data_resource.value.type
          values = data_resource.value.values
        }
      }
    }
  }

  # Advanced selectors: fine-grained field matching.
  dynamic "advanced_event_selector" {
    for_each = var.spec.advanced_event_selectors
    content {
      name = advanced_event_selector.value.name != "" ? advanced_event_selector.value.name : null

      dynamic "field_selector" {
        for_each = advanced_event_selector.value.field_selectors
        content {
          field           = field_selector.value.field
          equals          = length(field_selector.value.equals) > 0 ? field_selector.value.equals : null
          not_equals      = length(field_selector.value.not_equals) > 0 ? field_selector.value.not_equals : null
          starts_with     = length(field_selector.value.starts_with) > 0 ? field_selector.value.starts_with : null
          not_starts_with = length(field_selector.value.not_starts_with) > 0 ? field_selector.value.not_starts_with : null
          ends_with       = length(field_selector.value.ends_with) > 0 ? field_selector.value.ends_with : null
          not_ends_with   = length(field_selector.value.not_ends_with) > 0 ? field_selector.value.not_ends_with : null
        }
      }
    }
  }

  # Insights engines (anomaly detection; billed separately).
  dynamic "insight_selector" {
    for_each = var.spec.insight_types
    content {
      insight_type = insight_selector.value
    }
  }

  tags = local.aws_tags
}

# The organization's delegated CloudTrail administrator - an
# account-global registration (one per organization, performed from
# the management account). Reads resolve through the Organizations
# API, so the caller needs organizations:Describe* alongside
# cloudtrail:RegisterOrganizationDelegatedAdmin.
resource "aws_cloudtrail_organization_delegated_admin_account" "this" {
  count = var.spec.organization_delegated_admin_account_id != "" ? 1 : 0

  account_id = var.spec.organization_delegated_admin_account_id
}
