output "ip_set_arn" {
  description = "ARN of the IP set -- the identifier web ACL ip_set_reference statements point at."
  value       = aws_wafv2_ip_set.this.arn
}

output "ip_set_id" {
  description = "AWS-assigned IP set ID (UUID), used with name and scope in direct WAFv2 API calls."
  value       = aws_wafv2_ip_set.this.id
}

output "ip_set_name" {
  description = "The IP set name as created in AWS (derived from metadata.name)."
  value       = aws_wafv2_ip_set.this.name
}
