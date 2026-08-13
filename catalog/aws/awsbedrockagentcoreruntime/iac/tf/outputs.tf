output "agent_runtime_id" {
  description = "The unique runtime identifier"
  value       = aws_bedrockagentcore_agent_runtime.this.agent_runtime_id
}

output "agent_runtime_arn" {
  description = "The Amazon Resource Name of the runtime"
  value       = aws_bedrockagentcore_agent_runtime.this.agent_runtime_arn
}

output "agent_runtime_version" {
  description = "The runtime's current version number"
  value       = aws_bedrockagentcore_agent_runtime.this.agent_runtime_version
}

output "workload_identity_arn" {
  description = "ARN of the workload identity AWS created for this runtime"
  value       = one(aws_bedrockagentcore_agent_runtime.this.workload_identity_details[*].workload_identity_arn)
}

output "endpoint_arns" {
  description = "Endpoint ARNs keyed by each endpoints entry's name"
  value       = { for k, e in aws_bedrockagentcore_agent_runtime_endpoint.this : k => e.agent_runtime_endpoint_arn }
}
