output "domain_name" {
  description = "The custom domain name"
  value       = aws_api_gateway_domain_name.this.domain_name
}

output "domain_name_arn" {
  description = "The domain's ARN"
  value       = aws_api_gateway_domain_name.this.arn
}

output "domain_name_id" {
  description = "The domain name ID (distinguishes PRIVATE domains sharing one hostname)"
  value       = aws_api_gateway_domain_name.this.domain_name_id
}

output "regional_domain_name" {
  description = "The regional target hostname for Route 53 alias records (REGIONAL/PRIVATE domains)"
  value       = aws_api_gateway_domain_name.this.regional_domain_name
}

output "regional_zone_id" {
  description = "The Route 53 hosted zone ID of the regional endpoint"
  value       = aws_api_gateway_domain_name.this.regional_zone_id
}

output "cloudfront_domain_name" {
  description = "The CloudFront target hostname for alias records (EDGE domains)"
  value       = aws_api_gateway_domain_name.this.cloudfront_domain_name
}

output "cloudfront_zone_id" {
  description = "The CloudFront hosted zone ID"
  value       = aws_api_gateway_domain_name.this.cloudfront_zone_id
}

output "base_path_mapping_ids" {
  description = "Base-path mapping IDs keyed by base path ((root) for the empty path)"
  value       = { for k, m in aws_api_gateway_base_path_mapping.this : k => m.id }
}

output "access_association_arns" {
  description = "Access-association ARNs keyed by the granted VPC endpoint ID"
  value       = { for k, a in aws_api_gateway_domain_name_access_association.this : k => a.arn }
}
