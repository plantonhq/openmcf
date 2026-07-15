output "service_arn" {
  description = "The ARN of the ECS service keeping the runner running."
  value       = aws_ecs_service.runner.arn
}

output "service_name" {
  description = "The service's name (metadata.name)."
  value       = aws_ecs_service.runner.name
}

output "cluster_arn" {
  description = "The ARN of the dedicated ECS cluster the runner runs in."
  value       = aws_ecs_cluster.runner.arn
}

output "task_definition_arn" {
  description = "The task definition ARN (family:revision) of the running revision."
  value       = aws_ecs_task_definition.runner.arn
}

output "log_group_name" {
  description = "The CloudWatch log group carrying the runner's operation audit trail."
  value       = aws_cloudwatch_log_group.runner.name
}

output "security_group_id" {
  description = "The outbound-only security group -- private targets reference it to trust the runner."
  value       = aws_security_group.runner.id
}

output "execution_role_arn" {
  description = "The setup identity: pulls the runner image, writes logs, reads the credentials secret."
  value       = aws_iam_role.execution.arn
}

output "task_role_arn" {
  description = "The runner's runtime identity -- grant it permissions for keyless cloud access."
  value       = local.task_role_arn
}

output "credentials_secret_arn" {
  description = "The Secrets Manager secret holding the runner's credentials document."
  value       = aws_secretsmanager_secret.credentials.arn
}

output "region" {
  description = "The AWS region the runner was deployed in."
  value       = var.spec.region
}
