output "app_id" {
  description = "App Platform application UUID"
  value       = digitalocean_app.main.id
}

output "default_hostname" {
  description = "Default ondigitalocean.app hostname (no scheme)"
  value       = replace(replace(digitalocean_app.main.default_ingress, "https://", ""), "http://", "")
}

output "live_url" {
  description = "Public URL of the app, including https://"
  value       = digitalocean_app.main.live_url
}

output "live_domain" {
  description = "Live domain hostname without scheme"
  value       = digitalocean_app.main.live_domain
}

output "active_deployment_id" {
  description = "UUID of the currently live deployment"
  value       = digitalocean_app.main.active_deployment_id
}
