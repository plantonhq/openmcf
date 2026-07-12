# The configuration has no ARM object of its own -- the provider's
# resource id IS the namespace's ARM id.
output "customer_managed_key_id" {
  description = "The provider's identity for the CMK configuration (the namespace's ARM ID)"
  value       = azurerm_eventhub_namespace_customer_managed_key.main.id
}
