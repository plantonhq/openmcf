output "account_id" {
  description = "The Cloudflare account the Gateway configuration was applied to (the singleton's identity)"
  value       = var.spec.account_id
}
