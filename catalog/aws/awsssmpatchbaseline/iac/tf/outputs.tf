output "baseline_id" {
  description = "The baseline's AWS-generated ID (\"pb-...\" - also the provider's import ID)"
  value       = aws_ssm_patch_baseline.this.id
}

output "baseline_arn" {
  description = "The baseline's ARN"
  value       = aws_ssm_patch_baseline.this.arn
}

output "operating_system" {
  description = "The operating system the baseline governs (WINDOWS when the spec leaves it unset; also the default designation's import ID)"
  value       = aws_ssm_patch_baseline.this.operating_system
}
