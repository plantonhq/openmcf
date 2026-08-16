output "configured_region" {
  description = "The region whose invocation logging this instance owns (the singleton's identity and the provider's import ID)"
  value       = aws_bedrock_model_invocation_logging_configuration.this.id
}
