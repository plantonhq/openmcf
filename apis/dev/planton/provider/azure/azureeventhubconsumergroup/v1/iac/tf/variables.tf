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
  description = "Azure Event Hub Consumer Group specification"
  type = object({
    # The parent event hub's ARM ID. References are resolved to a
    # literal by the platform before the module runs.
    event_hub_id = string

    # The group's name -- unique within the hub. ForceNew: renaming
    # replaces the group and resets its consumers' stored offsets.
    consumer_group_name = string

    # Free-form metadata stored on the group (max 1024 characters) --
    # record the owning application or team so operators can tell whose
    # cursor this is.
    user_metadata = optional(string)
  })
}
