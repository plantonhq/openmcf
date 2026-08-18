output "account_id" {
  description = "The member account's 12-digit AWS account ID (also the provider's import ID)"
  value       = aws_organizations_account.this.id
}

output "arn" {
  description = "The member account's ARN"
  value       = aws_organizations_account.this.arn
}

output "state" {
  description = "The account's lifecycle state (ACTIVE, SUSPENDED, PENDING_CLOSURE)"
  value       = aws_organizations_account.this.state
}

output "govcloud_id" {
  description = "The companion GovCloud (US) account's ID, when create_govcloud was set"
  value       = aws_organizations_account.this.govcloud_id
}
