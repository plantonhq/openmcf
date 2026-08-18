output "account_id" {
  description = "The account these settings belong to (the singleton's identity; also the STS preference's resource ID at the provider)"
  value       = data.aws_caller_identity.this.account_id
}

output "account_alias" {
  description = "The applied sign-in alias (empty when the arm is unset)"
  value       = local.manage_alias ? aws_iam_account_alias.this[0].account_alias : ""
}

output "expire_passwords" {
  description = "Whether the applied password policy expires passwords (AWS derives it from max_password_age; empty when the arm is unset)"
  value       = local.manage_password_policy ? tostring(aws_iam_account_password_policy.this[0].expire_passwords) : ""
}
