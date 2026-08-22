output "portal_id" {
  description = "The portal's identifier (the user-chosen slug)"
  value       = cloudflare_zero_trust_access_ai_controls_mcp_portal.main.id
}

output "hostname" {
  description = "The hostname the portal is served on -- what users and AI clients connect to"
  value       = cloudflare_zero_trust_access_ai_controls_mcp_portal.main.hostname
}
