output "detector_id" {
  description = "The detector's AWS-assigned ID (the provider's import ID and every satellite's composition key)"
  value       = aws_guardduty_detector.this.id
}

output "detector_arn" {
  description = "The detector's ARN"
  value       = aws_guardduty_detector.this.arn
}

output "account_id" {
  description = "The AWS account ID the detector belongs to"
  value       = aws_guardduty_detector.this.account_id
}

output "ip_set_ids" {
  description = "Trusted IP list IDs keyed by each ip_sets entry's name"
  value       = { for k, v in aws_guardduty_ipset.this : k => v.ip_set_id }
}

output "threat_intel_set_ids" {
  description = "Threat intel list IDs keyed by each threat_intel_sets entry's name"
  value       = { for k, v in aws_guardduty_threatintelset.this : k => v.threat_intel_set_id }
}

output "publishing_destination_id" {
  description = "The findings-export destination ID (set only when spec.publishing_destination is configured)"
  value       = length(aws_guardduty_publishing_destination.this) > 0 ? aws_guardduty_publishing_destination.this[0].destination_id : ""
}
