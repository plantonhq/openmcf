output "account_id" {
  description = "The 12-digit AWS account ID the settings belong to (also the provider's import ID for the suppression singleton)"
  value       = data.aws_caller_identity.this.account_id
}
