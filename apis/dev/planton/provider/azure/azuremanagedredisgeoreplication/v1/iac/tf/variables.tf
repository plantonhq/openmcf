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
  description = "Azure Managed Redis geo-replication group specification"
  type = object({
    # The Managed Redis cluster through which the group is managed, by
    # ARM ID. References are resolved to a literal ID by the platform
    # before the module runs.
    managed_redis_id = string

    # The OTHER members of the group, by ARM ID -- 1 to 4 of them. The
    # managing cluster is always a member implicitly and must not be
    # repeated here (Azure rejects self-links at apply time).
    linked_managed_redis_ids = list(string)
  })
}
