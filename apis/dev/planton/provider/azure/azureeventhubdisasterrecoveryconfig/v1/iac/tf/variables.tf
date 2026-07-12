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
  description = "Azure Event Hub Disaster Recovery Config specification"
  type = object({
    # The failover-stable alias DNS name
    # ({alias_name}.servicebus.windows.net). ForceNew.
    alias_name = string

    # The PRIMARY namespace's ARM ID. References are resolved to literals
    # by the platform before the module runs. ForceNew.
    primary_namespace_id = string

    # The PARTNER namespace's ARM ID (different region, same tier, empty
    # at pairing time). Changing it breaks the pairing and re-pairs.
    partner_namespace_id = string
  })
}
