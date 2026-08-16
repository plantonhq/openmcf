output "usage_plan_id" {
  description = "The usage plan ID"
  value       = aws_api_gateway_usage_plan.this.id
}

output "usage_plan_arn" {
  description = "The usage plan ARN"
  value       = aws_api_gateway_usage_plan.this.arn
}

output "api_key_ids" {
  description = "API key IDs keyed by each api_keys entry's name"
  value       = { for k, v in aws_api_gateway_api_key.this : k => v.id }
}

output "api_key_arns" {
  description = "API key ARNs keyed by each api_keys entry's name"
  value       = { for k, v in aws_api_gateway_api_key.this : k => v.arn }
}
