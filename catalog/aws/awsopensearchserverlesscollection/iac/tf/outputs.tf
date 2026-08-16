output "collection_id" {
  description = "The unique ID of the collection (the API's own identifier, also the leading label of the collection endpoints)."
  value       = aws_opensearchserverless_collection.this.id
}

output "collection_arn" {
  description = "The ARN of the collection. Used in IAM policies and as the vector-store ARN when a Bedrock knowledge base consumes the collection."
  value       = aws_opensearchserverless_collection.this.arn
}

output "collection_name" {
  description = "The name of the collection. Matches metadata.name."
  value       = aws_opensearchserverless_collection.this.name
}

output "collection_endpoint" {
  description = "Collection-specific endpoint for OpenSearch API operations (HTTPS, SigV4-authenticated)."
  value       = aws_opensearchserverless_collection.this.collection_endpoint
}

output "dashboard_endpoint" {
  description = "Collection-specific endpoint for OpenSearch Dashboards."
  value       = aws_opensearchserverless_collection.this.dashboard_endpoint
}

output "kms_key_arn" {
  description = "The ARN of the KMS key encrypting the collection -- the AWS-owned key or the customer-managed key chosen in spec.encryption."
  value       = aws_opensearchserverless_collection.this.kms_key_arn
}
