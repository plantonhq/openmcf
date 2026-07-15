# The MemoryDB user is the leaf of MemoryDB's ACL graph: ACLs reference it
# by name, and clusters attach the ACLs -- so credential material lives
# here, membership lives on the ACL, and the cluster itself never changes
# when access is granted or revoked.
#
# Create-only in AWS: user_name (metadata.name). The access string and the
# authentication mode update in place, so tightening permissions or
# rotating passwords never replaces the user.
resource "aws_memorydb_user" "this" {
  user_name     = local.user_name
  access_string = var.spec.access_string

  # Exactly two input types exist ("password" carries 1-2 secrets; "iam"
  # carries none) -- CEL guarantees the coupling, so an empty password list
  # becomes null and AWS sees no credential material at all.
  authentication_mode {
    type      = var.spec.authentication_mode.type
    passwords = length(var.spec.authentication_mode.passwords) > 0 ? var.spec.authentication_mode.passwords : null
  }

  tags = local.aws_tags
}
