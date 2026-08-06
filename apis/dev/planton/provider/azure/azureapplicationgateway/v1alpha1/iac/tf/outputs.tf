output "application_gateway_id" {
  description = "The ARM ID of the Application Gateway"
  value       = azurerm_application_gateway.main.id
}

output "application_gateway_name" {
  description = "The name of the Application Gateway"
  value       = azurerm_application_gateway.main.name
}

output "backend_address_pool_ids" {
  description = "The ARM ID of each backend address pool, keyed by pool name -- the membership seam NICs and scale sets join through"
  value       = { for pool in azurerm_application_gateway.main.backend_address_pool : pool.name => pool.id }
}

output "frontend_ip_configuration_ids" {
  description = "The ARM ID of each frontend IP configuration, keyed by frontend name"
  value       = { for frontend in azurerm_application_gateway.main.frontend_ip_configuration : frontend.name => frontend.id }
}

output "private_ip_address" {
  description = "The first private frontend's address (what internal DNS records point at); empty when every frontend is public"
  value = try([
    for frontend in azurerm_application_gateway.main.frontend_ip_configuration : frontend.private_ip_address
    if frontend.private_ip_address != null && frontend.private_ip_address != ""
  ][0], "")
}

output "private_ip_addresses" {
  description = "The private addresses of ALL private frontends, in declaration order"
  value = [
    for frontend in azurerm_application_gateway.main.frontend_ip_configuration : frontend.private_ip_address
    if frontend.private_ip_address != null && frontend.private_ip_address != ""
  ]
}
