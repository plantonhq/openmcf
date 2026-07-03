output "cluster_identifier" {
  description = "The cluster identifier -- the handle DocumentDB APIs and the console use."
  value       = aws_docdb_cluster.this.cluster_identifier
}

output "arn" {
  description = "The cluster's Amazon Resource Name."
  value       = aws_docdb_cluster.this.arn
}

output "cluster_resource_id" {
  description = "The immutable cluster resource ID (cluster-...) -- survives identifier renames; the durable handle for point-in-time restores and CloudWatch dimensions."
  value       = aws_docdb_cluster.this.cluster_resource_id
}

output "endpoint" {
  description = "The writer endpoint -- connect here for reads and writes."
  value       = aws_docdb_cluster.this.endpoint
}

output "reader_endpoint" {
  description = "The reader endpoint -- load-balances connections across the cluster's reader instances."
  value       = aws_docdb_cluster.this.reader_endpoint
}

output "port" {
  description = "The port the cluster accepts connections on."
  value       = aws_docdb_cluster.this.port
}

output "hosted_zone_id" {
  description = "The Route53 hosted zone ID of the cluster endpoints, for DNS alias records."
  value       = aws_docdb_cluster.this.hosted_zone_id
}

output "engine_version_actual" {
  description = "The engine version actually running -- meaningful when the spec leaves engine_version to the AWS default."
  value       = aws_docdb_cluster.this.engine_version
}

output "master_user_secret_arn" {
  description = "The ARN of the AWS-managed master-user secret in Secrets Manager (only when manage_master_user_password is true) -- the handle applications use to fetch credentials at runtime."
  value       = try(aws_docdb_cluster.this.master_user_secret[0].secret_arn, "")
}

output "db_subnet_group_name" {
  description = "The DB subnet group the cluster runs in (managed here or referenced)."
  value       = try(coalesce(aws_docdb_cluster.this.db_subnet_group_name, ""), "")
}

output "db_cluster_parameter_group_name" {
  description = "The cluster parameter group in use (managed inline group, referenced group, or the engine default)."
  value       = try(coalesce(aws_docdb_cluster.this.db_cluster_parameter_group_name, ""), "")
}

output "instance_endpoints" {
  description = "Per-instance endpoints of the folded instances, in spec order. Empty for headless shapes (restores and global-cluster members created without instances)."
  value       = [for instance in var.spec.instances : aws_docdb_cluster_instance.this[instance.name].endpoint]
}
