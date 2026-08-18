output "budget_name" {
  description = "The budget's name (with the account ID, the provider's import ID: account_id:budget_name)"
  value       = aws_budgets_budget.this.name
}

output "budget_arn" {
  description = "The budget's ARN"
  value       = aws_budgets_budget.this.arn
}

output "account_id" {
  description = "The account the budget belongs to (the deploying account unless spec.account_id targeted a member account)"
  value       = aws_budgets_budget.this.account_id
}

output "action_ids" {
  description = "AWS-generated action IDs keyed by action name (each action imports as account_id:action_id:budget_name)"
  value       = { for name, action in aws_budgets_budget_action.this : name => action.action_id }
}
