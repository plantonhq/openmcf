output "launch_template_id" {
  description = "The launch template ID (what auto-scaling groups, EKS node groups, and Batch compute environments reference)."
  value       = aws_launch_template.this.id
}

output "launch_template_arn" {
  description = "The ARN of the launch template, for IAM policies that scope ec2:RunInstances to approved templates."
  value       = aws_launch_template.this.arn
}

output "latest_version" {
  description = "The latest version number of the template (every spec change creates a new immutable version)."
  value       = aws_launch_template.this.latest_version
}

output "default_version" {
  description = "The default version number -- what consumers referencing \"$Default\" launch from."
  value       = aws_launch_template.this.default_version
}
