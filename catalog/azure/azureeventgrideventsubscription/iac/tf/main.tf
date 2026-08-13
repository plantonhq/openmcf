# The addressing choice (spec-enforced to exactly one) selects which
# of the two provider resources below materializes. Azure's two
# subscription resources share ONE configuration grammar (the provider
# generates both from the same schema), so the two bodies here are
# deliberately identical except for the addressing arguments -- keep
# them in lockstep when either changes.
#
# Provider seams mirrored in BOTH bodies:
#   - expiration/batch/TTL/Entra fields are sent only when set (the
#     service owns their defaults);
#   - delivery properties are ignored by Azure on storage-queue
#     destinations (queue messages carry no custom properties);
#   - retry_policy is sent only when set -- Azure's defaults (30
#     attempts / 1440 minutes) echo back on read otherwise;
#   - labels and included_event_types pass through as lists (an empty
#     list means none / all event types respectively).

# A scope-addressed subscription: attaches to any ARM resource that
# emits Event Grid events (custom topic, domain, domain topic,
# resource group, subscription).
resource "azurerm_eventgrid_event_subscription" "main" {
  count = var.spec.scope != null ? 1 : 0

  name  = var.spec.name
  scope = var.spec.scope

  # Always sent (platform default mirrors Azure's). Create-only.
  event_delivery_schema = var.spec.event_delivery_schema

  expiration_time_utc = var.spec.expiration_time_utc != "" ? var.spec.expiration_time_utc : null

  included_event_types = length(var.spec.included_event_types) > 0 ? var.spec.included_event_types : null

  labels = var.spec.labels

  advanced_filtering_on_arrays_enabled = var.spec.advanced_filtering_on_arrays_enabled

  # The four id-arm destinations pass through (null when the arm is
  # not the chosen one).
  eventhub_id          = var.spec.destination.eventhub_id
  hybrid_connection_id = var.spec.destination.hybrid_connection_id
  service_bus_queue_id = var.spec.destination.service_bus_queue_id
  service_bus_topic_id = var.spec.destination.service_bus_topic_id

  dynamic "azure_function_endpoint" {
    for_each = var.spec.destination.azure_function != null ? [var.spec.destination.azure_function] : []
    content {
      function_id                       = azure_function_endpoint.value.function_id
      max_events_per_batch              = azure_function_endpoint.value.max_events_per_batch
      preferred_batch_size_in_kilobytes = azure_function_endpoint.value.preferred_batch_size_in_kilobytes
    }
  }

  dynamic "storage_queue_endpoint" {
    for_each = var.spec.destination.storage_queue != null ? [var.spec.destination.storage_queue] : []
    content {
      storage_account_id                    = storage_queue_endpoint.value.storage_account_id
      queue_name                            = storage_queue_endpoint.value.queue_name
      queue_message_time_to_live_in_seconds = storage_queue_endpoint.value.queue_message_time_to_live_in_seconds
    }
  }

  dynamic "webhook_endpoint" {
    for_each = var.spec.destination.webhook != null ? [var.spec.destination.webhook] : []
    content {
      url                               = webhook_endpoint.value.url
      max_events_per_batch              = webhook_endpoint.value.max_events_per_batch
      preferred_batch_size_in_kilobytes = webhook_endpoint.value.preferred_batch_size_in_kilobytes
      active_directory_tenant_id        = webhook_endpoint.value.active_directory_tenant_id != "" ? webhook_endpoint.value.active_directory_tenant_id : null
      active_directory_app_id_or_uri    = webhook_endpoint.value.active_directory_app_id_or_uri != "" ? webhook_endpoint.value.active_directory_app_id_or_uri : null
    }
  }

  dynamic "delivery_identity" {
    for_each = var.spec.delivery_identity != null ? [var.spec.delivery_identity] : []
    content {
      type                   = local.identity_type_map[delivery_identity.value.type]
      user_assigned_identity = delivery_identity.value.user_assigned_identity != "" ? delivery_identity.value.user_assigned_identity : null
    }
  }

  dynamic "delivery_property" {
    for_each = var.spec.delivery_properties
    content {
      header_name  = delivery_property.value.header_name
      type         = delivery_property.value.type
      value        = delivery_property.value.value != "" ? delivery_property.value.value : null
      source_field = delivery_property.value.source_field != "" ? delivery_property.value.source_field : null
      secret       = delivery_property.value.secret
    }
  }

  dynamic "storage_blob_dead_letter_destination" {
    for_each = var.spec.dead_letter != null ? [var.spec.dead_letter] : []
    content {
      storage_account_id          = storage_blob_dead_letter_destination.value.storage_account_id
      storage_blob_container_name = storage_blob_dead_letter_destination.value.storage_blob_container_name
    }
  }

  dynamic "dead_letter_identity" {
    for_each = var.spec.dead_letter_identity != null ? [var.spec.dead_letter_identity] : []
    content {
      type                   = local.identity_type_map[dead_letter_identity.value.type]
      user_assigned_identity = dead_letter_identity.value.user_assigned_identity != "" ? dead_letter_identity.value.user_assigned_identity : null
    }
  }

  dynamic "retry_policy" {
    for_each = var.spec.retry_policy != null ? [var.spec.retry_policy] : []
    content {
      max_delivery_attempts = retry_policy.value.max_delivery_attempts
      event_time_to_live    = retry_policy.value.event_time_to_live
    }
  }

  dynamic "subject_filter" {
    for_each = var.spec.subject_filter != null ? [var.spec.subject_filter] : []
    content {
      subject_begins_with = subject_filter.value.subject_begins_with != "" ? subject_filter.value.subject_begins_with : null
      subject_ends_with   = subject_filter.value.subject_ends_with != "" ? subject_filter.value.subject_ends_with : null
      case_sensitive      = subject_filter.value.case_sensitive
    }
  }

  dynamic "advanced_filter" {
    for_each = var.spec.advanced_filter != null ? [var.spec.advanced_filter] : []
    content {
      dynamic "bool_equals" {
        for_each = advanced_filter.value.bool_equals
        content {
          key   = bool_equals.value.key
          value = bool_equals.value.value
        }
      }
      dynamic "number_greater_than" {
        for_each = advanced_filter.value.number_greater_than
        content {
          key   = number_greater_than.value.key
          value = number_greater_than.value.value
        }
      }
      dynamic "number_greater_than_or_equals" {
        for_each = advanced_filter.value.number_greater_than_or_equals
        content {
          key   = number_greater_than_or_equals.value.key
          value = number_greater_than_or_equals.value.value
        }
      }
      dynamic "number_less_than" {
        for_each = advanced_filter.value.number_less_than
        content {
          key   = number_less_than.value.key
          value = number_less_than.value.value
        }
      }
      dynamic "number_less_than_or_equals" {
        for_each = advanced_filter.value.number_less_than_or_equals
        content {
          key   = number_less_than_or_equals.value.key
          value = number_less_than_or_equals.value.value
        }
      }
      dynamic "number_in" {
        for_each = advanced_filter.value.number_in
        content {
          key    = number_in.value.key
          values = number_in.value.values
        }
      }
      dynamic "number_not_in" {
        for_each = advanced_filter.value.number_not_in
        content {
          key    = number_not_in.value.key
          values = number_not_in.value.values
        }
      }
      # The provider's range shape is a list of [from, to] pairs; the
      # spec's named-message shape renders to it here.
      dynamic "number_in_range" {
        for_each = advanced_filter.value.number_in_range
        content {
          key    = number_in_range.value.key
          values = [for r in number_in_range.value.ranges : [r.from, r.to]]
        }
      }
      dynamic "number_not_in_range" {
        for_each = advanced_filter.value.number_not_in_range
        content {
          key    = number_not_in_range.value.key
          values = [for r in number_not_in_range.value.ranges : [r.from, r.to]]
        }
      }
      dynamic "string_begins_with" {
        for_each = advanced_filter.value.string_begins_with
        content {
          key    = string_begins_with.value.key
          values = string_begins_with.value.values
        }
      }
      dynamic "string_not_begins_with" {
        for_each = advanced_filter.value.string_not_begins_with
        content {
          key    = string_not_begins_with.value.key
          values = string_not_begins_with.value.values
        }
      }
      dynamic "string_ends_with" {
        for_each = advanced_filter.value.string_ends_with
        content {
          key    = string_ends_with.value.key
          values = string_ends_with.value.values
        }
      }
      dynamic "string_not_ends_with" {
        for_each = advanced_filter.value.string_not_ends_with
        content {
          key    = string_not_ends_with.value.key
          values = string_not_ends_with.value.values
        }
      }
      dynamic "string_contains" {
        for_each = advanced_filter.value.string_contains
        content {
          key    = string_contains.value.key
          values = string_contains.value.values
        }
      }
      dynamic "string_not_contains" {
        for_each = advanced_filter.value.string_not_contains
        content {
          key    = string_not_contains.value.key
          values = string_not_contains.value.values
        }
      }
      dynamic "string_in" {
        for_each = advanced_filter.value.string_in
        content {
          key    = string_in.value.key
          values = string_in.value.values
        }
      }
      dynamic "string_not_in" {
        for_each = advanced_filter.value.string_not_in
        content {
          key    = string_not_in.value.key
          values = string_not_in.value.values
        }
      }
      dynamic "is_not_null" {
        for_each = advanced_filter.value.is_not_null
        content {
          key = is_not_null.value.key
        }
      }
      dynamic "is_null_or_undefined" {
        for_each = advanced_filter.value.is_null_or_undefined
        content {
          key = is_null_or_undefined.value.key
        }
      }
    }
  }
}

