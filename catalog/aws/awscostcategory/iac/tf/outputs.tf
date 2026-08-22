output "category_arn" {
  description = "The category's ARN (also the provider's import ID)"
  value       = aws_ce_cost_category.this.arn
}

output "category_name" {
  description = "The category's name (the key other expressions reference the category by)"
  value       = aws_ce_cost_category.this.name
}

output "effective_start" {
  description = "When the category's rules take effect (AWS-normalized month start)"
  value       = aws_ce_cost_category.this.effective_start
}

output "effective_end" {
  description = "When the category stops applying (set by AWS on deletion; normally empty)"
  value       = aws_ce_cost_category.this.effective_end
}
