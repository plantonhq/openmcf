output "secret_arn" {
  description = "The ARN of the secret -- the canonical join key for IAM policies and cross-service references (AWS appends a random 6-character suffix, so the ARN is not derivable from the name)."
  value       = aws_secretsmanager_secret.this.arn
}

output "secret_name" {
  description = "The name of the secret. Matches metadata.name."
  value       = aws_secretsmanager_secret.this.name
}

output "version_id" {
  description = "The version ID of the managed secret version (empty for a shell secret with no value)."
  value       = local.create_version ? aws_secretsmanager_secret_version.this[0].version_id : ""
}
