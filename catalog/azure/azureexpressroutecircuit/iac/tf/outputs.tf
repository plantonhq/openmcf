output "express_route_circuit_id" {
  description = "The Azure Resource Manager ID of the ExpressRoute circuit"
  value       = azurerm_express_route_circuit.main.id
}

output "express_route_circuit_name" {
  description = "The name of the circuit -- what peerings reference as express_route_circuit_name"
  value       = azurerm_express_route_circuit.main.name
}

output "service_key" {
  description = "The circuit's provisioning credential, handed to the connectivity provider to complete the cross-connect"
  value       = azurerm_express_route_circuit.main.service_key
  sensitive   = true
}

output "service_provider_provisioning_state" {
  description = "The provider side's provisioning state: NotProvisioned, Provisioning, Provisioned, or Deprovisioning"
  value       = azurerm_express_route_circuit.main.service_provider_provisioning_state
}

output "authorization_keys" {
  description = "The generated key of each authorization issued by this circuit, keyed by the authorization's name"
  value       = { for name, authorization in azurerm_express_route_circuit_authorization.authorizations : name => authorization.authorization_key }
  sensitive   = true
}
