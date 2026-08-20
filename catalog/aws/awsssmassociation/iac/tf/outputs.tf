output "association_id" {
  description = "The association's AWS-generated ID (a UUID - also the provider's import ID)"
  value       = aws_ssm_association.this.association_id
}

output "association_arn" {
  description = "The association's ARN"
  value       = aws_ssm_association.this.arn
}

output "document_name" {
  description = "The document name the association resolved to"
  value       = aws_ssm_association.this.name
}
