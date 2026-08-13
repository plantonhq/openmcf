output "online_deployment_id" {
  description = "The Azure Resource Manager ID of the online deployment"
  value       = azapi_resource.main.id
}

output "online_deployment_name" {
  description = "The deployment's name -- the key the endpoint's traffic map routes scoring requests by"
  value       = azapi_resource.main.name
}
