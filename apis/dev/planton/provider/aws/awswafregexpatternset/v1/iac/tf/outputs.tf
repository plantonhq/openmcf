output "regex_pattern_set_arn" {
  description = "ARN of the regex pattern set -- the identifier web ACL regex_pattern_set_reference statements point at."
  value       = aws_wafv2_regex_pattern_set.this.arn
}

output "regex_pattern_set_id" {
  description = "AWS-assigned pattern set ID (UUID), used with name and scope in direct WAFv2 API calls."
  value       = aws_wafv2_regex_pattern_set.this.id
}

output "regex_pattern_set_name" {
  description = "The pattern set name as created in AWS (derived from metadata.name)."
  value       = aws_wafv2_regex_pattern_set.this.name
}
