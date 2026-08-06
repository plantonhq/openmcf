# The MemoryDB ACL is the attachment unit between users and clusters:
# membership is modeled as references, so the ACL is the single place an
# application's cluster access is granted or revoked -- and this module
# never mutates the users it references.
#
# Membership add/remove applies in place (AWS diffs the user set on
# update). An empty ACL is valid: MemoryDB has no mandatory member, a
# cluster attached to an empty ACL simply accepts no authenticated
# connections.
resource "aws_memorydb_acl" "this" {
  name = local.acl_name

  # Refs arrive pre-resolved to plain user names (the platform flattens
  # valueFrom references before the module runs). An empty list becomes
  # null so the create call mirrors Pulumi's.
  user_names = length(var.spec.user_names) > 0 ? var.spec.user_names : null

  tags = local.aws_tags
}
