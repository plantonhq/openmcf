# The ElastiCache RBAC user is a leaf of the RBAC graph: user groups
# reference it by id, and caches reference the groups -- so credential
# material lives here, membership lives on the group, and the cache
# itself never changes when access is granted or revoked.
#
# Create-only in AWS: user_id (metadata.name), user_name, and engine.
# The access string and the authentication mode update in place.
resource "aws_elasticache_user" "this" {
  user_id = local.user_id

  # user_name is what clients present in AUTH; it is NOT unique per user
  # (AWS unions credentials of same-named users), which is why it is a
  # spec field instead of reusing metadata.name.
  user_name     = var.spec.user_name
  engine        = var.spec.engine
  access_string = var.spec.access_string

  # The nested authentication_mode block is the single authentication
  # surface (the provider's legacy flat passwords/no_password_required
  # arms model the same capability and are deliberately not used -- one
  # honest shape). Passwords are only present for the "password" type;
  # CEL guarantees the other types carry none, so an empty list becomes
  # null and AWS sees no password material at all.
  authentication_mode {
    type      = var.spec.authentication_mode.type
    passwords = length(var.spec.authentication_mode.passwords) > 0 ? var.spec.authentication_mode.passwords : null
  }

  tags = local.aws_tags
}
