output "express_route_port_id" {
  description = "The Azure Resource Manager ID of the ExpressRoute Port -- what an ExpressRoute Direct circuit references"
  value       = azurerm_express_route_port.main.id
}

output "express_route_port_name" {
  description = "The name of the port"
  value       = azurerm_express_route_port.main.name
}

output "guid" {
  description = "The port's globally unique resource GUID"
  value       = azurerm_express_route_port.main.guid
}

output "ethertype" {
  description = "The link ethertype, fixed by the encapsulation choice"
  value       = azurerm_express_route_port.main.ethertype
}

output "mtu" {
  description = "The maximum transmission unit of the links"
  value       = azurerm_express_route_port.main.mtu
}

output "system_assigned_identity_principal_id" {
  description = "The principal ID of the port's system-assigned identity (populated only when the identity type includes SYSTEM_ASSIGNED)"
  value       = try(azurerm_express_route_port.main.identity[0].principal_id, "")
}

# The per-link facility facts (router, interface, patch panel, rack) are
# the letter-of-authorization data handed to the colocation facility to
# order the physical cross-connects.
output "link1_id" {
  description = "The ARM ID of the first physical link"
  value       = try(azurerm_express_route_port.main.link1[0].id, "")
}

output "link1_router_name" {
  description = "The Microsoft edge router of the first link"
  value       = try(azurerm_express_route_port.main.link1[0].router_name, "")
}

output "link1_interface_name" {
  description = "The router interface of the first link"
  value       = try(azurerm_express_route_port.main.link1[0].interface_name, "")
}

output "link1_patch_panel_id" {
  description = "The patch panel the first link lands on"
  value       = try(azurerm_express_route_port.main.link1[0].patch_panel_id, "")
}

output "link1_rack_id" {
  description = "The rack of the first link's patch panel"
  value       = try(azurerm_express_route_port.main.link1[0].rack_id, "")
}

output "link1_connector_type" {
  description = "The physical connector type of the first link"
  value       = try(azurerm_express_route_port.main.link1[0].connector_type, "")
}

output "link2_id" {
  description = "The ARM ID of the second physical link"
  value       = try(azurerm_express_route_port.main.link2[0].id, "")
}

output "link2_router_name" {
  description = "The Microsoft edge router of the second link"
  value       = try(azurerm_express_route_port.main.link2[0].router_name, "")
}

output "link2_interface_name" {
  description = "The router interface of the second link"
  value       = try(azurerm_express_route_port.main.link2[0].interface_name, "")
}

output "link2_patch_panel_id" {
  description = "The patch panel the second link lands on"
  value       = try(azurerm_express_route_port.main.link2[0].patch_panel_id, "")
}

output "link2_rack_id" {
  description = "The rack of the second link's patch panel"
  value       = try(azurerm_express_route_port.main.link2[0].rack_id, "")
}

output "link2_connector_type" {
  description = "The physical connector type of the second link"
  value       = try(azurerm_express_route_port.main.link2[0].connector_type, "")
}

output "authorization_keys" {
  description = "The generated key of each authorization issued by this port, keyed by the authorization's name"
  value       = { for name, authorization in azurerm_express_route_port_authorization.authorizations : name => authorization.authorization_key }
  sensitive   = true
}
