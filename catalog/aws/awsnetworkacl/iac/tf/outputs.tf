output "network_acl_id" {
  description = "The ACL's id (acl-...) - the provider's import ID"
  value       = aws_network_acl.this.id
}

output "network_acl_arn" {
  description = "The ACL's ARN"
  value       = aws_network_acl.this.arn
}

output "owner_id" {
  description = "The AWS account that owns the ACL"
  value       = aws_network_acl.this.owner_id
}
