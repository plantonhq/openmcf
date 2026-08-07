output "firewall_id" {
  description = "The Azure Resource Manager ID of the firewall"
  value       = azurerm_firewall.main.id
}

output "firewall_name" {
  description = "The name of the firewall resource"
  value       = azurerm_firewall.main.name
}

# THE hub-spoke seam: spoke route tables send egress here via a
# VIRTUAL_APPLIANCE next hop. Azure computes it on the subnet-bearing ip
# configuration; empty for hub-deployed firewalls.
output "private_ip_address" {
  description = "The firewall's private IP on its data-path subnet"
  value = try(
    [for c in azurerm_firewall.main.ip_configuration : c.private_ip_address if c.private_ip_address != null && c.private_ip_address != ""][0],
    ""
  )
}

output "management_private_ip_address" {
  description = "The management path's private IP, when a management configuration is deployed"
  value = try(
    azurerm_firewall.main.management_ip_configuration[0].private_ip_address,
    ""
  )
}

# Hub-firewall addressing is Azure-assigned and only known after
# deployment -- exported for DNS records and route configuration. Empty
# for VNet firewalls.
output "virtual_hub_public_ip_addresses" {
  description = "The public IPs Azure assigned to a hub-deployed firewall"
  value = try(
    azurerm_firewall.main.virtual_hub[0].public_ip_addresses,
    []
  )
}

output "virtual_hub_private_ip_address" {
  description = "The private IP Azure assigned to a hub-deployed firewall"
  value = try(
    azurerm_firewall.main.virtual_hub[0].private_ip_address,
    ""
  )
}
