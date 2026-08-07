# The topic, addressed by the parent namespace's ARM ID (azurerm's v4
# child-addressing grain). Topics require STANDARD or PREMIUM (BASIC is
# queue-only -- Azure rejects the create). Premium-tier contracts the
# provider enforces at apply time -- partitioning must match the
# namespace's partition layout, large messages Premium-only -- are
# documented on the spec fields; the module passes the dials through
# unchanged so those contracts surface verbatim from Azure.
resource "azurerm_servicebus_topic" "main" {
  # ForceNew: renaming replaces the topic and every subscription under it.
  name         = var.spec.topic_name
  namespace_id = var.spec.namespace_id

  # Capacity dials. Unset sizes let Azure default for the namespace's
  # tier; for partitioned STANDARD topics Azure reports 16x the sold
  # size and the provider normalizes the read back.
  max_size_in_megabytes         = var.spec.max_size_in_megabytes
  max_message_size_in_kilobytes = var.spec.max_message_size_in_kilobytes

  # ForceNew pair: the storage layout and dedup model are fixed at
  # creation.
  partitioning_enabled         = var.spec.partitioning_enabled
  requires_duplicate_detection = var.spec.requires_duplicate_detection

  # Lifecycle dials -- unset leaves Azure's defaults in place (TTL
  # unbounded, dedup window PT10M).
  default_message_ttl                     = var.spec.default_message_ttl
  duplicate_detection_history_time_window = var.spec.duplicate_detection_history_time_window
  auto_delete_on_idle                     = var.spec.auto_delete_on_idle

  batched_operations_enabled = var.spec.batched_operations_enabled
  express_enabled            = var.spec.express_enabled

  # Publish-order preservation -- pair with session-aware subscriptions
  # for strictly-ordered publish-subscribe.
  support_ordering = var.spec.support_ordering

  status = local.status
}
