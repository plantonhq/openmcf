output "identity_provider_id" {
  description = "The UUID of the identity provider (what Access policy rules reference)"
  value       = cloudflare_zero_trust_access_identity_provider.main.id
}

# Direct attribute access via try(), NEVER a `scim_config != null` guard:
# scim_config contains a sensitive member (secret), so comparing the whole
# object yields a SENSITIVE boolean and the conditional taints the result --
# tofu then refuses the unmarked output at apply time ("Output refers to
# sensitive values", measured live 2026-08-26). try() reads the non-sensitive
# leaf directly and degrades to null when scim_config is absent.
output "scim_base_url" {
  description = "The base URL of Cloudflare's SCIM v2.0 endpoint for this provider (present when SCIM is enabled)"
  value       = try(cloudflare_zero_trust_access_identity_provider.main.scim_config.scim_base_url, null)
}

output "scim_secret" {
  description = "The SCIM bearer token minted when SCIM is first enabled. Returned only once -- capture it into a secret store now; it is redacted on later reads and does not survive import."
  value       = try(cloudflare_zero_trust_access_identity_provider.main.scim_config.secret, null)
  sensitive   = true
}
