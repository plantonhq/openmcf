locals {
  # The primary cache's name and resource group, parsed from its ARM ID --
  # the link is that cache's child, so neither is ever spelled twice. The
  # named-group regexes fail the plan loudly if the ID is not a Redis
  # cache ARM ID. The type segment is matched case-insensitively: ARM has
  # emitted both .../Microsoft.Cache/Redis/{name} and .../redis/{name}
  # over the API's life, and ARM ID comparison is case-insensitive on
  # type segments.
  target_cache_name = regex("(?i)/redis/(?P<name>[^/]+)$", var.spec.target_redis_cache_id)["name"]

  target_resource_group_name = regex("(?i)/resourcegroups/(?P<name>[^/]+)/", var.spec.target_redis_cache_id)["name"]

  # The spec's role enum arrives as the FULL proto value name; ARM wants
  # the capitalized word.
  server_role_map = {
    "PRIMARY"   = "Primary"
    "SECONDARY" = "Secondary"
  }
  server_role = local.server_role_map[var.spec.server_role]
}
