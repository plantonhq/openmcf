output "batch_deployment_id" {
  description = "The Azure Resource Manager ID of the batch deployment"
  value       = azapi_resource.main.id
}

output "batch_deployment_name" {
  description = "The deployment's name -- the key the endpoint's default-deployment pointer routes job submissions by"
  value       = azapi_resource.main.name
}
