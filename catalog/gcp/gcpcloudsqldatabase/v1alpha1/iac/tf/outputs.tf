output "database_name" {
  description = "Name of the database inside the instance"
  value       = google_sql_database.this.name
}

output "self_link" {
  description = "Self-link URL of the database resource"
  value       = google_sql_database.this.self_link
}
