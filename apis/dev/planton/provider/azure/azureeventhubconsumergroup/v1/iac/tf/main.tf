# The consumer group: one application's independent read cursor over
# the hub's partitions. Tier limits are enforced by Azure at apply time
# (BASIC hubs allow no additional groups beyond the service-created
# $Default; STANDARD allows 20 per hub), so quota errors surface
# verbatim from Azure.
resource "azurerm_eventhub_consumer_group" "main" {
  # ForceNew: the group is its name -- renaming replaces it and resets
  # its consumers' stored offsets.
  name = var.spec.consumer_group_name

  # Parent addressing by discrete names, derived from the spec's single
  # event hub ARM ID in locals.
  resource_group_name = local.resource_group_name
  namespace_name      = local.namespace_name
  eventhub_name       = local.event_hub_name

  user_metadata = var.spec.user_metadata
}
