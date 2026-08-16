output "bucket_name" {
  description = "The bucket's FULL name ({base}--{zone_id}--x-s3) - derived by the module, what S3 Express clients address, and the provider's import ID"
  value       = aws_s3_directory_bucket.this.bucket
}

output "bucket_arn" {
  description = "The bucket's ARN"
  value       = aws_s3_directory_bucket.this.arn
}
