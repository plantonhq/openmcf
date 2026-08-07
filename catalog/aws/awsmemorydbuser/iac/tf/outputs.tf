output "user_name" {
  description = "The user's name -- the AUTH identity ACLs reference in their membership list."
  value       = aws_memorydb_user.this.user_name
}

output "user_arn" {
  description = "The user's Amazon Resource Name -- an IAM-authenticated client needs memorydb:Connect on both the user ARN and the cluster ARN."
  value       = aws_memorydb_user.this.arn
}

output "minimum_engine_version" {
  description = "The minimum engine version the user's configuration requires -- an ACL (and the cluster it attaches to) must run at least this version."
  value       = aws_memorydb_user.this.minimum_engine_version
}
