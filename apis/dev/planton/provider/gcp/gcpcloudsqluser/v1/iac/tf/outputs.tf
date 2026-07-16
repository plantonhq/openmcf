output "user_name" {
  description = "The user name as stored by Cloud SQL (IAM users on MySQL are stored truncated before the @)"
  value       = google_sql_user.this.name
}

output "instance_name" {
  description = "Name of the Cloud SQL instance this user belongs to"
  value       = google_sql_user.this.instance
}
