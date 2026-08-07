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
  description = "Azure Service Bus Subscription specification"
  type = object({
    # The parent topic's ARM ID. References are resolved to a literal
    # by the platform before the module runs.
    topic_id = string

    # The subscription's name -- unique within the topic. ForceNew.
    subscription_name = string

    # Delivery attempts before dead-lettering (required -- Azure has no
    # server default on subscriptions).
    max_delivery_count = number

    # PeekLock lock duration, ISO 8601, PT5S-PT5M (Azure default PT1M).
    lock_duration = optional(string)

    # Message time-to-live, ISO 8601 (unset inherits the topic's TTL).
    default_message_ttl = optional(string)

    # Auto-delete after idling this long, ISO 8601, min PT5M (unset =
    # never).
    auto_delete_on_idle = optional(string)

    # Move expired messages to the dead-letter sub-queue.
    dead_lettering_on_message_expiration = optional(bool)

    # Dead-letter messages that FAIL filter evaluation (Azure default
    # true).
    dead_lettering_on_filter_evaluation_error = optional(bool, true)

    # Enable sessions (FIFO + exclusive consumption). ForceNew.
    requires_session = optional(bool)

    # Allow batched broker operations (Azure's subscription-side
    # default is false).
    batched_operations_enabled = optional(bool)

    # Auto-forward to another entity in the same namespace, by name.
    forward_to = optional(string)

    # Auto-forward dead-lettered messages, by entity name.
    forward_dead_lettered_messages_to = optional(string)

    # Gate state, as the spec enum's value name (ACTIVE, DISABLED,
    # RECEIVE_DISABLED). Unset deploys ACTIVE.
    status = optional(string)

    # Client-scoped (JMS 2.0 client-affine) binding. ForceNew fields.
    client_scoped_subscription = optional(object({
      client_id = optional(string)
      shareable = optional(bool, true)
    }))

    # Filter rules (folded -- rules have no life outside their
    # subscription). OR semantics; Azure's auto-created $Default
    # catch-all remains unless a declared rule is itself named
    # "$Default".
    rules = optional(list(object({
      rule_name = string
      # SQL_FILTER or CORRELATION_FILTER (the spec enum's value name).
      filter_type = string
      # Required with SQL_FILTER.
      sql_filter = optional(string)
      # Required with CORRELATION_FILTER; at least one matcher set.
      correlation_filter = optional(object({
        correlation_id      = optional(string)
        message_id          = optional(string)
        to                  = optional(string)
        reply_to            = optional(string)
        label               = optional(string)
        session_id          = optional(string)
        reply_to_session_id = optional(string)
        content_type        = optional(string)
        properties          = optional(map(string), {})
      }))
      # Optional SQL action annotating matched messages.
      action = optional(string)
    })), [])
  })
}
