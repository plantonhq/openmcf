output "gateway_id" {
  description = "The unique gateway identifier"
  value       = aws_bedrockagentcore_gateway.this.gateway_id
}

output "gateway_arn" {
  description = "The Amazon Resource Name of the gateway"
  value       = aws_bedrockagentcore_gateway.this.gateway_arn
}

output "gateway_url" {
  description = "The MCP URL agents connect to"
  value       = aws_bedrockagentcore_gateway.this.gateway_url
}

output "workload_identity_arn" {
  description = "ARN of the workload identity AWS created for this gateway"
  value       = one(aws_bedrockagentcore_gateway.this.workload_identity_details[*].workload_identity_arn)
}

output "target_ids" {
  description = "Target IDs keyed by each targets entry's name"
  value       = { for k, t in aws_bedrockagentcore_gateway_target.this : k => t.target_id }
}
