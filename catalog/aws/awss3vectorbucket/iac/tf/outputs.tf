output "vector_bucket_arn" {
  description = "The vector bucket's ARN - what policies and Bedrock knowledge bases reference, and the provider's import ID"
  value       = aws_s3vectors_vector_bucket.this.vector_bucket_arn
}

output "index_arns" {
  description = "Each index's ARN, keyed by index name - what a Bedrock knowledge base's s3_vectors arm points at"
  value       = { for name, index in aws_s3vectors_index.this : name => index.index_arn }
}
