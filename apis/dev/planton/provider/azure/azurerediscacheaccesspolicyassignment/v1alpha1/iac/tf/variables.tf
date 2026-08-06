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
  description = "Azure Redis Cache Access Policy Assignment (data-plane grant) specification"
  type = object({
    # The cache the grant applies to, by ARM ID. References are resolved
    # to a literal ARM ID by the platform before the module runs.
    redis_cache_id = string

    # The assignment's name -- a label for the grant, unique within the
    # cache.
    assignment_name = string

    # The policy being granted: a built-in name ("Data Owner",
    # "Data Contributor", "Data Reader") or a custom
    # AzureRedisCacheAccessPolicy's name. References are resolved to a
    # literal name by the platform before the module runs.
    access_policy_name = string

    # The Entra object (principal) ID being granted -- for a managed
    # identity this is the PRINCIPAL id, not the client id.
    object_id = string

    # A human-readable alias that doubles as an alternative Redis
    # username at connect time.
    object_id_alias = string
  })
}
