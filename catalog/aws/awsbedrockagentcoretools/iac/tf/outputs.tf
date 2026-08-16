output "browser_ids" {
  description = "Browser IDs keyed by each browsers entry's name"
  value       = { for k, b in aws_bedrockagentcore_browser.this : k => b.browser_id }
}

output "browser_arns" {
  description = "Browser ARNs keyed by each browsers entry's name"
  value       = { for k, b in aws_bedrockagentcore_browser.this : k => b.browser_arn }
}

output "browser_profile_ids" {
  description = "Browser profile IDs keyed by each browser_profiles entry's name"
  value       = { for k, p in aws_bedrockagentcore_browser_profile.this : k => p.profile_id }
}

output "browser_profile_arns" {
  description = "Browser profile ARNs keyed by each browser_profiles entry's name"
  value       = { for k, p in aws_bedrockagentcore_browser_profile.this : k => p.profile_arn }
}

output "code_interpreter_ids" {
  description = "Code interpreter IDs keyed by each code_interpreters entry's name"
  value       = { for k, c in aws_bedrockagentcore_code_interpreter.this : k => c.code_interpreter_id }
}

output "code_interpreter_arns" {
  description = "Code interpreter ARNs keyed by each code_interpreters entry's name"
  value       = { for k, c in aws_bedrockagentcore_code_interpreter.this : k => c.code_interpreter_arn }
}
