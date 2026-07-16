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
  description = "Azure Service Bus Queue specification"
  type = object({
    # The parent namespace's ARM ID. References are resolved to a
    # literal by the platform before the module runs.
    namespace_id = string

    # The queue's name -- unique within the namespace. ForceNew.
    queue_name = string

    # Maximum queue size in MB (Azure's fixed ladder; large sizes are
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

    # PeekLock lock duration, ISO 8601, PT5S-PT5M (Azure default PT1M).
    lock_duration = optional(string)

    # Delivery attempts before dead-lettering (Azure default 10).
    max_delivery_count = optional(number)

    # Move expired messages to the dead-letter sub-queue.
    dead_lettering_on_message_expiration = optional(bool)

    # Enable sessions (FIFO + exclusive consumption). ForceNew.
    requires_session = optional(bool)

    # Auto-delete after idling this long, ISO 8601, min PT5M (unset =
    # never).
    auto_delete_on_idle = optional(string)

    # Allow batched broker operations (Azure default true).
    batched_operations_enabled = optional(bool, true)

    # Express Entities: in-memory buffering. BASIC/STANDARD only;
    # incompatible with duplicate detection.
    express_enabled = optional(bool)

    # Auto-forward to another entity in the same namespace, by name.
    forward_to = optional(string)

    # Auto-forward dead-lettered messages, by entity name.
    forward_dead_lettered_messages_to = optional(string)

    # Gate state, as the spec enum's value name (ACTIVE, DISABLED,
    # SEND_DISABLED, RECEIVE_DISABLED). Unset deploys ACTIVE.
    status = optional(string)
  })
}
