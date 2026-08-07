output "function_arn" {
  description = "The function ARN -- the join key for event-source mappings, trigger configurations, and IAM policy resources."
  value       = aws_lambda_function.this.arn
}

output "function_name" {
  description = "The function name -- what SDK invoke calls and CLI commands reference."
  value       = aws_lambda_function.this.function_name
}

output "invoke_arn" {
  description = "The apigateway-shaped invocation ARN -- what API Gateway integrations reference."
  value       = aws_lambda_function.this.invoke_arn
}

output "qualified_arn" {
  description = "The qualified ARN of the most recently published version. Empty when publish is disabled."
  value       = var.spec.publish ? aws_lambda_function.this.qualified_arn : ""
}

output "version" {
  description = "The most recently published version number. Empty when publish is disabled."
  value       = var.spec.publish ? aws_lambda_function.this.version : ""
}

output "function_url" {
  description = "The HTTPS endpoint of the function URL. Empty when no function URL is configured."
  value       = var.spec.function_url != null ? aws_lambda_function_url.this[0].function_url : ""
}

output "alias_arns" {
  description = "ARNs of the function's aliases keyed by alias name -- the stable invocation targets for traffic-shifted rollouts."
  value       = { for name, alias in aws_lambda_alias.this : name => alias.arn }
}

output "log_group_name" {
  description = "The CloudWatch log group receiving the function's logs -- the AWS default /aws/lambda/<name> or the custom group from logging_config."
  value       = local.log_group_name
}
