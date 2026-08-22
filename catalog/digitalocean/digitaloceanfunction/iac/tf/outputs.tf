output "function_id" {
  description = "App Platform application UUID that hosts the functions component"
  value       = digitalocean_app.main.id
}

output "https_endpoint" {
  description = "Public HTTPS URL of the app"
  value       = digitalocean_app.main.live_url
}

output "default_hostname" {
  description = "Default ondigitalocean.app hostname (no scheme)"
  value       = replace(replace(digitalocean_app.main.default_ingress, "https://", ""), "http://", "")
}
