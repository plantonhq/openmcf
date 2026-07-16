# The data-plane grant: assigns an access policy (built-in or an
# AzureRedisCacheAccessPolicy) to a Microsoft Entra identity on the cache
# -- the Redis analog of a role assignment. The granted identity connects
# with its object ID (or the alias below) as the Redis username and an
# Entra token as the password; Entra auth must be enabled on the cache
# for the grant to matter.
#
# Every argument is ForceNew: replacing the assignment momentarily
# revokes and re-grants, which is safe for the grant class. No tags: ARM
# does not support tags on access policy assignments (cache children).
resource "azurerm_redis_cache_access_policy_assignment" "main" {
  name           = var.spec.assignment_name
  redis_cache_id = var.spec.redis_cache_id

  access_policy_name = var.spec.access_policy_name

  # For a managed identity this must be the PRINCIPAL id -- granting the
  # client id fails at connect time, not at deploy time.
  object_id       = var.spec.object_id
  object_id_alias = var.spec.object_id_alias
}
