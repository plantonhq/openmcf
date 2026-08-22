locals {
  # None of IAM's account-settings resources is taggable at AWS, so this
  # module carries no tag map - the one deliberate absence against the
  # catalog's tag convention (mirrored in the Pulumi module).

  # Arm presence: each arm renders ONLY when its spec field/message is
  # present - an omitted arm leaves the account's current setting
  # untouched, and that omission is meaningful and deliberate.
  manage_alias           = var.spec.account_alias != ""
  manage_password_policy = var.spec.password_policy != null
  manage_sts             = var.spec.sts != null
}
