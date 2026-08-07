# The MongoDB API collection -- the unit of storage and scale-out where
# documents live. Addressed by the (resource group, account, database,
# name) tuple azurerm requires, parsed from the parent database's ARM
# ID. No Azure tags: ARM does not support tags on Cosmos child
# resources, so the platform's identity tags live on the account.
resource "azurerm_cosmosdb_mongo_collection" "main" {
  name                = var.spec.collection_name
  resource_group_name = local.resource_group_name
  account_name        = local.cosmosdb_account_name
  database_name       = local.mongo_database_name

  # The shard key -- the MongoDB face of the partition key, fixed at
  # creation. Sent only when set: empty creates an unsharded collection
  # confined to one physical partition.
  shard_key = var.spec.shard_key != "" ? var.spec.shard_key : null

  # Dedicated throughput. Sent only when set: serverless accounts
  # reject provisioned throughput, and unset means the collection
  # shares the database's throughput. The spec enforces the
  # fixed-XOR-autoscale contract before the plan ever runs.
  throughput = var.spec.throughput

  dynamic "autoscale_settings" {
    for_each = var.spec.autoscale_max_throughput != null ? [var.spec.autoscale_max_throughput] : []
    content {
      max_throughput = autoscale_settings.value
    }
  }

  # Document TTL (implemented by Cosmos DB as an expireAfter index on
  # _ts): -1 turns TTL on with per-document expiry only; a positive
  # value expires documents after their last write. Never 0 (the spec
  # rejects it -- ARM's contract).
  default_ttl_seconds = var.spec.default_ttl_seconds

  # Analytical-store TTL (requires analytical storage on the account):
  # -1 keeps analytical data forever.
  analytical_storage_ttl = var.spec.analytical_storage_ttl

  # Indexes, including the ["_id"] unique index Azure requires on
  # every Mongo collection (spec-enforced) -- declared explicitly,
  # never injected.
  dynamic "index" {
    for_each = var.spec.indexes
    content {
      keys   = index.value.keys
      unique = index.value.unique
    }
  }
}
