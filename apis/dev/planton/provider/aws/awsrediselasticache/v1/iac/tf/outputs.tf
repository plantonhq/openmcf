output "replication_group_id" {
  description = "The identifier of the replication group."
  value       = aws_elasticache_replication_group.this.replication_group_id
}

output "primary_endpoint_address" {
  description = "The primary (writer) endpoint DNS name."
  value       = aws_elasticache_replication_group.this.primary_endpoint_address
}

output "reader_endpoint_address" {
  description = "The reader endpoint DNS name for read replicas."
  value       = aws_elasticache_replication_group.this.reader_endpoint_address
}

output "configuration_endpoint_address" {
  description = "The configuration endpoint for Cluster Mode Enabled."
  value       = aws_elasticache_replication_group.this.configuration_endpoint_address
}

output "arn" {
  description = "The ARN of the replication group."
  value       = aws_elasticache_replication_group.this.arn
}

output "port" {
  description = "The port on which the cluster accepts connections."
  value       = aws_elasticache_replication_group.this.port
}

output "engine_version_actual" {
  description = "The engine version actually running -- meaningful when the spec leaves engine_version to the AWS default."
  value       = aws_elasticache_replication_group.this.engine_version_actual
}

output "subnet_group_name" {
  description = "The name of the created subnet group (empty if none created)."
  value       = local.manage_subnet_group ? aws_elasticache_subnet_group.this[0].name : ""
}

output "parameter_group_name" {
  description = "The name of the created parameter group (empty if none created)."
  value       = local.manage_parameter_group ? aws_elasticache_parameter_group.this[0].name : ""
}
