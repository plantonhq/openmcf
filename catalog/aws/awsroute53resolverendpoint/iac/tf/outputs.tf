output "endpoint_id" {
  description = "The endpoint's id (rslvr-in-.../rslvr-out-...) - the provider's import ID"
  value       = aws_route53_resolver_endpoint.this.id
}

output "endpoint_arn" {
  description = "The endpoint's ARN"
  value       = aws_route53_resolver_endpoint.this.arn
}

output "host_vpc_id" {
  description = "The VPC the endpoint's subnets belong to (AWS derives it)"
  value       = aws_route53_resolver_endpoint.this.host_vpc_id
}

output "ip_addresses" {
  description = "The ENI IP addresses the endpoint answers or originates on"
  value       = [for entry in aws_route53_resolver_endpoint.this.ip_address : entry.ip]
}

output "rule_ids" {
  description = "AWS-generated rule IDs (rslvr-rr-...) keyed by rule name"
  value       = { for name, rule in aws_route53_resolver_rule.this : name => rule.id }
}

output "rule_association_ids" {
  description = "AWS-generated rule association IDs (rslvr-rrassoc-...) keyed by rule_name//vpc_id"
  value       = { for key, assoc in aws_route53_resolver_rule_association.this : key => assoc.id }
}
