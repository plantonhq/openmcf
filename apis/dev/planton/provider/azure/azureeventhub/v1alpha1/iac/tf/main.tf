# One partitioned, replayable event stream inside an Event Hubs
# namespace. Consumer groups, SAS rules, and data-plane role assignments
# are first-class kinds that reference this hub's ARM id -- nothing is
# bundled here.
resource "azurerm_eventhub" "main" {
  # ForceNew: renaming replaces the hub and its retained events. The
  # name is also the Kafka topic name on the namespace's Kafka endpoint.
  name = var.spec.event_hub_name

  # ForceNew: the hub cannot move between namespaces.
  namespace_id = var.spec.namespace_id

  # The unit of parallelism and ordering. Azure enforces the tier caps
  # (32 on shared namespaces, 1024 on PREMIUM/dedicated) and the
  # never-decrease / increase-only-on-PREMIUM-or-dedicated contracts at
  # apply time -- they depend on the parent namespace's tier, which this
  # module cannot see.
  partition_count = var.spec.partition_count

  # Simple retention in days. Exactly one of message_retention and
  # retention_description is set (spec-enforced XOR); Azure caps the
  # window by tier (1/7/90 days) at apply time.
  message_retention = var.spec.message_retention

  # Rich retention: hour-granular windows and Kafka-style compaction.
  # The spec pairs the hour field to the policy (DELETE takes
  # retention_time_in_hours, COMPACT takes the tombstone window), so
  # each variant sends only its own field -- Azure silently ignores the
  # mismatched one, which the spec rejects up front instead.
  dynamic "retention_description" {
    for_each = var.spec.retention_description != null ? [var.spec.retention_description] : []
    content {
      # ForceNew: the cleanup policy is fixed at creation.
      cleanup_policy                    = local.cleanup_policy_map[retention_description.value.cleanup_policy]
      retention_time_in_hours           = retention_description.value.retention_time_in_hours
      tombstone_retention_time_in_hours = retention_description.value.tombstone_retention_time_in_hours
    }
  }

  # The administrative gate: Active, Disabled (sends and receives
  # rejected; events retained), or SendDisabled (receive-only drain).
  status = local.status

  # Capture: the built-in streaming-to-batch bridge -- every event is
  # archived to Blob Storage in Avro on a size-or-interval cadence, with
  # no consumer application to run.
  dynamic "capture_description" {
    for_each = var.spec.capture_description != null ? [var.spec.capture_description] : []
    content {
      enabled             = capture_description.value.enabled
      encoding            = local.capture_encoding_map[capture_description.value.encoding]
      interval_in_seconds = capture_description.value.interval_in_seconds
      size_limit_in_bytes = capture_description.value.size_limit_in_bytes
      skip_empty_archives = capture_description.value.skip_empty_archives

      destination {
        # Azure's destination name accepts exactly one value (Blob
        # Storage; the Data Lake variant retired with Gen1) -- a
        # constant, not configuration, so the module sends it
        # unconditionally.
        name                = "EventHubArchive.AzureBlockBlob"
        archive_name_format = capture_description.value.destination.archive_name_format
        blob_container_name = capture_description.value.destination.blob_container_name
        storage_account_id  = capture_description.value.destination.storage_account_id

        # "StorageSAS" (Azure's default) means service-managed SAS --
        # the provider sends no identity for it. The identity paths are
        # keyless: grant the chosen identity Storage Blob Data
        # Contributor on the account and attach it via the namespace's
        # identity block.
        storage_authentication_type = local.capture_auth
        storage_authentication_id   = capture_description.value.destination.storage_authentication_id
      }
    }
  }
}
