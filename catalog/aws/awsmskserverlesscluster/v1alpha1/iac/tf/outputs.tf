output "cluster_arn" {
  description = "ARN of the MSK Serverless cluster (also the resource identifier) -- referenced in IAM policies and Lambda event source mappings."
  value       = aws_msk_serverless_cluster.this.arn
}

output "cluster_name" {
  description = "Human-readable name of the cluster."
  value       = aws_msk_serverless_cluster.this.cluster_name
}

output "cluster_uuid" {
  description = "Unique identifier extracted from the cluster ARN."
  value       = aws_msk_serverless_cluster.this.cluster_uuid
}

output "bootstrap_brokers_sasl_iam" {
  description = "Comma-separated SASL/IAM broker endpoint list (port 9098) -- the only connection string serverless MSK exposes."
  value       = aws_msk_serverless_cluster.this.bootstrap_brokers_sasl_iam
}
