# The access entry composes onto its neighbors instead of embedding
# them: the cluster and the IAM principal attach by reference. AWS keys
# the entry on (cluster, principal) -- both create-only -- while groups,
# username, and policy associations update in place.
#
# The cluster must have API authentication enabled
# (access_config.authentication_mode "API" or "API_AND_CONFIG_MAP");
# CONFIG_MAP-only clusters reject access entries at create time.
resource "aws_eks_access_entry" "this" {
  cluster_name  = var.spec.cluster_name
  principal_arn = var.spec.principal_arn

  # Empty means STANDARD (the AWS default). The node types exist for
  # self-managed/hybrid node registration; the spec's CEL rules keep
  # groups/username/associations off them, mirroring AWS's runtime rules.
  type = var.spec.type != "" ? var.spec.type : null

  kubernetes_groups = length(var.spec.kubernetes_groups) > 0 ? var.spec.kubernetes_groups : null

  # Empty lets AWS default the username (principal ARN for users; a
  # session-templated name for roles, which preserves the session name
  # in audit logs).
  user_name = var.spec.user_name != "" ? var.spec.user_name : null

  tags = local.aws_tags
}

# The folded policy associations: AWS sub-resources of exactly this
# (cluster, principal) pair, materialized one provider resource per spec
# entry and keyed by the policy name -- adding, re-scoping, or removing
# one association diffs in place and never touches the entry or its
# siblings. AWS allows one association per policy per principal, so the
# key is unique by construction. The explicit depends_on makes the
# ordering AWS requires (entry before association) visible.
resource "aws_eks_access_policy_association" "this" {
  for_each = local.policy_associations_by_name

  cluster_name  = var.spec.cluster_name
  principal_arn = var.spec.principal_arn
  policy_arn    = each.value.policy_arn

  access_scope {
    type       = each.value.access_scope.type
    namespaces = length(each.value.access_scope.namespaces) > 0 ? each.value.access_scope.namespaces : null
  }

  depends_on = [aws_eks_access_entry.this]
}
