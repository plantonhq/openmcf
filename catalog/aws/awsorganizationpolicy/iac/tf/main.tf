# An Organizations policy (SCP or any of its twelve sibling types)
# with its folded attachments.
#
# Lifecycle facts the render below depends on:
#   - the policy type is immutable (forces replacement); name, content,
#     and description update in place (content diffs are
#     JSON-equivalence-suppressed by the provider);
#   - both attachment leaves are immutable - changing a target
#     re-attaches (detach + attach);
#   - the type must be enabled on the organization
#     (AwsOrganization.enabled_policy_types) before any attachment
#     succeeds - AWS state, not validation, is the referee;
#   - the provider's skip_destroy escape hatches (leave the policy or
#     attachment in place on destroy) are deliberately not modeled -
#     destroy means detach and delete;
#   - AWS-managed policies (FullAWSAccess, ...) can never be adopted -
#     the provider refuses to import them.

resource "aws_organizations_policy" "this" {
  name = var.spec.policy_name

  # Empty keeps the provider default (SERVICE_CONTROL_POLICY).
  type = var.spec.type != "" ? var.spec.type : null

  # The spec carries the document structured; it is serialized to JSON
  # here.
  content = jsonencode(var.spec.content)

  description = var.spec.description != "" ? var.spec.description : null

  tags = local.aws_tags
}

# Each attachment binds the policy to one root, OU, or member account.
resource "aws_organizations_policy_attachment" "this" {
  for_each = local.attachments

  policy_id = aws_organizations_policy.this.id
  target_id = each.value.target_id
}
