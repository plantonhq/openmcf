output "service_token_id" {
  description = "The UUID of the service token (the API identity used for import and policy service_token rules)"
  value       = cloudflare_zero_trust_access_service_token.main.id
}

output "client_id" {
  description = "The Client ID presented in the CF-Access-Client-ID request header"
  value       = cloudflare_zero_trust_access_service_token.main.client_id
}

output "client_secret" {
  description = "The Client Secret presented in the CF-Access-Client-Secret request header. Returned only at creation and rotation -- capture it into a secret store now; it can never be read back."
  value       = cloudflare_zero_trust_access_service_token.main.client_secret
  sensitive   = true
}

output "expires_at" {
  description = "When the token expires (RFC3339)"
  value       = cloudflare_zero_trust_access_service_token.main.expires_at
}
