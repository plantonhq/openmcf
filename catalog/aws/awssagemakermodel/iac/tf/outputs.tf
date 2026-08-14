output "model_name" {
  description = "The model name (the AWS identity endpoint variants reference)"
  value       = aws_sagemaker_model.this.name
}

output "model_arn" {
  description = "The Amazon Resource Name of the model"
  value       = aws_sagemaker_model.this.arn
}
