output "rule_group_id" {
  description = "The rule group's id (rslvr-frg-...) - the provider's import ID and half of each rule's composite import ID"
  value       = aws_route53_resolver_firewall_rule_group.this.id
}

output "rule_group_arn" {
  description = "The rule group's ARN"
  value       = aws_route53_resolver_firewall_rule_group.this.arn
}

output "share_status" {
  description = "Whether the group is shared via RAM (NOT_SHARED / SHARED_BY_ME / SHARED_WITH_ME)"
  value       = aws_route53_resolver_firewall_rule_group.this.share_status
}

output "domain_list_ids" {
  description = "AWS-generated domain list IDs (rslvr-fdl-...) keyed by list name"
  value       = { for name, list in aws_route53_resolver_firewall_domain_list.this : name => list.id }
}

output "association_ids" {
  description = "AWS-generated VPC association IDs (rslvr-frgassoc-...) keyed by association name"
  value       = { for name, assoc in aws_route53_resolver_firewall_rule_group_association.this : name => assoc.id }
}

output "rule_match_ids" {
  description = "Each rule's match identity keyed by rule name - the domain list ID for standard rules, the threat-protection ID for advanced rules (the second half of the rule's composite import ID)"
  value       = { for name, rule in aws_route53_resolver_firewall_rule.this : name => (rule.firewall_domain_list_id != null ? rule.firewall_domain_list_id : rule.firewall_threat_protection_id) }
}
