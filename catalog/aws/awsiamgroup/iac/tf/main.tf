# An IAM group: the container that grants a set of users a shared
# permission set - its policies are the permissions, the declarative
# users list is the membership.
#
# Lifecycle facts the render below depends on:
#   - the group's name comes from metadata.name and updates IN PLACE at
#     AWS (UpdateGroup renames; the ARN recomputes but members and
#     policies persist);
#   - membership renders as ONE aws_iam_group_membership resource
#     carrying the whole users list - the AUTHORITATIVE form: users
#     added out-of-band are removed on the next apply, and clearing the
#     spec list removes the resource (and with it every membership).
#     The users must already exist at IAM;
#   - IAM groups and every satellite here are untaggable at AWS.
resource "aws_iam_group" "this" {
  name = var.metadata.name
  path = var.spec.path != "" ? var.spec.path : null
}

# One membership resource owns the whole users list (the group-centric
# declarative form). The membership resource's own name is a
# provider-side label (never reaches AWS).
resource "aws_iam_group_membership" "this" {
  count = length(var.spec.users) > 0 ? 1 : 0

  name  = "${var.metadata.name}-membership"
  group = aws_iam_group.this.name
  users = var.spec.users
}

# Each managed-policy attachment is its own resource so attachments
# reconcile individually: adding or removing an entry attaches or
# detaches just that policy, and attachments made outside this resource
# are left alone.
resource "aws_iam_group_policy_attachment" "managed" {
  for_each   = toset(var.spec.managed_policy_arns)
  group      = aws_iam_group.this.name
  policy_arn = each.value
}

# Inline policies live and die with the group - permissions unique to
# this one group that would be noise as standalone AwsIamPolicy
# resources.
resource "aws_iam_group_policy" "inline" {
  for_each = local.inline_policies_json
  name     = each.key
  group    = aws_iam_group.this.id
  # each.value is already a JSON-encoded policy string (see locals.inline_policies_json).
  policy = each.value
}
