# The ElastiCache RBAC user group is the attachment unit between users and
# caches. Membership is modeled here as references (the
# aws_elasticache_user_group_association glue resource is deliberately not
# used): the group is the single place an application's cache access is
# granted or revoked, and this module never mutates the users it references.
#
# Create-only in AWS: user_group_id (metadata.name) and engine. Membership
# (user_ids) updates in place.
resource "aws_elasticache_user_group" "this" {
  user_group_id = local.user_group_id
  engine        = var.spec.engine
  user_ids      = local.user_ids

  tags = local.aws_tags
}
