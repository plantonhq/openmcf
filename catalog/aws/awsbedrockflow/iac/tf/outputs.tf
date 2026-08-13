output "flow_id" {
  description = "The unique flow identifier"
  value       = aws_bedrockagent_flow.this.id
}

output "flow_arn" {
  description = "The Amazon Resource Name of the flow"
  value       = aws_bedrockagent_flow.this.arn
}

output "draft_version" {
  description = "The flow's mutable working version (always DRAFT)"
  value       = "DRAFT"
}
