# Stack outputs flatten onto AwsEcsClusterStackOutputs field-for-field;
# both engines export the same names so composition never depends on the
# engine.

output "cluster_name" {
  description = "The cluster name (mirrors metadata.name)"
  value       = aws_ecs_cluster.this.name
}

output "cluster_arn" {
  description = "The ARN of the cluster -- the join key AwsEcsService references"
  value       = aws_ecs_cluster.this.arn
}

output "capacity_provider_names" {
  description = "Every associated capacity provider: the Fargate built-ins plus the folded EC2 provider names, in spec order"
  value       = local.associated_capacity_providers
}

output "capacity_provider_arns" {
  description = "The ARNs of the EC2 capacity providers this cluster defines (empty for Fargate-only clusters), in spec order"
  value       = [for p in var.spec.ec2_capacity_providers : aws_ecs_capacity_provider.this[p.name].arn]
}
