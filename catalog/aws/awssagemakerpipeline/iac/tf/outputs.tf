output "pipeline_name" {
  description = "The pipeline name (the AWS identity executions start against)"
  value       = aws_sagemaker_pipeline.this.pipeline_name
}

output "pipeline_arn" {
  description = "The Amazon Resource Name of the pipeline"
  value       = aws_sagemaker_pipeline.this.arn
}
