output "registry_name" {
  description = "The registry name (also the registry's resource identifier in DigitalOcean)."
  value       = digitalocean_container_registry.registry.name
}

output "server_url" {
  description = "The registry host, always \"registry.digitalocean.com\"."
  value       = digitalocean_container_registry.registry.server_url
}

output "endpoint" {
  description = "The full endpoint for docker push/pull, i.e. \"registry.digitalocean.com/<registry_name>\"."
  value       = digitalocean_container_registry.registry.endpoint
}

output "region" {
  description = "Region slug where the registry is hosted (reported by DigitalOcean, covering the DigitalOcean-chooses case)."
  value       = digitalocean_container_registry.registry.region
}

output "docker_credentials" {
  description = "Base64-encoded Docker config.json for this registry -- a SECRET. Empty when the spec's docker_credentials block is unset."
  value       = local.create_docker_credentials ? digitalocean_container_registry_docker_credentials.credentials[0].docker_credentials : ""
  sensitive   = true
}

output "credential_expiration_time" {
  description = "RFC 3339 timestamp at which the minted docker credentials expire. Empty when the spec's docker_credentials block is unset."
  value       = local.create_docker_credentials ? digitalocean_container_registry_docker_credentials.credentials[0].credential_expiration_time : ""
}
