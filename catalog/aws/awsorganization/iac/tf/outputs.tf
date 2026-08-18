output "organization_id" {
  description = "The organization's AWS-generated ID (o-... - also the provider's import ID)"
  value       = aws_organizations_organization.this.id
}

output "arn" {
  description = "The organization's ARN"
  value       = aws_organizations_organization.this.arn
}

output "management_account_id" {
  description = "The management account's 12-digit account ID"
  value       = aws_organizations_organization.this.master_account_id
}

output "management_account_arn" {
  description = "The management account's ARN"
  value       = aws_organizations_organization.this.master_account_arn
}

output "management_account_email" {
  description = "The management account's email address"
  value       = aws_organizations_organization.this.master_account_email
}

output "root_id" {
  description = "The organization root's ID (r-... - the parent for first-level OUs and root-scoped policy attachments)"
  value       = aws_organizations_organization.this.roots[0].id
}

output "resource_policy_id" {
  description = "The folded resource policy's AWS-generated ID (rp-... - its import ID), empty when the arm is absent"
  value       = var.spec.resource_policy != null ? aws_organizations_resource_policy.this[0].id : ""
}
