output "group_arn" {
  description = "The group's ARN"
  value       = aws_iam_group.this.arn
}

output "group_name" {
  description = "The group's name (also the provider's import ID, and the value policies and budget actions reference the group by)"
  value       = aws_iam_group.this.name
}

output "group_id" {
  description = "The group's AWS-generated stable unique ID (survives renames)"
  value       = aws_iam_group.this.unique_id
}
