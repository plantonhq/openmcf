# The queue, addressed by the parent namespace's ARM ID (azurerm's v4
# child-addressing grain). Premium-tier contracts the provider enforces
# at apply time -- express rejected on Premium, partitioning must match
# the namespace's partition layout, large messages Premium-only -- are
# documented on the spec fields; the module passes the dials through
# unchanged so those contracts surface verbatim from Azure.
resource "azurerm_servicebus_queue" "main" {
  # ForceNew: renaming replaces the queue and drops any messages in it.
  name         = var.spec.queue_name
  namespace_id = var.spec.namespace_id

  # Capacity dials. Unset sizes let Azure default for the namespace's
  # tier (1024 MB multi-tenant, 81920 MB premium); for partitioned
  # multi-tenant queues Azure reports 16x the sold size and the provider
  # normalizes the read back.
  max_size_in_megabytes         = var.spec.max_size_in_megabytes
  max_message_size_in_kilobytes = var.spec.max_message_size_in_kilobytes

  # ForceNew trio: the storage layout and dedup/session models are fixed
  # at creation.
  partitioning_enabled         = var.spec.partitioning_enabled
  requires_duplicate_detection = var.spec.requires_duplicate_detection
  requires_session             = var.spec.requires_session

  # Lifecycle dials -- unset leaves Azure's defaults in place (TTL
  # unbounded, dedup window PT10M, lock PT1M, 10 deliveries).
  default_message_ttl                     = var.spec.default_message_ttl
  duplicate_detection_history_time_window = var.spec.duplicate_detection_history_time_window
  lock_duration                           = var.spec.lock_duration
  max_delivery_count                      = var.spec.max_delivery_count
  dead_lettering_on_message_expiration    = var.spec.dead_lettering_on_message_expiration
  auto_delete_on_idle                     = var.spec.auto_delete_on_idle

  batched_operations_enabled = var.spec.batched_operations_enabled
  express_enabled            = var.spec.express_enabled

  # Routing chains: targets are entity NAMES in the same namespace (not
  # ARM ids) -- Azure's own addressing for auto-forwarding. The target
  # must exist first; compose with depends-on ordering in charts.
  forward_to                        = var.spec.forward_to
  forward_dead_lettered_messages_to = var.spec.forward_dead_lettered_messages_to

  status = local.status
}
