output "auth_domain" {
  description = "The team domain users sign in through, without the .cloudflareaccess.com suffix"
  value       = cloudflare_zero_trust_organization.main.auth_domain
}

output "account_id" {
  description = "The Cloudflare account the organization was applied to (empty for a zone-scoped organization)"
  value       = var.spec.account_id
}
