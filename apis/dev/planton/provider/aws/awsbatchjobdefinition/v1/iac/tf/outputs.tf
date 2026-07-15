output "job_definition_arn" {
  description = "The full ARN of the registered revision (name:revision) -- changes on every revision, rolling referencing consumers."
  value       = aws_batch_job_definition.this.arn
}

output "arn_without_revision" {
  description = "The ARN without the revision suffix, for consumers that track the name's latest ACTIVE revision."
  value       = aws_batch_job_definition.this.arn_prefix
}

output "job_definition_name" {
  description = "The job definition name (metadata.name) revisions register under."
  value       = aws_batch_job_definition.this.name
}

output "revision" {
  description = "The revision number this deployment registered."
  value       = aws_batch_job_definition.this.revision
}
