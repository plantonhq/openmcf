locals {
  # The account's use-case form renders only when the manifest declares it
  # (required once per account for Anthropic models; other vendors need no
  # form).
  has_use_case_form = var.spec.use_case_form != null
}
