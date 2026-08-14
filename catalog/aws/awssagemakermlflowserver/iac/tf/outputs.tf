output "tracking_server_name" {
  description = "The tracking server name (the AWS identity)"
  value       = aws_sagemaker_mlflow_tracking_server.this.tracking_server_name
}

output "tracking_server_arn" {
  description = "The Amazon Resource Name of the tracking server"
  value       = aws_sagemaker_mlflow_tracking_server.this.arn
}

output "tracking_server_url" {
  description = "URL of the MLflow UI served by this tracking server"
  value       = aws_sagemaker_mlflow_tracking_server.this.tracking_server_url
}
