output "account_id" {
  description = "The Cloudflare account the default profile was applied to -- the profile is an account singleton, so the account is its identity"
  value       = var.spec.account_id
}

output "gateway_unique_id" {
  description = "The Gateway-side identifier Cloudflare assigns the profile"
  value       = cloudflare_zero_trust_device_default_profile.main.gateway_unique_id
}

output "policy_id" {
  description = "The profile's policy identifier as reported by the device policy API"
  value       = cloudflare_zero_trust_device_default_profile.main.policy_id
}
