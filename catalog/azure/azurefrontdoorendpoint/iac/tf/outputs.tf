output "endpoint_id" {
  description = "The Azure Resource Manager ID of the endpoint (what routes reference as their parent)"
  value       = azurerm_cdn_frontdoor_endpoint.main.id
}

output "endpoint_name" {
  description = "The endpoint's name -- unique within its profile"
  value       = azurerm_cdn_frontdoor_endpoint.main.name
}

output "host_name" {
  description = "The generated, globally unique *.azurefd.net hostname clients connect to -- the CNAME target for custom-domain DNS records"
  value       = azurerm_cdn_frontdoor_endpoint.main.host_name
}
