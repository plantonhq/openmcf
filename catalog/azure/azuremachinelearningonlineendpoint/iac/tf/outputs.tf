output "online_endpoint_id" {
  description = "The Azure Resource Manager ID of the online endpoint"
  value       = azapi_resource.main.id
}

output "online_endpoint_name" {
  description = "The endpoint's name -- what deployments attach to and what the traffic map's keys route across"
  value       = azapi_resource.main.name
}

output "scoring_uri" {
  description = "The HTTPS address applications POST scoring requests to"
  value       = try(azapi_resource.main.output.properties.scoringUri, "")
}

output "swagger_uri" {
  description = "The endpoint's OpenAPI (Swagger) document address"
  value       = try(azapi_resource.main.output.properties.swaggerUri, "")
}

output "system_assigned_identity_principal_id" {
  description = "The principal (object) ID of the endpoint's system-assigned identity, when one is enabled"
  value       = try(azapi_resource.main.identity[0].principal_id, "")
}
