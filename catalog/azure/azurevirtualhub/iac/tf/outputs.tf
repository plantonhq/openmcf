output "virtual_hub_id" {
  description = "The Azure Resource Manager ID of the hub -- what connections, gateways, and firewalls reference as their virtual_hub_id"
  value       = azurerm_virtual_hub.main.id
}

output "virtual_hub_name" {
  description = "The name of the hub"
  value       = azurerm_virtual_hub.main.name
}

output "default_route_table_id" {
  description = "The ARM ID of the hub's built-in default route table -- what a connection's routing associates with when no custom table is chosen"
  value       = azurerm_virtual_hub.main.default_route_table_id
}

output "virtual_router_asn" {
  description = "The hub router's BGP autonomous system number (always 65515) -- the remote ASN when peering NVAs with the hub"
  value       = azurerm_virtual_hub.main.virtual_router_asn
}

output "virtual_router_ips" {
  description = "The hub router's peering IPv4 addresses (a pair) -- the BGP neighbor addresses an NVA peers with"
  value       = azurerm_virtual_hub.main.virtual_router_ips
}

output "route_table_ids" {
  description = "The ARM ID of each custom route table on the hub, keyed by the table's name from the spec"
  value       = { for name, route_table in azurerm_virtual_hub_route_table.route_tables : name => route_table.id }
}

output "route_map_ids" {
  description = "The ARM ID of each route map on the hub, keyed by the map's name from the spec"
  value       = { for name, route_map in azurerm_route_map.route_maps : name => route_map.id }
}

output "bgp_connection_ids" {
  description = "The ARM ID of each BGP connection on the hub, keyed by the peering's name from the spec"
  value       = { for name, bgp_connection in azurerm_virtual_hub_bgp_connection.bgp_connections : name => bgp_connection.id }
}

output "routing_intent_id" {
  description = "The ARM ID of the hub's routing intent -- empty when no routing intent is configured"
  value       = try(values(azurerm_virtual_hub_routing_intent.routing_intent)[0].id, "")
}
