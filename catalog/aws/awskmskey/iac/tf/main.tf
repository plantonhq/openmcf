# A customer-managed KMS key: cryptographic shape, key policy, rotation,
# multi-region designation, and friendly aliases. key_spec and key_usage are
# create-time immutable; rotation, disabled state, policy, and aliases edit
# in place.
resource "aws_kms_key" "this" {
  description = var.spec.description != "" ? var.spec.description : null

  customer_master_key_spec = var.spec.key_spec != "" ? var.spec.key_spec : null
  key_usage                = var.spec.key_usage != "" ? var.spec.key_usage : null

  policy                             = var.spec.policy != "" ? var.spec.policy : null
  bypass_policy_lockout_safety_check = var.spec.bypass_policy_lockout_safety_check

  is_enabled          = !var.spec.disabled
  enable_key_rotation = var.spec.enable_key_rotation

  rotation_period_in_days = var.spec.rotation_period_in_days != 0 ? var.spec.rotation_period_in_days : null

  multi_region = var.spec.multi_region

  deletion_window_in_days = var.spec.deletion_window_days != 0 ? var.spec.deletion_window_days : null

  # Custom key store surface: setting the store id makes KMS create the key
  # material in the CloudHSM cluster (or, with xks_key_id, forward operations
  # to the named key in an external key manager). Both create-time immutable.
  custom_key_store_id = var.spec.custom_key_store_id != "" ? var.spec.custom_key_store_id : null
  xks_key_id          = var.spec.xks_key_id != "" ? var.spec.xks_key_id : null

  tags = local.aws_tags
}

resource "aws_kms_alias" "this" {
  for_each = local.aliases

  name          = each.value
  target_key_id = aws_kms_key.this.key_id
}

# One KMS grant per spec entry: scoped, revocable permissions for a principal
# to use this key without editing the key policy. Every grant argument is
# create-time immutable (a change replaces the grant -- safe, grants carry no
# state). Entries are keyed by list position; grant identity lives in AWS's
# generated grant id, not the key.
resource "aws_kms_grant" "this" {
  for_each = { for idx, grant in var.spec.grants : tostring(idx) => grant }

  key_id            = aws_kms_key.this.key_id
  grantee_principal = each.value.grantee_principal
  operations        = each.value.operations

  name               = try(each.value.name, "") != "" ? each.value.name : null
  retiring_principal = try(each.value.retiring_principal, "") != "" ? each.value.retiring_principal : null

  # false REVOKES the grant at teardown (immediate hard stop); true RETIRES it
  # (the graceful path AWS recommends once the grant's work is done).
  retire_on_delete = try(each.value.retire_on_delete, false)

  # At most one encryption-context constraint per grant (spec CEL enforces the
  # exclusivity at validate time; the provider only fails it at apply).
  dynamic "constraints" {
    for_each = (length(try(each.value.encryption_context_equals, {})) > 0 ||
    length(try(each.value.encryption_context_subset, {})) > 0) ? [1] : []
    content {
      encryption_context_equals = length(try(each.value.encryption_context_equals, {})) > 0 ? each.value.encryption_context_equals : null
      encryption_context_subset = length(try(each.value.encryption_context_subset, {})) > 0 ? each.value.encryption_context_subset : null
    }
  }
}
