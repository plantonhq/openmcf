output "pipe_arn" {
  description = "The pipe's ARN"
  value       = aws_pipes_pipe.this.arn
}

output "pipe_name" {
  description = "The pipe's name - the provider's import ID"
  value       = aws_pipes_pipe.this.name
}
