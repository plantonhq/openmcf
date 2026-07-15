output "pipeline_arn" {
  description = "ARN of the CodePipeline pipeline -- the handle IAM policies and EventBridge targets reference"
  value       = aws_codepipeline.this.arn
}

output "pipeline_name" {
  description = "Name of the CodePipeline pipeline -- what CLI commands and other pipelines' action configurations reference"
  value       = aws_codepipeline.this.name
}
