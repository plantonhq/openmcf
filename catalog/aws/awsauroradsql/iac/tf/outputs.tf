output "identifier" {
  description = "The AWS-generated cluster identifier - the provider's import ID"
  value       = aws_dsql_cluster.this.identifier
}

output "cluster_arn" {
  description = "The cluster's ARN - what a peer cluster's multi_region configuration references"
  value       = aws_dsql_cluster.this.arn
}

output "endpoint" {
  description = "The PostgreSQL connection host (AWS exposes no endpoint attribute; both modules derive the documented DNS shape)"
  value       = "${aws_dsql_cluster.this.identifier}.dsql.${var.spec.region}.on.aws"
}

output "vpc_endpoint_service_name" {
  description = "The VPC endpoint service name for PrivateLink connectivity"
  value       = aws_dsql_cluster.this.vpc_endpoint_service_name
}

output "encryption_type" {
  description = "How the cluster is encrypted as AWS reports it (AWS_OWNED_KMS_KEY or CUSTOMER_MANAGED_KMS_KEY)"
  value       = try(aws_dsql_cluster.this.encryption_details[0].encryption_type, "")
}
