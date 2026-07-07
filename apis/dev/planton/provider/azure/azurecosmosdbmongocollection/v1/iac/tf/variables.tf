variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Azure Cosmos DB MongoDB collection specification"
  type = object({
    # The Mongo database the collection lives in. References are
    # resolved to a literal ARM ID by the platform before the module
    # runs.
    mongo_database_id = string

    # The collection's name -- unique within the database.
    collection_name = string

    # The shard (partition) key property. Empty creates an unsharded
    # collection confined to one physical partition.
    shard_key = optional(string, "")

    # Fixed dedicated throughput in RU/s. Mutually exclusive with
    # autoscale_max_throughput (enforced by the spec); leave both unset
    # to share the database's throughput.
    throughput = optional(number)

    # Autoscale ceiling in RU/s. Mutually exclusive with throughput.
    autoscale_max_throughput = optional(number)

    # Default document TTL in seconds (-1 = on with no default expiry;
    # never 0).
    default_ttl_seconds = optional(number)

    # Analytical-store TTL in seconds (-1 = keep forever).
    analytical_storage_ttl = optional(number)

    # Indexes on the collection (Azure requires the ["_id"] unique
    # index on every Mongo collection; the spec enforces it).
    indexes = optional(list(object({
      keys   = list(string)
      unique = optional(bool, false)
    })), [])
  })
}
