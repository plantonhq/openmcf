output "access_entry_arn" {
  description = "The Amazon Resource Name of the access entry."
  value       = aws_eks_access_entry.this.access_entry_arn
}

output "principal_arn" {
  description = "The IAM principal the entry grants access to, as resolved at provisioning time."
  value       = aws_eks_access_entry.this.principal_arn
}
