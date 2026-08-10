# The full resource name (projects/{p}/locations/{region}/workflows/{name})
# with the ambient project/region resolved — exactly what Eventarc
# destinations consume (the resource ID carries this form).
output "workflow_id" {
  description = "Full workflow resource name"
  value       = google_workflows_workflow.this.id
}

output "workflow_name" {
  description = "The short workflow name"
  value       = google_workflows_workflow.this.name
}

# A new revision is minted on every source / env-var / service-account
# change — compare across applies to confirm a deploy actually rolled.
output "revision_id" {
  description = "The deployed workflow revision"
  value       = google_workflows_workflow.this.revision_id
}

output "state" {
  description = "Workflow state (ACTIVE when executable)"
  value       = google_workflows_workflow.this.state
}
