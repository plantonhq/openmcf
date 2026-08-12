# Create one named event stream (domain topic) inside an Azure Event
# Grid domain -- the per-tenant mailbox of the multi-tenant pattern.
# Publishers address it by naming the topic in events sent to the
# DOMAIN's endpoint; subscribers attach event subscriptions to the
# topic's own ARM id ({domain_id}/topics/{name}). Free at rest;
# everything here is create-only (the topic is pure addressing).
resource "azurerm_eventgrid_domain_topic" "main" {
  name                = var.spec.name
  domain_name         = local.domain_name
  resource_group_name = local.resource_group_name
}
