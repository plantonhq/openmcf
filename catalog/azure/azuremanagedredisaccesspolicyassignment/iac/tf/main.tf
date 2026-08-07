# The data-plane grant: assigns Managed Redis's built-in "default"
# access policy to a Microsoft Entra identity -- the Redis analog of a
# role assignment, and the grant half of the keyless-by-default story
# (access keys are off unless enabled, so grants are how clients connect
# at all). The granted identity presents its object ID as the Redis
# username and an Entra token as the password.
#
# Azure names the assignment after the object ID, so an identity is
# granted at most once per database -- there is nothing else to name.
# Every argument is ForceNew: replacing the assignment momentarily
# revokes and re-grants, which is safe for the grant class. No tags: ARM
# does not support tags on access policy assignments (database
# children).
resource "azurerm_managed_redis_access_policy_assignment" "main" {
  managed_redis_id = var.spec.managed_redis_id

  # For a managed identity this must be the PRINCIPAL id -- granting the
  # client id fails at connect time, not at deploy time.
  object_id = var.spec.object_id
}
