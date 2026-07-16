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
  description = "Azure Service Bus Authorization Rule specification"
  type = object({
    # The rule's name -- unique within its scope. ForceNew.
    rule_name = string

    # Exactly ONE of the three scopes is set (spec-enforced XOR); each
    # is resolved to a literal ARM ID by the platform before the module
    # runs. The scope picks which azurerm resource the module creates.
    namespace_id = optional(string)
    queue_id     = optional(string)
    topic_id     = optional(string)

    # The rights trio (Azure's contract: at least one true; manage
    # requires listen AND send -- both spec-enforced).
    listen = optional(bool)
    send   = optional(bool)
    manage = optional(bool)
  })
}
