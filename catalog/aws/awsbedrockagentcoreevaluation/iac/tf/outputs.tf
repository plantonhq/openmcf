output "evaluator_ids" {
  description = "Evaluator IDs keyed by each evaluators entry's name"
  value       = { for k, e in aws_bedrockagentcore_evaluator.this : k => e.evaluator_id }
}

output "evaluator_arns" {
  description = "Evaluator ARNs keyed by each evaluators entry's name"
  value       = { for k, e in aws_bedrockagentcore_evaluator.this : k => e.evaluator_arn }
}

output "harness_ids" {
  description = "Harness IDs keyed by each harnesses entry's name"
  value       = { for k, h in aws_bedrockagentcore_harness.this : k => h.harness_id }
}

output "harness_arns" {
  description = "Harness ARNs keyed by each harnesses entry's name"
  value       = { for k, h in aws_bedrockagentcore_harness.this : k => h.arn }
}

output "online_evaluation_config_ids" {
  description = "Online evaluation config IDs keyed by each online_evaluation_configs entry's name"
  value       = { for k, c in aws_bedrockagentcore_online_evaluation_config.this : k => c.online_evaluation_config_id }
}

output "online_evaluation_config_arns" {
  description = "Online evaluation config ARNs keyed by each online_evaluation_configs entry's name"
  value       = { for k, c in aws_bedrockagentcore_online_evaluation_config.this : k => c.online_evaluation_config_arn }
}

output "online_evaluation_output_log_groups" {
  description = "CloudWatch log group each online config writes its evaluation results to (server-assigned), keyed by the entry's name"
  value       = { for k, c in aws_bedrockagentcore_online_evaluation_config.this : k => try(c.output_config[0].cloudwatch_config[0].log_group_name, "") }
}
