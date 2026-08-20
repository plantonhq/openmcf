output "proxy_name" {
  description = "The proxy's name - the provider's import ID and the AWS CLI/API join key"
  value       = aws_db_proxy.this.name
}

output "proxy_arn" {
  description = "The proxy's ARN"
  value       = aws_db_proxy.this.arn
}

output "endpoint" {
  description = "The proxy's DEFAULT endpoint DNS name - what applications connect to"
  value       = aws_db_proxy.this.endpoint
}

output "default_target_group_arn" {
  description = "The default target group's ARN"
  value       = aws_db_proxy_default_target_group.this.arn
}

output "default_target_group_name" {
  description = "The default target group's name (always \"default\" at AWS)"
  value       = aws_db_proxy_default_target_group.this.name
}

output "endpoint_addresses" {
  description = "Additional endpoints' DNS names keyed by endpoint name"
  value       = { for name, endpoint in aws_db_proxy_endpoint.this : name => endpoint.endpoint }
}

output "endpoint_arns" {
  description = "Additional endpoints' ARNs keyed by endpoint name"
  value       = { for name, endpoint in aws_db_proxy_endpoint.this : name => endpoint.arn }
}

output "target_type" {
  description = "The registered target's type as AWS reports it (RDS_INSTANCE or TRACKED_CLUSTER)"
  value       = var.spec.target != null ? aws_db_proxy_target.this[0].type : ""
}

output "target_rds_resource_id" {
  description = "The registered target's RDS resource id - part of the target's import ID"
  value       = var.spec.target != null ? aws_db_proxy_target.this[0].rds_resource_id : ""
}
