output "project_arn" {
  description = "ARN of the CodeBuild project -- the handle IAM policies and EventBridge targets reference"
  value       = aws_codebuild_project.this.arn
}

output "project_name" {
  description = "Name of the CodeBuild project -- what CodePipeline Build actions reference via the ProjectName configuration key"
  value       = aws_codebuild_project.this.name
}

output "service_role_arn" {
  description = "IAM service role ARN used by the project"
  value       = aws_codebuild_project.this.service_role
}

output "badge_url" {
  description = "Build badge URL (empty when badge_enabled is false)"
  value       = aws_codebuild_project.this.badge_url != null ? aws_codebuild_project.this.badge_url : ""
}

output "public_project_alias" {
  description = "Public alias of the project (empty unless project_visibility is PUBLIC_READ)"
  value       = aws_codebuild_project.this.public_project_alias != null ? aws_codebuild_project.this.public_project_alias : ""
}

output "webhook_url" {
  description = "Webhook URL at the source provider (empty if no webhook)"
  value       = local.has_webhook ? aws_codebuild_webhook.this[0].url : ""
}

output "webhook_payload_url" {
  description = "Webhook payload URL -- with manual_creation, register this on the repository by hand (empty if no webhook)"
  value       = local.has_webhook ? aws_codebuild_webhook.this[0].payload_url : ""
}

output "webhook_secret" {
  description = "Webhook HMAC signing secret, only minted at webhook creation. Sensitive -- treat as a credential (empty if no webhook)"
  value       = local.has_webhook ? aws_codebuild_webhook.this[0].secret : ""
  sensitive   = true
}
