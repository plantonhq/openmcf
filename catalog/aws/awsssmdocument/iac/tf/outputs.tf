output "document_name" {
  description = "The document's name (also the provider's import ID, and what associations reference)"
  value       = aws_ssm_document.this.name
}

output "document_arn" {
  description = "The document's ARN"
  value       = aws_ssm_document.this.arn
}

output "default_version" {
  description = "The default document version (updates promote the new version here)"
  value       = aws_ssm_document.this.default_version
}

output "latest_version" {
  description = "The latest document version"
  value       = aws_ssm_document.this.latest_version
}

output "document_hash" {
  description = "The Sha256 digest of the document content"
  value       = aws_ssm_document.this.hash
}

output "status" {
  description = "The document's lifecycle status (\"Creating\", \"Active\", ...)"
  value       = aws_ssm_document.this.status
}
