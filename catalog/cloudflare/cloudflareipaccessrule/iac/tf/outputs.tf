output "rule_id" {
  description = "The ID of the created IP Access rule"
  value       = cloudflare_access_rule.main.id
}

output "zone_id" {
  description = "The zone the rule was created on (empty for account-wide rules)"
  value       = var.spec.zone_id
}

output "account_id" {
  description = "The account the rule was created on (empty for zone-scoped rules)"
  value       = var.spec.account_id
}