# A system-topic subscription: a child of the system topic, addressed
# by (resource group, system topic name) parsed from the referenced
# topic's ARM id in locals.tf. The configuration body is identical to
# the scope-addressed resource above by design.
resource "azurerm_eventgrid_system_topic_event_subscription" "main" {
  count = var.spec.system_topic_id != null ? 1 : 0

  name                = var.spec.name
  system_topic        = local.system_topic_name
  resource_group_name = local.system_topic_resource_group

  # Always sent (platform default mirrors Azure's). Create-only.
  event_delivery_schema = var.spec.event_delivery_schema

  expiration_time_utc = var.spec.expiration_time_utc != "" ? var.spec.expiration_time_utc : null

  included_event_types = length(var.spec.included_event_types) > 0 ? var.spec.included_event_types : null

  labels = var.spec.labels

  advanced_filtering_on_arrays_enabled = var.spec.advanced_filtering_on_arrays_enabled

  eventhub_id          = var.spec.destination.eventhub_id
  hybrid_connection_id = var.spec.destination.hybrid_connection_id
  service_bus_queue_id = var.spec.destination.service_bus_queue_id
  service_bus_topic_id = var.spec.destination.service_bus_topic_id

  dynamic "azure_function_endpoint" {
    for_each = var.spec.destination.azure_function != null ? [var.spec.destination.azure_function] : []
    content {
      function_id                       = azure_function_endpoint.value.function_id
      max_events_per_batch              = azure_function_endpoint.value.max_events_per_batch
      preferred_batch_size_in_kilobytes = azure_function_endpoint.value.preferred_batch_size_in_kilobytes
    }
  }

  dynamic "storage_queue_endpoint" {
    for_each = var.spec.destination.storage_queue != null ? [var.spec.destination.storage_queue] : []
    content {
      storage_account_id                    = storage_queue_endpoint.value.storage_account_id
      queue_name                            = storage_queue_endpoint.value.queue_name
      queue_message_time_to_live_in_seconds = storage_queue_endpoint.value.queue_message_time_to_live_in_seconds
    }
  }

  dynamic "webhook_endpoint" {
    for_each = var.spec.destination.webhook != null ? [var.spec.destination.webhook] : []
    content {
      url                               = webhook_endpoint.value.url
      max_events_per_batch              = webhook_endpoint.value.max_events_per_batch
      preferred_batch_size_in_kilobytes = webhook_endpoint.value.preferred_batch_size_in_kilobytes
      active_directory_tenant_id        = webhook_endpoint.value.active_directory_tenant_id != "" ? webhook_endpoint.value.active_directory_tenant_id : null
      active_directory_app_id_or_uri    = webhook_endpoint.value.active_directory_app_id_or_uri != "" ? webhook_endpoint.value.active_directory_app_id_or_uri : null
    }
  }

  dynamic "delivery_identity" {
    for_each = var.spec.delivery_identity != null ? [var.spec.delivery_identity] : []
    content {
      type                   = local.identity_type_map[delivery_identity.value.type]
      user_assigned_identity = delivery_identity.value.user_assigned_identity != "" ? delivery_identity.value.user_assigned_identity : null
    }
  }

  dynamic "delivery_property" {
    for_each = var.spec.delivery_properties
    content {
      header_name  = delivery_property.value.header_name
      type         = delivery_property.value.type
      value        = delivery_property.value.value != "" ? delivery_property.value.value : null
      source_field = delivery_property.value.source_field != "" ? delivery_property.value.source_field : null
      secret       = delivery_property.value.secret
    }
  }

  dynamic "storage_blob_dead_letter_destination" {
    for_each = var.spec.dead_letter != null ? [var.spec.dead_letter] : []
    content {
      storage_account_id          = storage_blob_dead_letter_destination.value.storage_account_id
      storage_blob_container_name = storage_blob_dead_letter_destination.value.storage_blob_container_name
    }
  }

  dynamic "dead_letter_identity" {
    for_each = var.spec.dead_letter_identity != null ? [var.spec.dead_letter_identity] : []
    content {
      type                   = local.identity_type_map[dead_letter_identity.value.type]
      user_assigned_identity = dead_letter_identity.value.user_assigned_identity != "" ? dead_letter_identity.value.user_assigned_identity : null
    }
  }

  dynamic "retry_policy" {
    for_each = var.spec.retry_policy != null ? [var.spec.retry_policy] : []
    content {
      max_delivery_attempts = retry_policy.value.max_delivery_attempts
      event_time_to_live    = retry_policy.value.event_time_to_live
    }
  }

  dynamic "subject_filter" {
    for_each = var.spec.subject_filter != null ? [var.spec.subject_filter] : []
    content {
      subject_begins_with = subject_filter.value.subject_begins_with != "" ? subject_filter.value.subject_begins_with : null
      subject_ends_with   = subject_filter.value.subject_ends_with != "" ? subject_filter.value.subject_ends_with : null
      case_sensitive      = subject_filter.value.case_sensitive
    }
  }

  dynamic "advanced_filter" {
    for_each = var.spec.advanced_filter != null ? [var.spec.advanced_filter] : []
    content {
      dynamic "bool_equals" {
        for_each = advanced_filter.value.bool_equals
        content {
          key   = bool_equals.value.key
          value = bool_equals.value.value
        }
      }
      dynamic "number_greater_than" {
        for_each = advanced_filter.value.number_greater_than
        content {
          key   = number_greater_than.value.key
          value = number_greater_than.value.value
        }
      }
      dynamic "number_greater_than_or_equals" {
        for_each = advanced_filter.value.number_greater_than_or_equals
        content {
          key   = number_greater_than_or_equals.value.key
          value = number_greater_than_or_equals.value.value
        }
      }
      dynamic "number_less_than" {
        for_each = advanced_filter.value.number_less_than
        content {
          key   = number_less_than.value.key
          value = number_less_than.value.value
        }
      }
      dynamic "number_less_than_or_equals" {
        for_each = advanced_filter.value.number_less_than_or_equals
        content {
          key   = number_less_than_or_equals.value.key
          value = number_less_than_or_equals.value.value
        }
      }
      dynamic "number_in" {
        for_each = advanced_filter.value.number_in
        content {
          key    = number_in.value.key
          values = number_in.value.values
        }
      }
      dynamic "number_not_in" {
        for_each = advanced_filter.value.number_not_in
        content {
          key    = number_not_in.value.key
          values = number_not_in.value.values
        }
      }
      dynamic "number_in_range" {
        for_each = advanced_filter.value.number_in_range
        content {
          key    = number_in_range.value.key
          values = [for r in number_in_range.value.ranges : [r.from, r.to]]
        }
      }
      dynamic "number_not_in_range" {
        for_each = advanced_filter.value.number_not_in_range
        content {
          key    = number_not_in_range.value.key
          values = [for r in number_not_in_range.value.ranges : [r.from, r.to]]
        }
      }
      dynamic "string_begins_with" {
        for_each = advanced_filter.value.string_begins_with
        content {
          key    = string_begins_with.value.key
          values = string_begins_with.value.values
        }
      }
      dynamic "string_not_begins_with" {
        for_each = advanced_filter.value.string_not_begins_with
        content {
          key    = string_not_begins_with.value.key
          values = string_not_begins_with.value.values
        }
      }
      dynamic "string_ends_with" {
        for_each = advanced_filter.value.string_ends_with
        content {
          key    = string_ends_with.value.key
          values = string_ends_with.value.values
        }
      }
      dynamic "string_not_ends_with" {
        for_each = advanced_filter.value.string_not_ends_with
        content {
          key    = string_not_ends_with.value.key
          values = string_not_ends_with.value.values
        }
      }
      dynamic "string_contains" {
        for_each = advanced_filter.value.string_contains
        content {
          key    = string_contains.value.key
          values = string_contains.value.values
        }
      }
      dynamic "string_not_contains" {
        for_each = advanced_filter.value.string_not_contains
        content {
          key    = string_not_contains.value.key
          values = string_not_contains.value.values
        }
      }
      dynamic "string_in" {
        for_each = advanced_filter.value.string_in
        content {
          key    = string_in.value.key
          values = string_in.value.values
        }
      }
      dynamic "string_not_in" {
        for_each = advanced_filter.value.string_not_in
        content {
          key    = string_not_in.value.key
          values = string_not_in.value.values
        }
      }
      dynamic "is_not_null" {
        for_each = advanced_filter.value.is_not_null
        content {
          key = is_not_null.value.key
        }
      }
      dynamic "is_null_or_undefined" {
        for_each = advanced_filter.value.is_null_or_undefined
        content {
          key = is_null_or_undefined.value.key
        }
      }
    }
  }
}
