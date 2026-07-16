output "load_balancer_id" {
  description = "The Azure Resource Manager ID of the load balancer"
  value       = azurerm_lb.main.id
}

output "load_balancer_name" {
  description = "The name of the load balancer"
  value       = azurerm_lb.main.name
}

output "private_ip_address" {
  description = "The first internal frontend's private IP address (empty when every frontend is public)"
  value       = azurerm_lb.main.private_ip_address == null ? "" : azurerm_lb.main.private_ip_address
}

output "private_ip_addresses" {
  description = "The private IP addresses of all internal frontends, in declaration order"
  value       = azurerm_lb.main.private_ip_addresses == null ? [] : azurerm_lb.main.private_ip_addresses
}

output "frontend_ip_configuration_ids" {
  description = "The ARM ID of each frontend IP configuration, keyed by frontend name"
  value       = { for f in azurerm_lb.main.frontend_ip_configuration : f.name => f.id }
}

output "backend_pool_ids" {
  description = "The ARM ID of each backend address pool, keyed by pool name -- the member-side association seam"
  value       = { for name, pool in azurerm_lb_backend_address_pool.pools : name => pool.id }
}

output "probe_ids" {
  description = "The ARM ID of each health probe, keyed by probe name"
  value       = { for name, probe in azurerm_lb_probe.probes : name => probe.id }
}

output "nat_rule_ids" {
  description = "The ARM ID of each inbound NAT rule, keyed by rule name -- the NIC NAT-rule association seam"
  value       = { for name, rule in azurerm_lb_nat_rule.nat_rules : name => rule.id }
}
