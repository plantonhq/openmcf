output "workflow_name" {
  description = "The workflow's name -- its identity within the account, and what Worker workflow bindings reference"
  value       = cloudflare_workflow.main.workflow_name
}

output "version_id" {
  description = "The ID of the workflow version the registration produced"
  value       = cloudflare_workflow.main.version_id
}
