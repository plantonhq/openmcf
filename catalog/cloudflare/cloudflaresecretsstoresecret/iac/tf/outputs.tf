output "secret_id" {
  description = "The secret's ID within its store"
  value       = cloudflare_secrets_store_secret.main.id
}

output "store_id" {
  description = "The ID of the store holding the secret -- echoed for consumers that need the store/secret pair"
  value       = cloudflare_secrets_store_secret.main.store_id
}
