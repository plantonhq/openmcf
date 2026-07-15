output "cluster_endpoint_address" {
  description = "The DNS address of the cluster endpoint -- applications connect here; MemoryDB handles slot discovery and routing behind it."
  value       = try(aws_memorydb_cluster.this.cluster_endpoint[0].address, "")
}

output "cluster_endpoint_port" {
  description = "The port of the cluster endpoint."
  value       = try(aws_memorydb_cluster.this.cluster_endpoint[0].port, 0)
}

output "cluster_arn" {
  description = "The ARN of the MemoryDB cluster -- used in IAM policies (memorydb:Connect) and cross-service permissions."
  value       = aws_memorydb_cluster.this.arn
}

output "cluster_name" {
  description = "The name of the MemoryDB cluster (metadata.name)."
  value       = aws_memorydb_cluster.this.name
}

output "engine_patch_version" {
  description = "The actual engine patch version running on the cluster -- may differ from the requested engine_version due to automatic patching."
  value       = aws_memorydb_cluster.this.engine_patch_version
}

output "subnet_group_name" {
  description = "The subnet group the cluster is placed in -- module-managed, bring-your-own, or empty when AWS's account default was used."
  value       = local.effective_subnet_group_name != null ? local.effective_subnet_group_name : ""
}

output "parameter_group_name" {
  description = "The parameter group attached to the cluster -- module-managed, bring-your-own, or empty when the family default applies."
  value       = local.effective_parameter_group_name != null ? local.effective_parameter_group_name : ""
}
