# An SSM patch baseline with its folded patch-group registrations and
# the optional account/region default-baseline designation.
#
# Lifecycle facts the render below depends on:
#   - metadata.name is the baseline name on both engines; AWS
#     identifies the baseline as "pb-..." (the import ID);
#   - operating_system is rendered only on an explicit choice (unset =
#     WINDOWS, the provider default) and forces replacement;
#   - patch groups are fully ForceNew registrations keyed by the group
#     name; a group registers with only ONE baseline per OS
#     account-wide (AWS state is the referee);
#   - the default designation is a REVERSIBLE pointer: destroying it
#     RESTORES AWS's own predefined default baseline for the OS (the
#     provider looks it up and re-registers it) - the TGW default-table
#     class, not the App Runner one-way class;
#   - if the baseline is deleted while holding the designation, the
#     provider restores AWS's default and retries the delete.

resource "aws_ssm_patch_baseline" "this" {
  # metadata.name is the baseline name on both engines.
  name = var.metadata.name

  operating_system = var.spec.operating_system != "" ? var.spec.operating_system : null
  description      = var.spec.description != "" ? var.spec.description : null

  dynamic "approval_rule" {
    for_each = var.spec.approval_rules
    content {
      approve_after_days  = approval_rule.value.approve_after_days
      approve_until_date  = approval_rule.value.approve_until_date != "" ? approval_rule.value.approve_until_date : null
      compliance_level    = approval_rule.value.compliance_level != "" ? approval_rule.value.compliance_level : null
      enable_non_security = approval_rule.value.enable_non_security

      dynamic "patch_filter" {
        for_each = approval_rule.value.patch_filters
        content {
          key    = patch_filter.value.key
          values = patch_filter.value.values
        }
      }
    }
  }

  dynamic "global_filter" {
    for_each = var.spec.global_filters
    content {
      key    = global_filter.value.key
      values = global_filter.value.values
    }
  }

  approved_patches                  = length(var.spec.approved_patches) > 0 ? var.spec.approved_patches : null
  approved_patches_compliance_level = var.spec.approved_patches_compliance_level != "" ? var.spec.approved_patches_compliance_level : null
  approved_patches_enable_non_security = var.spec.approved_patches_enable_non_security

  rejected_patches        = length(var.spec.rejected_patches) > 0 ? var.spec.rejected_patches : null
  rejected_patches_action = var.spec.rejected_patches_action != "" ? var.spec.rejected_patches_action : null

  available_security_updates_compliance_status = var.spec.available_security_updates_compliance_status != "" ? var.spec.available_security_updates_compliance_status : null

  dynamic "source" {
    for_each = var.spec.sources
    content {
      name          = source.value.name
      configuration = source.value.configuration
      products      = source.value.products
    }
  }

  tags = local.aws_tags
}

# Folded patch-group registrations, keyed by group name (fully
# ForceNew at the provider - any change replaces the registration).
resource "aws_ssm_patch_group" "this" {
  for_each = toset(var.spec.patch_groups)

  patch_group = each.value
  baseline_id = aws_ssm_patch_baseline.this.id
}

# The account/region default-baseline designation for this baseline's
# OS. Destroying it restores AWS's own predefined default (the
# provider records and reverts - a true revert, unlike the App Runner
# one-way designation).
resource "aws_ssm_default_patch_baseline" "this" {
  count = var.spec.set_as_default_baseline ? 1 : 0

  baseline_id      = aws_ssm_patch_baseline.this.id
  operating_system = aws_ssm_patch_baseline.this.operating_system
}
