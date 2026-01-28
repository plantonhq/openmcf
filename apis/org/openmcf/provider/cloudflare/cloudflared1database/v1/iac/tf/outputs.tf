# outputs.tf

output "database_id" {
  description = "The unique identifier of the created D1 database"
  value       = cloudflare_d1_database.main.id
}

output "database_name" {
  description = "The name of the database (same as the input name)"
  value       = cloudflare_d1_database.main.name
}

output "connection_string" {
  description = "The connection string for connecting to the D1 database (currently not available in Terraform provider)"
  value       = ""
}

