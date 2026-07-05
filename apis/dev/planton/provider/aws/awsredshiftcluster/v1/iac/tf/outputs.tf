output "cluster_identifier" {
  description = "The cluster identifier -- the handle Redshift APIs and the console use."
  value       = aws_redshift_cluster.this.cluster_identifier
}

output "cluster_arn" {
  description = "The cluster's Amazon Resource Name, for IAM policies and cross-service references."
  value       = aws_redshift_cluster.this.arn
}

output "cluster_namespace_arn" {
  description = "The namespace ARN of the cluster, used by Redshift data sharing and the Redshift Data API."
  value       = aws_redshift_cluster.this.cluster_namespace_arn
}

output "endpoint" {
  description = "The connection endpoint in address:port form, for SQL client connection strings."
  value       = aws_redshift_cluster.this.endpoint
}

output "dns_name" {
  description = "The DNS hostname of the cluster's leader node (without port), for connection strings and DNS alias records."
  value       = aws_redshift_cluster.this.dns_name
}

output "database_name" {
  description = "The name of the first database in the cluster."
  value       = aws_redshift_cluster.this.database_name
}

output "port" {
  description = "The port the cluster accepts connections on."
  value       = aws_redshift_cluster.this.port
}

output "subnet_group_name" {
  description = "The Redshift subnet group the cluster runs in (managed here or referenced)."
  value       = try(coalesce(aws_redshift_cluster.this.cluster_subnet_group_name, ""), "")
}

output "parameter_group_name" {
  description = "The parameter group in use (managed inline group, referenced group, or the Redshift default)."
  value       = try(coalesce(aws_redshift_cluster.this.cluster_parameter_group_name, ""), "")
}

output "master_password_secret_arn" {
  description = "The ARN of the AWS-managed admin-password secret in Secrets Manager (only when manage_master_password is true) -- the handle applications use to fetch credentials at runtime."
  value       = try(coalesce(aws_redshift_cluster.this.master_password_secret_arn, ""), "")
}
