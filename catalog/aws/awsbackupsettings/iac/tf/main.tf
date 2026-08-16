# AWS Backup account/region settings: two independent singleton arms.
#
# Lifecycle facts the render below depends on:
#   - BOTH resources' deletes are no-ops at the provider - destroying
#     this component changes NOTHING at AWS; the last-applied settings
#     stay in effect (taught on the spec arms);
#   - both are Required full maps at the provider: AWS returns every
#     supported key on read, so a key/type missing from the spec map
#     shows as a perpetual plan difference - the spec teaches listing
#     every key deliberately;
#   - the global arm is ACCOUNT-wide (its AWS identity is the account
#     ID, no region involved); the region arm's identity is the
#     region.

resource "aws_backup_global_settings" "this" {
  count = var.spec.global != null ? 1 : 0

  global_settings = var.spec.global.settings
}

resource "aws_backup_region_settings" "this" {
  count = var.spec.region_settings != null ? 1 : 0

  resource_type_opt_in_preference = var.spec.region_settings.resource_type_opt_in_preference

  # Rendered only when set - once set at AWS, the preference cannot be
  # cleared back to unset, only flipped per type.
  resource_type_management_preference = length(var.spec.region_settings.resource_type_management_preference) > 0 ? var.spec.region_settings.resource_type_management_preference : null
}

# The settings resources expose no ARNs - the identities below are the
# outputs' source.
data "aws_caller_identity" "this" {}
data "aws_region" "this" {}
