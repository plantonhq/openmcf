# One CloudWatch Logs account-level policy: a per-(name, type) policy
# object applied account-wide in the region.
#
# Lifecycle facts the render below depends on:
#   - policy_name and policy_type are BOTH identity (changing either
#     replaces the policy; the provider imports the pair as
#     "policy_name:policy_type");
#   - selection_criteria narrows the account-wide scope and also
#     replaces on change; only the document (and AWS's scope argument)
#     update in place;
#   - the provider's `scope` argument is deliberately pinned to its only
#     legal value (ALL) rather than modeled - a recorded exclusion.

resource "aws_cloudwatch_log_account_policy" "this" {
  policy_name     = var.spec.policy_name
  policy_type     = var.spec.policy_type
  policy_document = local.policy_document

  selection_criteria = var.spec.selection_criteria != "" ? var.spec.selection_criteria : null

  # ALL is the only value the provider's Scope enum carries at the pin.
  scope = "ALL"
}
