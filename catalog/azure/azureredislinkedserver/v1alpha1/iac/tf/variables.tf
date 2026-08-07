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
  description = "Azure Redis Linked Server (geo-replication link) specification"
  type = object({
    # The PRIMARY cache's ARM ID -- the cache serving writes. The link is
    # created as its child; the cache's name and resource group are parsed
    # from this ID. References are resolved to a literal ARM ID by the
    # platform before the module runs.
    target_redis_cache_id = string

    # The SECONDARY cache's ARM ID -- the DR replica in another region.
    linked_redis_cache_id = string

    # The SECONDARY cache's region. Normally referenced from the same
    # cache as linked_redis_cache_id so it can never disagree.
    linked_redis_cache_location = string

    # The linked cache's replication role, as the spec enum's name string
    # (PRIMARY / SECONDARY).
    server_role = string
  })
}
