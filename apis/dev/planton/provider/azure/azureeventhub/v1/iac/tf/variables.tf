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
  description = "Azure Event Hub specification"
  type = object({
    # The parent Event Hubs namespace ARM id. References are resolved to
    # a literal by the platform before the module runs. ForceNew.
    namespace_id = string

    # The hub's name, unique within the namespace (1-256 characters).
    # ForceNew.
    event_hub_name = string

    # The number of partitions -- the unit of parallelism and ordering.
    # Never decreasable; increasable only on PREMIUM or a dedicated
    # cluster (Azure enforces the tier caps at apply time).
    partition_count = number

    # Simple retention: how long events stay replayable, in days.
    # Exactly one of message_retention and retention_description is set
    # (spec-enforced XOR).
    message_retention = optional(number)

    # Rich retention: hour-granular windows and Kafka-style compaction.
    retention_description = optional(object({
      # DELETE (remove past the window) or COMPACT (keep the latest
      # event per partition key). ForceNew.
      cleanup_policy = string
      # The DELETE window in hours; absent with COMPACT.
      retention_time_in_hours = optional(number)
      # The COMPACT tombstone window in hours; absent with DELETE.
      tombstone_retention_time_in_hours = optional(number)
    }))

    # The gate state, as the spec enum's value name (ACTIVE, DISABLED,
    # SEND_DISABLED). Unset deploys Active.
    status = optional(string)

    # Capture: continuous archival of every event to Blob Storage.
    capture_description = optional(object({
      # Whether capture is running -- keeping the block with false
      # preserves the configuration while pausing archival.
      enabled = bool
      # AVRO or AVRO_DEFLATE.
      encoding = string
      # Window cadence in seconds (60-900); unset keeps Azure's 300.
      interval_in_seconds = optional(number)
      # Window size in bytes (10 MB - 500 MB); unset keeps Azure's
      # 300 MB.
      size_limit_in_bytes = optional(number)
      # Skip writing archive files for empty windows.
      skip_empty_archives = optional(bool)
      destination = object({
        # Must contain all nine {Namespace}/{EventHub}/{PartitionId}/
        # {Year}/{Month}/{Day}/{Hour}/{Minute}/{Second} tokens
        # (spec-enforced).
        archive_name_format = string
        # The blob container archives land in (resolved to a literal
        # name before the module runs).
        blob_container_name = string
        # The storage account ARM id (resolved before the module runs).
        storage_account_id = string
        # STORAGE_SAS (Azure's default), SYSTEM_ASSIGNED, or
        # USER_ASSIGNED.
        storage_authentication_type = optional(string)
        # The user-assigned identity ARM id -- required with
        # USER_ASSIGNED (spec-enforced pairing).
        storage_authentication_id = optional(string)
      })
    }))
  })
}
