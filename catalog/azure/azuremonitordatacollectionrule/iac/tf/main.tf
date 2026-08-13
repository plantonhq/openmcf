# Create the Azure Monitor data collection rule -- the routing table
# declaring what telemetry the Azure Monitor Agent collects (data
# sources), where it lands (destinations), and how the two wire
# together (data flows, optionally with an ingestion-time KQL
# transformation). Machines attach to the rule with separate
# association resources; the rule itself is free at rest.
#
# Names wire the rule together: flows reference destinations by their
# rule-local names, and destination names share one namespace across
# all arms -- Azure enforces both at deploy time. Platform
# compatibility (a Linux rule cannot carry windows_event_log sources,
# the *_direct destinations require kind AgentDirectToStore) is also
# enforced by Azure at deploy time; the provider performs no early
# check.
resource "azurerm_monitor_data_collection_rule" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  # Omitted when unset -- the default rule kind accepts every
  # platform's sources. Once set, changing (or clearing) the kind
  # forces a new rule (provider lifecycle).
  kind = var.spec.kind

  # Sent only when non-empty for a clean plan; Azure treats an absent
  # and an empty description identically.
  description = var.spec.description != "" ? var.spec.description : null

  # The DCE the rule ingests through -- required by Azure when the rule
  # declares custom streams; sent only when set.
  data_collection_endpoint_id = var.spec.data_collection_endpoint_id

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  dynamic "data_sources" {
    for_each = var.spec.data_sources != null ? [var.spec.data_sources] : []
    content {
      dynamic "syslog" {
        for_each = data_sources.value.syslogs
        content {
          name           = syslog.value.name
          facility_names = syslog.value.facility_names
          log_levels     = syslog.value.log_levels
          streams        = syslog.value.streams
        }
      }

      dynamic "performance_counter" {
        for_each = data_sources.value.performance_counters
        content {
          name                          = performance_counter.value.name
          sampling_frequency_in_seconds = performance_counter.value.sampling_frequency_in_seconds
          counter_specifiers            = performance_counter.value.counter_specifiers
          streams                       = performance_counter.value.streams
        }
      }

      dynamic "windows_event_log" {
        for_each = data_sources.value.windows_event_logs
        content {
          name           = windows_event_log.value.name
          x_path_queries = windows_event_log.value.x_path_queries
          streams        = windows_event_log.value.streams
        }
      }

      dynamic "extension" {
        for_each = data_sources.value.extensions
        content {
          name           = extension.value.name
          extension_name = extension.value.extension_name
          # The provider validates a non-empty JSON string -- sent only
          # when set.
          extension_json     = extension.value.extension_json != "" ? extension.value.extension_json : null
          input_data_sources = length(extension.value.input_data_sources) > 0 ? extension.value.input_data_sources : null
          streams            = extension.value.streams
        }
      }

      dynamic "iis_log" {
        for_each = data_sources.value.iis_logs
        content {
          name = iis_log.value.name
          # Omitted when empty -- the agent then reads the server's
          # configured IIS log location.
          log_directories = length(iis_log.value.log_directories) > 0 ? iis_log.value.log_directories : null
          streams         = iis_log.value.streams
        }
      }

      dynamic "log_file" {
        for_each = data_sources.value.log_files
        content {
          name          = log_file.value.name
          file_patterns = log_file.value.file_patterns
          format        = log_file.value.format
          streams       = log_file.value.streams

          dynamic "settings" {
            for_each = log_file.value.settings != null ? [log_file.value.settings] : []
            content {
              text {
                record_start_timestamp_format = settings.value.text.record_start_timestamp_format
              }
            }
          }
        }
      }

      dynamic "prometheus_forwarder" {
        for_each = data_sources.value.prometheus_forwarders
        content {
          name    = prometheus_forwarder.value.name
          streams = prometheus_forwarder.value.streams

          dynamic "label_include_filter" {
            for_each = prometheus_forwarder.value.label_include_filters
            content {
              label = label_include_filter.value.label
              value = label_include_filter.value.value
            }
          }
        }
      }

      dynamic "windows_firewall_log" {
        for_each = data_sources.value.windows_firewall_logs
        content {
          name    = windows_firewall_log.value.name
          streams = windows_firewall_log.value.streams
        }
      }

      dynamic "platform_telemetry" {
        for_each = data_sources.value.platform_telemetries
        content {
          name    = platform_telemetry.value.name
          streams = platform_telemetry.value.streams
        }
      }

      dynamic "data_import" {
        for_each = data_sources.value.data_import != null ? [data_sources.value.data_import] : []
        content {
          # Azure's rule model carries exactly ONE event-hub import; the
          # spec models it singular (the provider would silently drop
          # extra entries).
          event_hub_data_source {
            name   = data_import.value.event_hub_data_source.name
            stream = data_import.value.event_hub_data_source.stream
            # Sent only when set -- the provider omits an empty consumer
            # group (Azure then reads $Default).
            consumer_group = data_import.value.event_hub_data_source.consumer_group != "" ? data_import.value.event_hub_data_source.consumer_group : null
          }
        }
      }
    }
  }

  destinations {
    dynamic "log_analytics" {
      for_each = var.spec.destinations.log_analytics
      content {
        name                  = log_analytics.value.name
        workspace_resource_id = log_analytics.value.workspace_resource_id
      }
    }

    dynamic "azure_monitor_metrics" {
      for_each = var.spec.destinations.azure_monitor_metrics != null ? [var.spec.destinations.azure_monitor_metrics] : []
      content {
        name = azure_monitor_metrics.value.name
      }
    }

    dynamic "event_hub" {
      for_each = var.spec.destinations.event_hub != null ? [var.spec.destinations.event_hub] : []
      content {
        name         = event_hub.value.name
        event_hub_id = event_hub.value.event_hub_id
      }
    }

    dynamic "event_hub_direct" {
      for_each = var.spec.destinations.event_hub_direct != null ? [var.spec.destinations.event_hub_direct] : []
      content {
        name         = event_hub_direct.value.name
        event_hub_id = event_hub_direct.value.event_hub_id
      }
    }

    dynamic "monitor_account" {
      for_each = var.spec.destinations.monitor_accounts
      content {
        name               = monitor_account.value.name
        monitor_account_id = monitor_account.value.monitor_account_id
      }
    }

    dynamic "storage_blob" {
      for_each = var.spec.destinations.storage_blobs
      content {
        name               = storage_blob.value.name
        container_name     = storage_blob.value.container_name
        storage_account_id = storage_blob.value.storage_account_id
      }
    }

    dynamic "storage_blob_direct" {
      for_each = var.spec.destinations.storage_blob_directs
      content {
        name               = storage_blob_direct.value.name
        container_name     = storage_blob_direct.value.container_name
        storage_account_id = storage_blob_direct.value.storage_account_id
      }
    }

    dynamic "storage_table_direct" {
      for_each = var.spec.destinations.storage_table_directs
      content {
        name               = storage_table_direct.value.name
        table_name         = storage_table_direct.value.table_name
        storage_account_id = storage_table_direct.value.storage_account_id
      }
    }
  }

  dynamic "data_flow" {
    for_each = var.spec.data_flows
    content {
      streams      = data_flow.value.streams
      destinations = data_flow.value.destinations
      # The provider validates non-empty strings on all three -- each is
      # sent only when set.
      built_in_transform = data_flow.value.built_in_transform != "" ? data_flow.value.built_in_transform : null
      output_stream      = data_flow.value.output_stream != "" ? data_flow.value.output_stream : null
      transform_kql      = data_flow.value.transform_kql != "" ? data_flow.value.transform_kql : null
    }
  }

  dynamic "stream_declaration" {
    for_each = var.spec.stream_declarations
    content {
      stream_name = stream_declaration.value.stream_name

      dynamic "column" {
        for_each = stream_declaration.value.columns
        content {
          name = column.value.name
          type = column.value.type
        }
      }
    }
  }

  tags = local.final_tags
}
