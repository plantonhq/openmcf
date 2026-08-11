output "map_id" {
  description = "Full certificate map resource name"
  value       = local.map_id
}

# The form a GcpTargetHttpsProxy's certificate_map argument consumes.
output "map_uri" {
  description = "The //certificatemanager.googleapis.com/... map URI"
  value       = "//certificatemanager.googleapis.com/${local.map_id}"
}

output "map_name" {
  description = "The short map name"
  value       = google_certificate_manager_certificate_map.this.name
}
