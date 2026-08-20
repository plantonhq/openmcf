output "provisioned_model_arn" {
  description = "ARN of the provisioned model -- the modelId applications pass to InvokeModel/Converse to consume the dedicated capacity."
  value       = aws_bedrock_provisioned_model_throughput.this.provisioned_model_arn
}

output "provisioned_model_name" {
  description = "The provisioned model's name. Matches metadata.name."
  value       = aws_bedrock_provisioned_model_throughput.this.provisioned_model_name
}
