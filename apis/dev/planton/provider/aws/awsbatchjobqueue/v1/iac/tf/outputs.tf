output "job_queue_arn" {
  description = "The ARN of the job queue -- the handle jobs are submitted against and EventBridge Batch targets point at."
  value       = aws_batch_job_queue.this.arn
}

output "job_queue_name" {
  description = "The job queue's name (derived from metadata.name)."
  value       = aws_batch_job_queue.this.name
}
