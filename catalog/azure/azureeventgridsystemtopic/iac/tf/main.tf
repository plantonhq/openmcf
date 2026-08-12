# Create the Azure Event Grid system topic -- the subscription surface
# for events Azure itself publishes about the source resource. Azure
# allows ONE system topic per source resource per topic type; the
# region must match the source's (or be "Global" for global sources).
# Free at rest; billing is per operation.
resource "azurerm_eventgrid_system_topic" "main" {
  name                = var.spec.name
  resource_group_name = var.spec.resource_group
  location            = var.spec.region

  # Both create-only: the (source, type) pair IS the topic's identity.
  source_resource_id = var.spec.source_resource_id
  topic_type         = var.spec.topic_type

  dynamic "identity" {
    for_each = var.spec.identity != null ? [var.spec.identity] : []
    content {
      type         = local.identity_type_map[identity.value.type]
      identity_ids = length(identity.value.identity_ids) > 0 ? identity.value.identity_ids : null
    }
  }

  tags = local.final_tags
}
