output "cluster_identifier" {
  description = "The cluster identifier -- the handle Neptune APIs and the console use."
  value       = aws_neptune_cluster.this.cluster_identifier
}

output "arn" {
  description = "The cluster's Amazon Resource Name."
  value       = aws_neptune_cluster.this.arn
}

output "cluster_resource_id" {
  description = "The immutable cluster resource ID (cluster-...) -- survives identifier renames; the durable handle for CloudWatch dimensions and IAM database authentication policies."
  value       = aws_neptune_cluster.this.cluster_resource_id
}

output "endpoint" {
  description = "The writer endpoint -- send Gremlin/openCypher/SPARQL queries here for reads and writes."
  value       = aws_neptune_cluster.this.endpoint
}

output "reader_endpoint" {
  description = "The reader endpoint -- load-balances read-only queries across the cluster's reader instances."
  value       = aws_neptune_cluster.this.reader_endpoint
}

output "port" {
  description = "The port the cluster accepts connections on."
  value       = aws_neptune_cluster.this.port
}

output "hosted_zone_id" {
  description = "The Route53 hosted zone ID of the cluster endpoints, for DNS alias records."
  value       = aws_neptune_cluster.this.hosted_zone_id
}

output "engine_version_actual" {
  description = "The engine version actually running -- meaningful when the spec leaves engine_version to the AWS default."
  value       = aws_neptune_cluster.this.engine_version
}

output "neptune_subnet_group_name" {
  description = "The Neptune subnet group the cluster runs in (managed here or referenced)."
  value       = try(coalesce(aws_neptune_cluster.this.neptune_subnet_group_name, ""), "")
}

output "neptune_cluster_parameter_group_name" {
  description = "The cluster parameter group in use (managed inline group, referenced group, or the engine default)."
  value       = try(coalesce(aws_neptune_cluster.this.neptune_cluster_parameter_group_name, ""), "")
}

output "instance_endpoints" {
  description = "Per-instance endpoints of the folded instances, in spec order. Empty for headless shapes (restores, replicas, and global-cluster members created without instances)."
  value       = [for instance in var.spec.instances : aws_neptune_cluster_instance.this[instance.name].endpoint]
}

output "custom_endpoint_addresses" {
  description = "Custom endpoint DNS addresses keyed by spec.custom_endpoints[].name."
  value       = { for name, endpoint in aws_neptune_cluster_endpoint.this : name => endpoint.endpoint }
}

output "neptune_instance_parameter_group_name" {
  description = "The module-managed instance parameter group created from spec.instance_parameters (empty when not managed here)."
  value       = local.manage_instance_parameter_group ? aws_neptune_parameter_group.instance[0].name : ""
}
