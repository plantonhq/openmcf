output "server_id" {
  description = "The server's identifier -- what MCP portals reference in their servers[] rows"
  value       = cloudflare_zero_trust_access_ai_controls_mcp_server.main.id
}
