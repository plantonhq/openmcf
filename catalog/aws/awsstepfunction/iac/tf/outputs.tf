output "state_machine_arn" {
  description = "The ARN of the Step Functions state machine."
  value       = aws_sfn_state_machine.this.arn
}

output "state_machine_name" {
  description = "The name of the Step Functions state machine."
  value       = aws_sfn_state_machine.this.name
}

output "state_machine_version_arn" {
  description = "The ARN of the most recently published version (empty unless spec.publish is true)."
  value       = aws_sfn_state_machine.this.state_machine_version_arn
}

output "revision_id" {
  description = "The revision identifier of the current state machine definition."
  value       = aws_sfn_state_machine.this.revision_id
}

output "status" {
  description = "Lifecycle status of the state machine as reported by AWS (e.g. ACTIVE)."
  value       = aws_sfn_state_machine.this.status
}

output "creation_date" {
  description = "RFC3339 timestamp of when the state machine was created."
  value       = aws_sfn_state_machine.this.creation_date
}

output "alias_arns" {
  description = "Alias ARNs keyed by alias name (spec.aliases[].name)."
  value       = { for name, alias in aws_sfn_alias.this : name => alias.arn }
}
