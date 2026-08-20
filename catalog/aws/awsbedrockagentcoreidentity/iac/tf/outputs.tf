output "workload_identity_arns" {
  description = "Workload identity ARNs keyed by each workload_identities entry's name"
  value       = { for k, w in aws_bedrockagentcore_workload_identity.this : k => w.workload_identity_arn }
}

output "api_key_provider_arns" {
  description = "Credential provider ARNs keyed by each api_key_credential_providers entry's name"
  value       = { for k, p in aws_bedrockagentcore_api_key_credential_provider.this : k => p.credential_provider_arn }
}

output "api_key_secret_arns" {
  description = "Secrets Manager secret ARNs holding each vaulted API key, keyed by provider name"
  value       = { for k, p in aws_bedrockagentcore_api_key_credential_provider.this : k => one(p.api_key_secret_arn[*].secret_arn) }
}

output "oauth2_provider_arns" {
  description = "Credential provider ARNs keyed by each oauth2_credential_providers entry's name"
  value       = { for k, p in aws_bedrockagentcore_oauth2_credential_provider.this : k => p.credential_provider_arn }
}

output "oauth2_client_secret_arns" {
  description = "Secrets Manager secret ARNs holding each vaulted OAuth client secret, keyed by provider name"
  value       = { for k, p in aws_bedrockagentcore_oauth2_credential_provider.this : k => one(p.client_secret_arn[*].secret_arn) }
}

output "policy_engine_id" {
  description = "The policy engine's unique identifier (empty when the bundle has no engine)"
  value       = local.has_policy_engine ? aws_bedrockagentcore_policy_engine.this[0].policy_engine_id : ""
}

output "policy_engine_arn" {
  description = "The policy engine's ARN"
  value       = local.has_policy_engine ? aws_bedrockagentcore_policy_engine.this[0].policy_engine_arn : ""
}

output "policy_ids" {
  description = "Cedar policy IDs keyed by each policy_engine.policies entry's name"
  value       = { for k, p in aws_bedrockagentcore_policy.this : k => p.policy_id }
}
