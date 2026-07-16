output "service_arn" {
  description = "The ARN of the ECS service -- encodes both the cluster and service names."
  value       = aws_ecs_service.this.arn
}

output "service_name" {
  description = "The service's name, the ECS API's join key together with the cluster."
  value       = aws_ecs_service.this.name
}

output "cluster_arn" {
  description = "The ARN of the cluster the service runs in, republished for downstream joins."
  value       = var.spec.cluster_arn
}

output "task_definition_arn" {
  description = "The task definition ARN (family:revision) this deployment of the service is running."
  value       = var.spec.task_definition
}
