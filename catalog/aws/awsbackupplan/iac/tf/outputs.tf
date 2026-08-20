output "plan_id" {
  description = "The plan's AWS-generated ID (a UUID - also the provider's import ID)"
  value       = aws_backup_plan.this.id
}

output "plan_arn" {
  description = "The plan's ARN"
  value       = aws_backup_plan.this.arn
}

output "plan_version" {
  description = "The plan's version ID (changes on every plan update)"
  value       = aws_backup_plan.this.version
}

output "selection_ids" {
  description = "AWS-generated selection IDs keyed by selection name (each imports as \"{plan_id}|{selection_id}\")"
  value       = { for name, sel in aws_backup_selection.this : name => sel.id }
}
