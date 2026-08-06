# Fully qualified table resource path — the canonical identifier for
# Admin API calls and IAM bindings.
output "table_id" {
  description = "Fully qualified table resource path"
  value       = google_bigtable_table.this.id
}

# Short table name — what Bigtable client libraries open, together with
# the project and instance.
output "table_name" {
  description = "Short table name"
  value       = google_bigtable_table.this.name
}

# The parent instance, confirmed without chasing the reference chain.
output "instance_name" {
  description = "Short name of the instance the table lives in"
  value       = var.spec.instance
}
