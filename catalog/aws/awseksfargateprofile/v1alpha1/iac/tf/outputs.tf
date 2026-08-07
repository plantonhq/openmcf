output "fargate_profile_arn" {
  description = "The Amazon Resource Name of the Fargate profile."
  value       = aws_eks_fargate_profile.this.arn
}

output "fargate_profile_name" {
  description = "The name of the Fargate profile."
  value       = aws_eks_fargate_profile.this.fargate_profile_name
}

output "status" {
  description = "The profile's state after provisioning -- ACTIVE on a successful create."
  value       = aws_eks_fargate_profile.this.status
}
