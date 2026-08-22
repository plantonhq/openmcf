# DigitalOcean Database Kafka Topic
#
# Provisions a topic on a DigitalOcean managed Kafka cluster -- the complete
# digitalocean_database_kafka_topic resource surface, including the full
# per-topic configuration block. Topic creation is asynchronous on
# DigitalOcean's side; partition count, replication factor, and every config
# leaf update in place (Kafka only ever ADDS partitions -- lowering the
# count is rejected by the API).
#
# The provider carries the numeric config tunables as strings because
# Terraform's number type cannot hold 64-bit values; the spec models them as
# real integers and this module renders them to strings at the boundary.
# Every leaf is presence-gated: an unset leaf is never sent, so the Kafka
# server default applies.

resource "digitalocean_database_kafka_topic" "topic" {
  cluster_id = var.spec.cluster
  name       = var.spec.topic_name

  # Unset lets DigitalOcean create the API default of 3 partitions. The
  # provider never reads the live count back (partition changes apply
  # asynchronously), so this is deliberately config-authoritative.
  partition_count    = var.spec.partition_count
  replication_factor = var.spec.replication_factor

  # Per-topic Kafka configuration. NOTE: when a config block is present the
  # provider seeds cleanup_policy to "compact_delete" unless set explicitly.
  dynamic "config" {
    for_each = var.spec.config != null ? [var.spec.config] : []
    content {
      cleanup_policy   = local.config_strings.cleanup_policy
      compression_type = local.config_strings.compression_type

      delete_retention_ms   = local.config_numbers.delete_retention_ms
      file_delete_delay_ms  = local.config_numbers.file_delete_delay_ms
      flush_messages        = local.config_numbers.flush_messages
      flush_ms              = local.config_numbers.flush_ms
      index_interval_bytes  = local.config_numbers.index_interval_bytes
      max_compaction_lag_ms = local.config_numbers.max_compaction_lag_ms
      max_message_bytes     = local.config_numbers.max_message_bytes

      message_down_conversion_enable = config.value.message_down_conversion_enable

      message_format_version              = local.config_strings.message_format_version
      message_timestamp_difference_max_ms = local.config_numbers.message_timestamp_difference_max_ms
      message_timestamp_type              = local.config_strings.message_timestamp_type

      min_cleanable_dirty_ratio = config.value.min_cleanable_dirty_ratio
      min_compaction_lag_ms     = local.config_numbers.min_compaction_lag_ms

      # The block's only leaf the provider defaults locally (to 1) instead
      # of reading the server value back.
      min_insync_replicas = config.value.min_insync_replicas

      preallocate = config.value.preallocate

      retention_bytes = local.config_numbers.retention_bytes
      retention_ms    = local.config_numbers.retention_ms

      segment_bytes       = local.config_numbers.segment_bytes
      segment_index_bytes = local.config_numbers.segment_index_bytes
      segment_jitter_ms   = local.config_numbers.segment_jitter_ms
      segment_ms          = local.config_numbers.segment_ms
    }
  }
}
