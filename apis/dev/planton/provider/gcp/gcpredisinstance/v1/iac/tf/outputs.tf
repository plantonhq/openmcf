output "host" {
  description = "Hostname or IP address of the primary Redis endpoint"
  value       = google_redis_instance.this.host
}

output "port" {
  description = "Port number of the primary Redis endpoint"
  value       = google_redis_instance.this.port
}

output "read_endpoint" {
  description = "Hostname or IP address of the read replica endpoint (STANDARD_HA with read replicas only)"
  value       = google_redis_instance.this.read_endpoint
}

output "read_endpoint_port" {
  description = "Port number of the read replica endpoint"
  value       = google_redis_instance.this.read_endpoint_port
}

output "current_location_id" {
  description = "Zone where the Redis primary is currently running"
  value       = google_redis_instance.this.current_location_id
}

output "auth_string" {
  description = "Redis AUTH string (populated only when auth_enabled is true)"
  value       = google_redis_instance.this.auth_string
  sensitive   = true
}

output "server_ca_certs" {
  description = "PEM-encoded CA certificates clients must trust when transit encryption is enabled"
  value       = [for cert in google_redis_instance.this.server_ca_certs : cert.cert]
}

output "persistence_iam_identity" {
  description = "Cloud IAM identity (serviceAccount:<email>) used by import/export operations"
  value       = google_redis_instance.this.persistence_iam_identity
}

output "effective_reserved_ip_range" {
  description = "The CIDR range actually in use by the instance (explicit or auto-selected)"
  value       = google_redis_instance.this.effective_reserved_ip_range
}

output "instance_name" {
  description = "Name of the Redis instance in GCP"
  value       = google_redis_instance.this.name
}

output "region" {
  description = "Region hosting the instance (plain region name)"
  value       = var.spec.region
}
