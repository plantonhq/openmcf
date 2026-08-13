output "endpoint_name" {
  description = "The endpoint name (the AWS identity clients invoke)"
  value       = aws_sagemaker_endpoint.this.name
}

output "endpoint_arn" {
  description = "The Amazon Resource Name of the endpoint"
  value       = aws_sagemaker_endpoint.this.arn
}

output "endpoint_config_name" {
  description = "The name of the endpoint configuration currently in service"
  value       = aws_sagemaker_endpoint_configuration.this.name
}

output "endpoint_config_arn" {
  description = "The Amazon Resource Name of that endpoint configuration"
  value       = aws_sagemaker_endpoint_configuration.this.arn
}
