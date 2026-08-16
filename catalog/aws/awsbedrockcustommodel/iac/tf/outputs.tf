output "custom_model_arn" {
  description = "ARN of the resulting custom model -- the value provisioned throughput references to buy serving capacity."
  value       = aws_bedrock_custom_model.this.custom_model_arn
}

output "custom_model_name" {
  description = "The custom model's name. Matches metadata.name."
  value       = aws_bedrock_custom_model.this.custom_model_name
}

output "job_arn" {
  description = "ARN of the customization job that produced (or is producing) the model."
  value       = aws_bedrock_custom_model.this.job_arn
}

output "job_status" {
  description = "Customization job status at the end of the deploy (InProgress, Completed, Failed, Stopping, Stopped). Training continues asynchronously after the deploy returns."
  value       = aws_bedrock_custom_model.this.job_status
}
