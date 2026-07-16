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
  description = "Azure Service Bus Topic specification"
  type = object({
    # The parent namespace's ARM ID. References are resolved to a
    # literal by the platform before the module runs.
    namespace_id = string

    # The topic's name -- unique within the namespace. ForceNew.
    topic_name = string

    # Maximum topic size in MB (Azure's fixed ladder; large sizes are
    # PREMIUM-only). Unset lets Azure default for the tier.
    max_size_in_megabytes = optional(number)

    # Largest accepted message in KB (1024-102400). PREMIUM only --
    # multi-tenant tiers are fixed at 256 KB.
    max_message_size_in_kilobytes = optional(number)

    # Spread across multiple message stores. ForceNew. On PREMIUM the
    # namespace's partition layout dictates this (apply-time contract).
    partitioning_enabled = optional(bool)

    # Track MessageIds and drop duplicates. ForceNew.
    requires_duplicate_detection = optional(bool)

    # Message time-to-live, ISO 8601 (unset = unbounded).
    default_message_ttl = optional(string)

    # Duplicate-detection history window, ISO 8601 (Azure default PT10M).
    duplicate_detection_history_time_window = optional(string)

    # Auto-delete after idling this long, ISO 8601, min PT5M (unset =
    # never).
    auto_delete_on_idle = optional(string)

    # Allow batched broker operations (Azure default true).
    batched_operations_enabled = optional(bool, true)

    # Express Entities: in-memory buffering. STANDARD only; incompatible
    # with duplicate detection.
    express_enabled = optional(bool)

    # Preserve publish order for session-aware subscriptions.
    support_ordering = optional(bool)

    # Gate state, as the spec enum's value name (ACTIVE, DISABLED).
    # Unset deploys ACTIVE.
    status = optional(string)
  })
}
