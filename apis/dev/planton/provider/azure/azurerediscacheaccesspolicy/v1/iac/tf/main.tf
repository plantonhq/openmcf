# A CUSTOM data-plane access policy -- a named permission set in Redis's
# own ACL syntax that AzureRedisCacheAccessPolicyAssignment grants to
# Entra identities. The Redis analog of a custom role definition: this
# says WHAT is allowed; the assignment says WHO gets it.
#
# The three built-in policies ("Data Owner", "Data Contributor", "Data
# Reader") need no policy resource -- assignments reference them by name;
# a custom policy exists for finer grants (one key prefix, no admin
# commands, single commands). Permissions are updatable in place; the
# name and cache are fixed at creation. No tags: ARM does not support
# tags on access policies (they are cache children).
resource "azurerm_redis_cache_access_policy" "main" {
  name           = var.spec.policy_name
  redis_cache_id = var.spec.redis_cache_id
  permissions    = var.spec.permissions
}
