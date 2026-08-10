# The server-generated tenant ID (the last segment of tenant_name) — what
# client SDKs set as the tenantId to scope sign-in to this tenant.
output "tenant_id" {
  description = "The server-generated tenant ID"
  value       = google_identity_platform_tenant.this.name
}

# The tenant's full resource name: projects/{project}/tenants/{tenant_id}.
output "tenant_name" {
  description = "The full resource name of the tenant"
  value       = google_identity_platform_tenant.this.id
}
