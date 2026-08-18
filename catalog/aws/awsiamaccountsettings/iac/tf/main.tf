# IAM's account-level settings - a SETTINGS SINGLETON for a GLOBAL
# service: IAM keeps exactly one of each setting per ACCOUNT, so deploy
# at most one instance per account. metadata.name never reaches AWS.
#
# Destroy semantics DIFFER per arm:
#   - the alias truly DELETES (sign-in URLs revert to the account ID);
#   - the password policy RESETS to AWS's defaults;
#   - the STS preference's delete is a NO-OP (the last-applied token
#     version persists; reverting is an apply with the other version).
#
# The password policy is replaced WHOLE on every apply (AWS's
# UpdateAccountPasswordPolicy semantics): an unset field is AWS's
# default, never "keep the current setting". The provider also never
# sends false/0 values on the wire (indistinguishable from unset at the
# SDK layer) - AWS treats a missing toggle as false, so the rendered
# result is identical.

# The account's sign-in alias. An account has exactly ONE alias, and
# aliases are GLOBALLY unique across all of AWS - applying this arm
# REPLACES whatever alias the account already had, which changes the
# sign-in URL everyone uses.
resource "aws_iam_account_alias" "this" {
  count = local.manage_alias ? 1 : 0

  account_alias = var.spec.account_alias
}

resource "aws_iam_account_password_policy" "this" {
  count = local.manage_password_policy ? 1 : 0

  # Plain-bool toggles render unconditionally: false and unset are the
  # same posture at AWS (the policy is replaced whole), and rendering
  # both engines identically keeps state parity.
  require_lowercase_characters = var.spec.password_policy.require_lowercase_characters
  require_numbers              = var.spec.password_policy.require_numbers
  require_symbols              = var.spec.password_policy.require_symbols
  require_uppercase_characters = var.spec.password_policy.require_uppercase_characters
  hard_expiry                  = var.spec.password_policy.hard_expiry

  # Presence-typed knobs render only when set - their AWS defaults
  # (6-character minimum, self-service changes allowed, no expiry, no
  # reuse prevention) apply otherwise.
  minimum_password_length        = var.spec.password_policy.minimum_password_length
  allow_users_to_change_password = var.spec.password_policy.allow_users_to_change_password
  max_password_age               = var.spec.password_policy.max_password_age
  password_reuse_prevention      = var.spec.password_policy.password_reuse_prevention
}

# The STS global-endpoint token version. v2Token works in ALL regions
# including opt-in ones; v1Token (the AWS default) only in
# default-enabled regions.
resource "aws_iam_security_token_service_preferences" "this" {
  count = local.manage_sts ? 1 : 0

  global_endpoint_token_version = var.spec.sts.global_endpoint_token_version
}

# The account id feeds the output regardless of which arms render
# (every resource here is account-scoped).
data "aws_caller_identity" "this" {}
