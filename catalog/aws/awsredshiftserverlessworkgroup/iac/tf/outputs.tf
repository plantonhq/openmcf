output "workgroup_name" {
  description = "The workgroup name -- the handle the Redshift Serverless APIs and the credentials API address the workgroup by."
  value       = aws_redshiftserverless_workgroup.this.workgroup_name
}

output "workgroup_id" {
  description = "The unique identifier AWS assigns to the workgroup."
  value       = aws_redshiftserverless_workgroup.this.workgroup_id
}

output "arn" {
  description = "The workgroup's Amazon Resource Name, for IAM policies and usage limits."
  value       = aws_redshiftserverless_workgroup.this.arn
}

output "endpoint_address" {
  description = "The DNS hostname SQL clients connect to."
  value       = try(aws_redshiftserverless_workgroup.this.endpoint[0].address, "")
}

output "port" {
  description = "The port the workgroup accepts connections on."
  value       = aws_redshiftserverless_workgroup.this.port
}

output "endpoint_access_addresses" {
  description = "The private DNS addresses of the workgroup's VPC endpoints, keyed by endpoint name."
  value       = { for k, e in aws_redshiftserverless_endpoint_access.this : k => e.address }
}

output "usage_limit_ids" {
  description = "The AWS-generated usage-limit IDs, keyed by usage_type/period (unset period rendered as monthly) -- the handles delete-usage-limit and state import take."
  value       = { for k, l in aws_redshiftserverless_usage_limit.this : k => l.id }
}

output "custom_domain_certificate_expiry_time" {
  description = "When the custom domain's ACM certificate expires (RFC 3339) -- empty without a custom domain."
  value       = try(aws_redshiftserverless_custom_domain_association.this[0].custom_domain_certificate_expiry_time, "")
}
