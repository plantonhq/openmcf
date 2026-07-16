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
  description = "Azure Cosmos DB SQL container specification"
  type = object({
    # The SQL database the container lives in. References are resolved
    # to a literal ARM ID by the platform before the module runs.
    sql_database_id = string

    # The container's name -- unique within the database.
    container_name = string

    # The partition key paths (each starts with "/"). One path for HASH,
    # two or three for MULTI_HASH hierarchical keys.
    partition_key_paths = list(string)

    # The partition key kind enum name (HASH when empty, MULTI_HASH for
    # hierarchical keys).
    partition_key_kind = optional(string, "")

    # The partition key definition version: 1 (classic) or 2 (large
    # keys; required for MULTI_HASH).
    partition_key_version = optional(number)

    # Fixed dedicated throughput in RU/s. Mutually exclusive with
    # autoscale_max_throughput (enforced by the spec); leave both unset
    # to share the database's throughput.
    throughput = optional(number)

    # Autoscale ceiling in RU/s. Mutually exclusive with throughput.
    autoscale_max_throughput = optional(number)

    # Default document TTL in seconds (-1 = on with no default expiry).
    default_ttl = optional(number)

    # Analytical-store TTL in seconds (-1 = keep forever).
    analytical_storage_ttl = optional(number)

    # Unique key constraints (scoped to the logical partition).
    unique_keys = optional(list(object({
      paths = list(string)
    })), [])

    # The indexing policy; empty applies Azure's default (consistent,
    # index everything).
    indexing_policy = optional(object({
      indexing_mode  = optional(string, "")
      included_paths = optional(list(object({ path = string })), [])
      excluded_paths = optional(list(object({ path = string })), [])
      composite_indexes = optional(list(object({
        entries = list(object({
          path  = string
          order = optional(string, "")
        }))
      })), [])
      spatial_indexes = optional(list(object({ path = string })), [])
    }))

    # Conflict resolution for multi-region-write accounts; empty applies
    # Azure's default (last-writer-wins on /_ts).
    conflict_resolution_policy = optional(object({
      mode                          = string
      conflict_resolution_path      = optional(string, "")
      conflict_resolution_procedure = optional(string, "")
    }))
  })
}
