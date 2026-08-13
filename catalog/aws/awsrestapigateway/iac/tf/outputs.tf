output "rest_api_id" {
  description = "The REST API ID"
  value       = aws_api_gateway_rest_api.this.id
}

output "rest_api_arn" {
  description = "The REST API ARN"
  value       = aws_api_gateway_rest_api.this.arn
}

output "execution_arn" {
  description = "The execution ARN prefix Lambda resource policies and IAM invoke statements scope to"
  value       = aws_api_gateway_rest_api.this.execution_arn
}

output "root_resource_id" {
  description = "The root (/) resource ID"
  value       = aws_api_gateway_rest_api.this.root_resource_id
}

output "stage_name" {
  description = "The stage name serving the API"
  value       = aws_api_gateway_stage.this.stage_name
}

output "stage_arn" {
  description = "The stage ARN (the WAF web-ACL association target)"
  value       = aws_api_gateway_stage.this.arn
}

output "invoke_url" {
  description = "The stage invoke URL"
  value       = aws_api_gateway_stage.this.invoke_url
}

output "deployment_id" {
  description = "The deployment ID the stage currently serves"
  value       = aws_api_gateway_deployment.this.id
}

output "client_certificate_id" {
  description = "The generated client certificate's ID (empty unless generated)"
  value       = length(aws_api_gateway_client_certificate.this) > 0 ? aws_api_gateway_client_certificate.this[0].id : ""
}

output "client_certificate_pem" {
  description = "The generated client certificate's PEM body (empty unless generated)"
  value       = length(aws_api_gateway_client_certificate.this) > 0 ? aws_api_gateway_client_certificate.this[0].pem_encoded_certificate : ""
}

output "resource_ids" {
  description = "API Gateway resource IDs keyed by route path (every derived tree node)"
  value       = local.resource_id_by_path
}

output "authorizer_ids" {
  description = "Authorizer IDs keyed by each authorizers entry's name"
  value       = { for k, a in aws_api_gateway_authorizer.this : k => a.id }
}

output "model_ids" {
  description = "Model IDs keyed by each models entry's name"
  value       = { for k, m in aws_api_gateway_model.this : k => m.id }
}

output "request_validator_ids" {
  description = "Request validator IDs keyed by each request_validators entry's name"
  value       = { for k, v in aws_api_gateway_request_validator.this : k => v.id }
}

output "documentation_part_ids" {
  description = "Documentation part IDs keyed by declaration position"
  value       = { for k, p in aws_api_gateway_documentation_part.this : k => p.id }
}
