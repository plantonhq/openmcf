# Exactly one of the two resources exists (the addressing choice is
# spec-enforced to one), so the one() over the union is total.
output "event_subscription_id" {
  description = "The Azure Resource Manager ID of the event subscription"
  value = one(concat(
    azurerm_eventgrid_event_subscription.main[*].id,
    azurerm_eventgrid_system_topic_event_subscription.main[*].id,
  ))
}

output "event_subscription_name" {
  description = "The event subscription's name"
  value = one(concat(
    azurerm_eventgrid_event_subscription.main[*].name,
    azurerm_eventgrid_system_topic_event_subscription.main[*].name,
  ))
}
