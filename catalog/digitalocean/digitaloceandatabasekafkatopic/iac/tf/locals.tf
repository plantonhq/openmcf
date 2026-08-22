locals {
  # The spec's config message (null when absent; the locals below are only
  # consumed inside the dynamic block, which does not expand in that case --
  # the try() guards keep evaluation safe either way).
  config = var.spec.config

  # Enum-walled string leaves: empty means unset; the provider must not
  # receive "" (it would fail the provider-side value validation).
  config_strings = {
    cleanup_policy         = try(local.config.cleanup_policy, "") != "" ? local.config.cleanup_policy : null
    compression_type       = try(local.config.compression_type, "") != "" ? local.config.compression_type : null
    message_format_version = try(local.config.message_format_version, "") != "" ? local.config.message_format_version : null
    message_timestamp_type = try(local.config.message_timestamp_type, "") != "" ? local.config.message_timestamp_type : null
  }

  # 64-bit numeric leaves: the provider's schema carries them as strings
  # (Terraform numbers are not 64-bit safe), so render each present value
  # to its decimal string and keep absent values null (never sent).
  config_numbers = {
    delete_retention_ms                 = try(local.config.delete_retention_ms, null) != null ? tostring(local.config.delete_retention_ms) : null
    file_delete_delay_ms                = try(local.config.file_delete_delay_ms, null) != null ? tostring(local.config.file_delete_delay_ms) : null
    flush_messages                      = try(local.config.flush_messages, null) != null ? tostring(local.config.flush_messages) : null
    flush_ms                            = try(local.config.flush_ms, null) != null ? tostring(local.config.flush_ms) : null
    index_interval_bytes                = try(local.config.index_interval_bytes, null) != null ? tostring(local.config.index_interval_bytes) : null
    max_compaction_lag_ms               = try(local.config.max_compaction_lag_ms, null) != null ? tostring(local.config.max_compaction_lag_ms) : null
    max_message_bytes                   = try(local.config.max_message_bytes, null) != null ? tostring(local.config.max_message_bytes) : null
    message_timestamp_difference_max_ms = try(local.config.message_timestamp_difference_max_ms, null) != null ? tostring(local.config.message_timestamp_difference_max_ms) : null
    min_compaction_lag_ms               = try(local.config.min_compaction_lag_ms, null) != null ? tostring(local.config.min_compaction_lag_ms) : null
    retention_bytes                     = try(local.config.retention_bytes, null) != null ? tostring(local.config.retention_bytes) : null
    retention_ms                        = try(local.config.retention_ms, null) != null ? tostring(local.config.retention_ms) : null
    segment_bytes                       = try(local.config.segment_bytes, null) != null ? tostring(local.config.segment_bytes) : null
    segment_index_bytes                 = try(local.config.segment_index_bytes, null) != null ? tostring(local.config.segment_index_bytes) : null
    segment_jitter_ms                   = try(local.config.segment_jitter_ms, null) != null ? tostring(local.config.segment_jitter_ms) : null
    segment_ms                          = try(local.config.segment_ms, null) != null ? tostring(local.config.segment_ms) : null
  }
}
