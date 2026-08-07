output "user_id" {
  description = "The user's AWS identifier -- what user groups reference in their membership list and what the AWS CLI/API address."
  value       = aws_elasticache_user.this.user_id
}

output "arn" {
  description = "The user's Amazon Resource Name -- an IAM-authenticated client needs elasticache:Connect on both the user ARN and the cache ARN."
  value       = aws_elasticache_user.this.arn
}

output "user_name" {
  description = "The name this user presents in the AUTH command -- exported so application configuration can be wired from the resource graph."
  value       = aws_elasticache_user.this.user_name
}
