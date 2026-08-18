output "layer_arn" {
  description = "The layer's unversioned ARN - the identity that persists across versions"
  value       = aws_lambda_layer_version.this.layer_arn
}

output "layer_version_arn" {
  description = "The published version's ARN - what functions attach"
  value       = aws_lambda_layer_version.this.arn
}

output "version" {
  description = "The published version number"
  value       = aws_lambda_layer_version.this.version
}

output "code_sha256" {
  description = "Base64-encoded SHA256 of the archive as Lambda stored it"
  value       = aws_lambda_layer_version.this.code_sha256
}

output "permission_revision_ids" {
  description = "Policy revision ids keyed by each grant's statement_id"
  value       = { for statement_id, permission in aws_lambda_layer_version_permission.this : statement_id => permission.revision_id }
}
