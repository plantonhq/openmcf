# Fully qualified policy resource name
# (projects/{project}/locations/{location}/autoscalingPolicies/{policy_id}) —
# the handle a Dataproc cluster's autoscaling_policy_uri reference
# resolves to.
output "name" {
  description = "Fully qualified policy resource name"
  value       = google_dataproc_autoscaling_policy.policy.name
}

output "policy_id" {
  description = "The policy ID"
  value       = google_dataproc_autoscaling_policy.policy.policy_id
}

# The plain spec region name (not a provider-computed attribute), so API
# callers and verifiers can address the policy without parsing paths.
output "location" {
  description = "Region the policy lives in"
  value       = var.spec.location
}
