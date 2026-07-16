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

  tags = local.aws_tags
}

resource "aws_kms_alias" "this" {
  for_each = local.aliases

  name          = each.value
  target_key_id = aws_kms_key.this.key_id
}
