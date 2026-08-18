output "prefix_list_id" {
  description = "The prefix list's id (pl-...) - what security-group rules, NACL rules, and route tables reference, and the provider's import ID"
  value       = aws_ec2_managed_prefix_list.this.id
}

output "prefix_list_arn" {
  description = "The prefix list's ARN"
  value       = aws_ec2_managed_prefix_list.this.arn
}

output "owner_id" {
  description = "The AWS account that owns the list"
  value       = aws_ec2_managed_prefix_list.this.owner_id
}

output "version" {
  description = "The list's current version - AWS increments it on every entry change"
  value       = tostring(aws_ec2_managed_prefix_list.this.version)
}
