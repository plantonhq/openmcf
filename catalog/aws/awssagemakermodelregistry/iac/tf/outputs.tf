output "model_package_group_name" {
  description = "The model package group name (the AWS identity training pipelines register packages into)"
  value       = aws_sagemaker_model_package_group.this.model_package_group_name
}

output "model_package_group_arn" {
  description = "The Amazon Resource Name of the model package group"
  value       = aws_sagemaker_model_package_group.this.arn
}
