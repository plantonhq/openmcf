output "store_id" {
  description = "The store's ID -- what store secrets, Worker bindings, and AI Gateway authentication reference"
  value       = cloudflare_secrets_store.main.id
}
