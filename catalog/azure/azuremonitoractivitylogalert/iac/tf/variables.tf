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
  description = "Azure Monitor Activity Log Alert specification"
  type = object({
    resource_group = string
    name           = string

    # The alert-resource region as the spec enum name (GLOBAL / WEST_EUROPE
    # / NORTH_EUROPE / EAST_US_2_EUAP); unset means GLOBAL.
    location = optional(string)

    # The scopes watched, resolved to literal ARM IDs.
    scopes = list(string)

    # The matching criteria. Enum fields are the spec enum name strings; the
    # plural list fields carry the single case too.
    criteria = object({
      category                = string
      operation_name          = optional(string, "")
      caller                  = optional(string, "")
      levels                  = optional(list(string), [])
      resource_providers      = optional(list(string), [])
      resource_types          = optional(list(string), [])
      resource_groups         = optional(list(string), [])
      resource_ids            = optional(list(string), [])
      statuses                = optional(list(string), [])
      sub_statuses            = optional(list(string), [])
      recommendation_category = optional(string, "")
      recommendation_impact   = optional(string, "")
      recommendation_type     = optional(string, "")
      resource_health = optional(object({
        current  = optional(list(string), [])
        previous = optional(list(string), [])
        reason   = optional(list(string), [])
      }))
      service_health = optional(object({
        events    = optional(list(string), [])
        locations = optional(list(string), [])
        services  = optional(list(string), [])
      }))
    })

    # The actions fired on match; action_group_id resolved to a literal ARM
    # ID.
    actions = optional(list(object({
      action_group_id    = string
      webhook_properties = optional(map(string), {})
    })), [])

    description = optional(string, "")
    enabled     = optional(bool)

    tags = optional(map(string), {})
  })
}
