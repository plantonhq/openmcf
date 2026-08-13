output "dns_resolver_id" {
  description = "The Azure Resource Manager ID of the DNS Private Resolver"
  value       = azurerm_private_dns_resolver.main.id
}

output "dns_resolver_name" {
  description = "The name of the DNS Private Resolver resource"
  value       = azurerm_private_dns_resolver.main.name
}

output "inbound_endpoint_ip" {
  description = "The private IP of the FIRST inbound endpoint declared in the spec -- what on-premises forwarders point at; empty when no inbound endpoints are declared"
  value = (
    length(var.spec.inbound_endpoints) > 0
    ? azurerm_private_dns_resolver_inbound_endpoint.main[var.spec.inbound_endpoints[0].name].ip_configurations[0].private_ip_address
    : ""
  )
}

output "inbound_endpoint_ips" {
  description = "The private IPs of ALL inbound endpoints, keyed by endpoint name"
  value = {
    for name, endpoint in azurerm_private_dns_resolver_inbound_endpoint.main :
    name => endpoint.ip_configurations[0].private_ip_address
  }
}

output "outbound_endpoint_id" {
  description = "The ARM id of the FIRST outbound endpoint declared in the spec -- what forwarding rulesets reference by default; empty when no outbound endpoints are declared"
  value = (
    length(var.spec.outbound_endpoints) > 0
    ? azurerm_private_dns_resolver_outbound_endpoint.main[var.spec.outbound_endpoints[0].name].id
    : ""
  )
}

output "outbound_endpoint_ids" {
  description = "The ARM ids of ALL outbound endpoints, keyed by endpoint name"
  value = {
    for name, endpoint in azurerm_private_dns_resolver_outbound_endpoint.main :
    name => endpoint.id
  }
}
