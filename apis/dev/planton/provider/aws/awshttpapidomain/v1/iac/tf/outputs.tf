output "domain_name" {
  description = "The custom domain name (the domain's join key for downstream references)."
  value       = aws_apigatewayv2_domain_name.this.id
}

output "domain_name_arn" {
  description = "The ARN of the domain name resource."
  value       = aws_apigatewayv2_domain_name.this.arn
}

output "target_domain_name" {
  description = "The API Gateway-managed regional domain to target from DNS (alias/CNAME the custom domain to this value)."
  value       = one(aws_apigatewayv2_domain_name.this.domain_name_configuration[*].target_domain_name)
}

output "hosted_zone_id" {
  description = "The Route 53 hosted zone ID of the API Gateway regional endpoint (the alias target zone)."
  value       = one(aws_apigatewayv2_domain_name.this.domain_name_configuration[*].hosted_zone_id)
}
