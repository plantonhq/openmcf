output "rule_arn" {
  description = "The ARN of the listener rule."
  value       = aws_lb_listener_rule.this.arn
}

output "priority" {
  description = "The priority AWS assigned to the rule (meaningful when the spec left it unset)."
  value       = tostring(aws_lb_listener_rule.this.priority)
}
