output "namespace_name" {
  description = "The namespace name -- the join key workgroups attach with."
  value       = aws_redshiftserverless_namespace.this.namespace_name
}

output "namespace_id" {
  description = "The unique identifier AWS assigns to the namespace."
  value       = aws_redshiftserverless_namespace.this.namespace_id
}

output "arn" {
  description = "The namespace's Amazon Resource Name, for IAM policies, usage limits, and resource policies."
  value       = aws_redshiftserverless_namespace.this.arn
}

output "db_name" {
  description = "The name of the first database in the namespace."
  value       = try(coalesce(aws_redshiftserverless_namespace.this.db_name, ""), "")
}

output "admin_password_secret_arn" {
  description = "The ARN of the AWS-managed admin-password secret in Secrets Manager (only when manage_admin_password is true) -- the handle applications use to fetch credentials at runtime."
  value       = try(coalesce(aws_redshiftserverless_namespace.this.admin_password_secret_arn, ""), "")
}
