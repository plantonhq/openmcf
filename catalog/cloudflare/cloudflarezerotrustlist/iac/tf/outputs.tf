output "list_id" {
  description = "The UUID of the created list (what Gateway policies and posture rules reference)"
  value       = cloudflare_zero_trust_list.main.id
}
