# A CloudTrail Lake event data store: a queryable, immutable store of
# AWS activity events with its own retention and billing - no trail
# and no S3 bucket involved.
#
# Lifecycle facts the render below depends on:
#   - deletion is REFUSED while termination protection is on (AWS
#     behavior, not a module choice) - the teardown is two steps:
#     apply with termination_protection_enabled = false, then destroy;
#   - a destroyed store lingers in PENDING_DELETION for 7 days and
#     its name stays reserved until the purge completes;
#   - "suspend" is write-only at AWS (never reported back), so it is
#     asserted on every apply and invisible to imports;
#   - an omitted selector list makes AWS materialize a default
#     all-management-events selector - the first import after an
#     omitted-selector create shows that server-side default.

resource "aws_cloudtrail_event_data_store" "this" {
  # metadata.name is the store name on both engines (AWS: 3-128
  # characters).
  name = var.metadata.name

  # Rendered only on an explicit choice so the module never fights the
  # provider defaults (billing_mode EXTENDABLE_RETENTION_PRICING,
  # multi_region true, retention 2555, termination protection true).
  billing_mode                   = var.spec.billing_mode != "" ? var.spec.billing_mode : null
  kms_key_id                     = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null
  multi_region_enabled           = var.spec.multi_region_enabled
  organization_enabled           = var.spec.organization_enabled
  retention_period               = var.spec.retention_period_days != 0 ? var.spec.retention_period_days : null
  termination_protection_enabled = var.spec.termination_protection_enabled

  # The provider models suspend as a nullable string ("true"/"false").
  suspend = var.spec.suspend != null ? tostring(var.spec.suspend) : null

  # Ingestion scope: fine-grained field matching. AWS requires every
  # selector to carry an eventCategory condition (server-side rule).
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

  tags = local.aws_tags
}
