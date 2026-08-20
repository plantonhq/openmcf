output "policy_id" {
  description = "The Cloudflare-assigned identifier of the profile (its policy ID)"
  value       = cloudflare_zero_trust_device_custom_profile.main.id
}

output "gateway_unique_id" {
  description = "The Gateway-side identifier Cloudflare assigns the profile"
  value       = cloudflare_zero_trust_device_custom_profile.main.gateway_unique_id
}
