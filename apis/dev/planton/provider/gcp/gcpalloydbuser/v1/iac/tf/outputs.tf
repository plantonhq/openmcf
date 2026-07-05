output "name" {
  description = "Fully qualified user resource name"
  value       = google_alloydb_user.this.name
}

output "user_id" {
  description = "The user_id as stored by AlloyDB"
  value       = google_alloydb_user.this.user_id
}

output "cluster_id" {
  description = "Fully qualified cluster resource name"
  value       = google_alloydb_user.this.cluster
}
