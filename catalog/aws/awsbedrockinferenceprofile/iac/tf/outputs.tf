output "inference_profile_arn" {
  description = "ARN of the inference profile -- the modelId applications pass to InvokeModel/Converse and the resource IAM policies scope to."
  value       = aws_bedrock_inference_profile.this.arn
}

output "inference_profile_id" {
  description = "The inference profile's unique identifier."
  value       = aws_bedrock_inference_profile.this.id
}

output "status" {
  description = "Profile status as reported by AWS (ACTIVE when usable)."
  value       = aws_bedrock_inference_profile.this.status
}

output "type" {
  description = "Profile type as reported by AWS -- APPLICATION for profiles this component creates."
  value       = aws_bedrock_inference_profile.this.type
}
