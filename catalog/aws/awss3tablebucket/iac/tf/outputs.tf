output "table_bucket_arn" {
  description = "The table bucket's ARN - what catalog integrations, policies, and replication destinations reference, and the provider's import ID"
  value       = aws_s3tables_table_bucket.this.arn
}

output "owner_account_id" {
  description = "The AWS account that owns the bucket"
  value       = aws_s3tables_table_bucket.this.owner_account_id
}

output "table_arns" {
  description = "Each table's ARN, keyed namespace//table - what per-table policies and table-level replication reference"
  value       = { for key, table in aws_s3tables_table.this : key => table.arn }
}

output "table_warehouse_locations" {
  description = "Each table's warehouse location (s3://... metadata root), keyed namespace//table - what manually-configured query engines point at"
  value       = { for key, table in aws_s3tables_table.this : key => table.warehouse_location }
}
