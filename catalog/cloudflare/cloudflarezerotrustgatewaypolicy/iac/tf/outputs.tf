output "policy_id" {
  description = "The UUID of the created Gateway policy"
  value       = cloudflare_zero_trust_gateway_policy.main.id
}

output "precedence" {
  description = "The policy's evaluation precedence (assigned by Cloudflare when the spec leaves it unset; lower runs earlier)"
  value       = cloudflare_zero_trust_gateway_policy.main.precedence
}
