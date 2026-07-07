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
  description = "Azure Front Door origin group specification"
  type = object({
    # The Front Door profile the origin group lives in, by ARM ID.
    # References are resolved to a literal ID by the platform before the
    # module runs. ForceNew.
    profile_id = string

    # The origin group's name -- unique within the profile. ForceNew.
    origin_group_name = string

    # Traffic distribution across the group's origins. Azure requires
    # load-balancing settings on every group, so the module always sends
    # the block -- with these values when set, Azure's defaults
    # (4 / 3 / 50 ms) otherwise.
    load_balancing = optional(object({
      sample_size                        = optional(number)
      successful_samples_required        = optional(number)
      additional_latency_in_milliseconds = optional(number)
    }))

    # Periodic origin health probing. Absent means probing disabled
    # (all origins assumed healthy) -- a real behavior, not a defaults
    # shortcut.
    health_probe = optional(object({
      # HTTP / HTTPS (spec enum value names).
      protocol            = string
      interval_in_seconds = number
      # HEAD / GET (spec enum value names). Absent means HEAD.
      request_type = optional(string)
      path         = optional(string)
    }))

    # Cookie-based session affinity. Azure defaults to true when omitted.
    session_affinity_enabled = optional(bool)

    # Minutes to ramp traffic onto a just-healed or new origin (0-50).
    # Azure defaults to 10 when omitted.
    restore_traffic_time_to_healed_or_new_endpoint_in_minutes = optional(number)
  })
}
