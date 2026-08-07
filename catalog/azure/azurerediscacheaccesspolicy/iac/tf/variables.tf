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
  description = "Azure Redis Cache Access Policy (custom data-plane permission set) specification"
  type = object({
    # The cache the policy is defined on, by ARM ID. References are
    # resolved to a literal ARM ID by the platform before the module runs.
    redis_cache_id = string

    # The policy's name -- what assignments reference; unique within the
    # cache and never one of the built-in names (spec-enforced).
    policy_name = string

    # The permission set in Redis ACL syntax, e.g.
    # "+@read +@connection ~app:*".
    permissions = string
  })
}
