output "user_group_id" {
  description = "The group's AWS identifier -- what caches reference to attach RBAC and what the AWS CLI/API address."
  value       = aws_elasticache_user_group.this.user_group_id
}

output "arn" {
  description = "The group's Amazon Resource Name -- used in IAM policies and cross-service permissions."
  value       = aws_elasticache_user_group.this.arn
}
