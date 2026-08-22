output "identity_provider_id" {
  description = "The UUID of the identity provider (what Access policy rules reference)"
  value       = cloudflare_zero_trust_access_identity_provider.main.id
}

output "scim_base_url" {
  description = "The base URL of Cloudflare's SCIM v2.0 endpoint for this provider (present when SCIM is enabled)"
  value       = cloudflare_zero_trust_access_identity_provider.main.scim_config != null ? cloudflare_zero_trust_access_identity_provider.main.scim_config.scim_base_url : null
}

output "scim_secret" {
  description = "The SCIM bearer token minted when SCIM is first enabled. Returned only once -- capture it into a secret store now; it is redacted on later reads and does not survive import."
  value       = cloudflare_zero_trust_access_identity_provider.main.scim_config != null ? cloudflare_zero_trust_access_identity_provider.main.scim_config.secret : null
  sensitive   = true
}
