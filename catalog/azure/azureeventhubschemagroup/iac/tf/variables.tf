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
  description = "Azure Event Hub Schema Group specification"
  type = object({
    # The Event Hubs namespace whose registry holds the group, by ARM
    # ID. References are resolved to a literal by the platform before
    # the module runs. ForceNew.
    namespace_id = string

    # The group's name -- unique within the namespace; serializers
    # address the group by this name. ForceNew: renaming replaces the
    # group and drops its registered schemas.
    schema_group_name = string

    # The evolution policy, as the spec enum's value name (NONE,
    # BACKWARD, FORWARD). ForceNew.
    schema_compatibility = string

    # The serialization format, as the spec enum's value name (AVRO,
    # JSON). ForceNew.
    schema_type = string
  })
}
