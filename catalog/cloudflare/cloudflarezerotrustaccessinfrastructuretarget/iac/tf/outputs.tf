output "target_id" {
  description = "The Cloudflare-assigned UUID of the target"
  value       = cloudflare_zero_trust_access_infrastructure_target.main.id
}
