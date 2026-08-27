output "account_id" {
  description = "The Cloudflare account the Gateway configuration was applied to (the singleton's identity)"
  value       = var.spec.account_id
}

output "pacfile_ids" {
  description = "Cloudflare-assigned ids of the PAC files, keyed by file name (the module's for_each key) -- import recipes derive per-file import IDs from it"
  value       = { for k, p in cloudflare_zero_trust_gateway_pacfile.main : k => p.id }
}
