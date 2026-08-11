# The configuration's resource name: projects/{project}/config.
output "config_name" {
  description = "The resource name of the Identity Platform configuration"
  value       = google_identity_platform_config.this.name
}

# The auto-provisioned API key client apps initialize the Identity Platform
# SDK with. A live credential in the API-key sense — marked sensitive so it
# never prints in plans; restrict it by domain/app in the console.
output "api_key" {
  description = "The API key for client-side Identity Platform SDK initialization"
  value       = try(google_identity_platform_config.this.client[0].api_key, "")
  sensitive   = true
}

# The project's Firebase subdomain — the default hosted sign-in domain.
output "firebase_subdomain" {
  description = "The Firebase subdomain of the project"
  value       = try(google_identity_platform_config.this.client[0].firebase_subdomain, "")
}
