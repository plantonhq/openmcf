output "repository_name" {
  description = "The ECR repository name (the slash-namespaced registry path)."
  value       = aws_ecr_repository.this.name
}

output "repository_url" {
  description = "The repository URL images are pushed to and pulled from (<registry_id>.dkr.ecr.<region>.amazonaws.com/<name>)."
  value       = aws_ecr_repository.this.repository_url
}

output "repository_arn" {
  description = "The ARN of the repository, for IAM policies scoping ECR actions."
  value       = aws_ecr_repository.this.arn
}

output "registry_id" {
  description = "The AWS account (registry) ID the repository belongs to."
  value       = aws_ecr_repository.this.registry_id
}
