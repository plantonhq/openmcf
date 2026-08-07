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
  description = "Azure Managed Redis access policy assignment (data-plane grant) specification"
  type = object({
    # The Managed Redis instance the grant applies to, by ARM ID (the
    # grant is created on its default database). References are
    # resolved to a literal ID by the platform before the module runs.
    managed_redis_id = string

    # The Entra object (principal) ID being granted, as a GUID. For a
    # managed identity this is the PRINCIPAL id, never the client id.
    object_id = string
  })
}
