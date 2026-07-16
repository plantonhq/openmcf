output "domain_id" {
  description = "The unique identifier of the OpenSearch domain."
  value       = aws_opensearch_domain.this.domain_id
}

output "domain_name" {
  description = "The name of the OpenSearch domain."
  value       = aws_opensearch_domain.this.domain_name
}

output "domain_arn" {
  description = "The ARN of the OpenSearch domain."
  value       = aws_opensearch_domain.this.arn
}

output "endpoint" {
  description = "The domain-specific endpoint for index and search requests."
  value       = aws_opensearch_domain.this.endpoint
}

output "dashboard_endpoint" {
  description = "The OpenSearch Dashboards endpoint."
  value       = aws_opensearch_domain.this.dashboard_endpoint
}

output "endpoint_v2" {
  description = "The dual-stack (IPv4 + IPv6) V2 domain endpoint."
  value       = aws_opensearch_domain.this.endpoint_v2
}

output "dashboard_endpoint_v2" {
  description = "The OpenSearch Dashboards endpoint on the dual-stack V2 domain endpoint."
  value       = aws_opensearch_domain.this.dashboard_endpoint_v2
}

output "domain_endpoint_v2_hosted_zone_id" {
  description = "The Route 53 hosted zone ID for aliasing DNS records at the V2 domain endpoint."
  value       = aws_opensearch_domain.this.domain_endpoint_v2_hosted_zone_id
}
