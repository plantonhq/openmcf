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
  description = "Azure Event Hub Authorization Rule specification"
  type = object({
    # The rule's name -- the SharedAccessKeyName clients present.
    # ForceNew: renaming replaces the rule and regenerates its keys.
    rule_name = string

    # Namespace-wide scope, by ARM id (resolved to a literal before the
    # module runs). Exactly one scope is set (spec-enforced XOR).
    namespace_id = optional(string)

    # Single-hub scope, by ARM id.
    event_hub_id = optional(string)

    # The rights trio. At least one is true; manage requires listen and
    # send (spec-enforced, mirroring Azure's own contract).
    listen = optional(bool)
    send   = optional(bool)
    manage = optional(bool)
  })
}
