output "instance_name" {
  description = "Fully qualified instance resource name"
  value       = google_alloydb_instance.this.name
}

output "ip_address" {
  description = "Private IP address of the instance"
  value       = google_alloydb_instance.this.ip_address
}

output "state" {
  description = "Current state of the AlloyDB instance"
  value       = google_alloydb_instance.this.state
}
