# Create the Event Grid namespace topic -- one named CloudEvents
# stream inside a namespace. Azure pins the event schema to CloudEvents
# v1.0 and the publisher type to Custom (the provider sends both;
# neither is configurable). Retention is the topic's ONLY updatable
# property; name and namespace replace it.
resource "azurerm_eventgrid_namespace_topic" "main" {
  name                   = local.namespace_topic_name
  eventgrid_namespace_id = var.spec.namespace_id

  # Platform default 7 days (the provider's own) -- always sent, so the
  # rendered plan states the retention.
  event_retention_in_days = coalesce(var.spec.event_retention_in_days, 7)
}
