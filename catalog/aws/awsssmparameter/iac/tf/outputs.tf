output "parameter_name" {
  description = "The parameter's name (also the provider's import ID; ARNs import too)"
  value       = aws_ssm_parameter.this.name
}

output "parameter_arn" {
  description = "The parameter's ARN"
  value       = aws_ssm_parameter.this.arn
}

output "version" {
  description = "The parameter's version number (increments on every value write)"
  value       = tostring(aws_ssm_parameter.this.version)
}

output "tier" {
  description = "The tier AWS resolved for the parameter (Standard or Advanced)"
  value       = aws_ssm_parameter.this.tier
}
