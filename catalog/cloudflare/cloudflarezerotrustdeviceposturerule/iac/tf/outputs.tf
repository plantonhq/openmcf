output "rule_id" {
  description = "The Cloudflare-assigned UUID of the posture rule (what Access and Gateway policies reference)"
  value       = cloudflare_zero_trust_device_posture_rule.main.id
}
