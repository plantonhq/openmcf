output "instance_identifier" {
  description = "The instance identifier -- the handle RDS APIs and the console use."
  value       = aws_db_instance.this.identifier
}

output "arn" {
  description = "The instance's Amazon Resource Name."
  value       = aws_db_instance.this.arn
}

output "resource_id" {
  description = "The immutable DB instance resource ID (db-...) -- survives identifier renames; the durable handle for point-in-time restores, IAM auth policies, and CloudWatch dimensions."
  value       = aws_db_instance.this.resource_id
}

output "endpoint" {
  description = "The connection endpoint in address:port form."
  value       = aws_db_instance.this.endpoint
}

output "address" {
  description = "The DNS address of the instance (endpoint without the port)."
  value       = aws_db_instance.this.address
}

output "port" {
  description = "The port the instance accepts connections on."
  value       = aws_db_instance.this.port
}

output "hosted_zone_id" {
  description = "The Route53 hosted zone ID of the endpoint, for DNS alias records."
  value       = aws_db_instance.this.hosted_zone_id
}

output "engine_version_actual" {
  description = "The engine version actually running -- meaningful when the spec leaves engine_version to the AWS default."
  value       = aws_db_instance.this.engine_version_actual
}

output "master_user_secret_arn" {
  description = "The ARN of the AWS-managed master-user secret in Secrets Manager (only when manage_master_user_password is true) -- the handle applications use to fetch credentials at runtime."
  value       = try(aws_db_instance.this.master_user_secret[0].secret_arn, "")
}

output "db_subnet_group_name" {
  description = "The DB subnet group the instance runs in (managed here or referenced)."
  value       = try(coalesce(aws_db_instance.this.db_subnet_group_name, ""), "")
}

output "db_parameter_group_name" {
  description = "The DB parameter group in use (managed inline group, referenced group, or empty for the engine default)."
  value       = try(coalesce(aws_db_instance.this.parameter_group_name, ""), "")
}

output "option_group_name" {
  description = "The option group in use (managed inline group, referenced group, or empty for the engine default)."
  value       = try(coalesce(aws_db_instance.this.option_group_name, ""), "")
}
