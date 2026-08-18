output "ou_id" {
  description = "The OU's AWS-generated ID (ou-... - also the provider's import ID)"
  value       = aws_organizations_organizational_unit.this.id
}

output "arn" {
  description = "The OU's ARN"
  value       = aws_organizations_organizational_unit.this.arn
}
