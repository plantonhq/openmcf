output "provider_arn" {
  description = "The provider's ARN (also the provider's import ID, and the value role trust policies reference via the SAML principal)"
  value       = aws_iam_saml_provider.this.arn
}

output "saml_provider_uuid" {
  description = "The provider's AWS-assigned UUID"
  value       = aws_iam_saml_provider.this.saml_provider_uuid
}

output "valid_until" {
  description = "When the metadata document's certificates expire - rotate the document before this date or federation breaks"
  value       = aws_iam_saml_provider.this.valid_until
}
