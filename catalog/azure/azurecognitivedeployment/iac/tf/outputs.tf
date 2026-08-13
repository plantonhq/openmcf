output "deployment_id" {
  description = "The Azure Resource Manager ID of the deployment"
  value       = azurerm_cognitive_deployment.main.id
}

output "deployment_name" {
  description = "The deployment's name -- what applications pass as the model/deployment parameter when calling the account's endpoint"
  value       = azurerm_cognitive_deployment.main.name
}

output "model_version" {
  description = "The deployed model's version as ARM reports it -- the resolved value when the spec left version unset"
  value       = azurerm_cognitive_deployment.main.model[0].version
}
