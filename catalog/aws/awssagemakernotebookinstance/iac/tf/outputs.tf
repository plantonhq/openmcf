output "notebook_instance_name" {
  description = "The notebook instance name (the AWS identity)"
  value       = aws_sagemaker_notebook_instance.this.name
}

output "notebook_instance_arn" {
  description = "The Amazon Resource Name of the notebook instance"
  value       = aws_sagemaker_notebook_instance.this.arn
}

output "url" {
  description = "URL to open the Jupyter notebook"
  value       = aws_sagemaker_notebook_instance.this.url
}

output "network_interface_id" {
  description = "The ENI SageMaker created in your subnet (VPC notebooks only)"
  value       = aws_sagemaker_notebook_instance.this.network_interface_id
}

output "lifecycle_config_name" {
  description = "The folded lifecycle configuration's name (empty when not configured)"
  value       = local.has_lifecycle ? aws_sagemaker_notebook_instance_lifecycle_configuration.this[0].name : ""
}
