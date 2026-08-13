output "app_arn" {
  description = "The Amazon Resource Name of the MLflow app (the AWS identity)"
  value       = aws_sagemaker_mlflow_app.this.arn
}

output "app_name" {
  description = "The app name"
  value       = aws_sagemaker_mlflow_app.this.name
}
