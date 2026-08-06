# The active geo-replication group: links Managed Redis instances whose
# default databases declare the same geo_replication_group_name into a
# multi-primary replica set -- every member accepts writes and Azure
# merges the datasets with conflict-free semantics.
#
# Membership is managed as its own resource because linking mutates the
# replication state of EVERY member out of band -- it is a group-wide
# operation performed through one member, not a property of any single
# instance. ONE resource manages the whole group (linking is reciprocal;
# never create one per member). Deleting it force-unlinks the members,
# each keeping its own copy of the data; removing a single ID from the
# linked list evacuates just that member.
#
# Cross-resource contracts Azure enforces at link time: every member
# carries the SAME group name, is BALANCED_B3 or larger, has no
# persistence, and uses only the RediSearch/RedisJSON modules. No tags:
# the group has no ARM object of its own.
resource "azurerm_managed_redis_geo_replication" "main" {
  managed_redis_id         = var.spec.managed_redis_id
  linked_managed_redis_ids = var.spec.linked_managed_redis_ids
}
