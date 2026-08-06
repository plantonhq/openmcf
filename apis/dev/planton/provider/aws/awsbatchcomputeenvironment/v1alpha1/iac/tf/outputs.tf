output "compute_environment_arn" {
  description = "The ARN of the compute environment -- what job queues reference in their compute_environment_order."
  value       = aws_batch_compute_environment.this.arn
}

output "compute_environment_name" {
  description = "The compute environment's name (derived from metadata.name)."
  value       = aws_batch_compute_environment.this.name
}

output "ecs_cluster_arn" {
  description = "The ARN of the ECS cluster AWS Batch provisions behind the MANAGED environment."
  value       = aws_batch_compute_environment.this.ecs_cluster_arn
}

output "status" {
  description = "The environment's current status (VALID / INVALID); queues can only associate VALID environments."
  value       = aws_batch_compute_environment.this.status
}
