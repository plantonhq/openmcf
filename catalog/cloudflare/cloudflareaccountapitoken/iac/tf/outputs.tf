output "token_id" {
  description = "The Cloudflare-assigned token ID (the token's identity for management calls, not the credential)"
  value       = cloudflare_account_token.main.id
}

output "value" {
  description = "The token's secret value -- returned by Cloudflare exactly once, on create; if lost, rotate"
  value       = cloudflare_account_token.main.value
  sensitive   = true
}
