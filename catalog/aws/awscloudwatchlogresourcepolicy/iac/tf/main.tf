# One CloudWatch Logs resource policy: the service log-write grant, in
# exactly one scope.
#
# Lifecycle facts the render below depends on:
#   - the spec guarantees exactly one of policy_name (account scope) and
#     resource_arn (resource scope); both are identity - changing either
#     replaces the policy;
#   - updates pass AWS's revision ID from state (optimistic concurrency),
#     so concurrent out-of-band edits fail loudly instead of being
#     overwritten;
#   - resource-scoped deletes REQUIRE the tracked revision ID - never
#     clear it from state by hand.

resource "aws_cloudwatch_log_resource_policy" "this" {
  policy_document = local.policy_document

  policy_name  = var.spec.policy_name != "" ? var.spec.policy_name : null
  resource_arn = var.spec.resource_arn != "" ? var.spec.resource_arn : null
}
