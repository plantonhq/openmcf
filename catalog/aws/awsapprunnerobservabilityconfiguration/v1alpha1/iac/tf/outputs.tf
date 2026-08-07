output "configuration_arn" {
  description = "The revision-carrying ARN of this observability configuration -- the identifier App Runner services reference."
  value       = aws_apprunner_observability_configuration.this.arn
}

output "configuration_revision" {
  description = "The revision this deployment registered under the configuration name."
  value       = aws_apprunner_observability_configuration.this.observability_configuration_revision
}

output "latest" {
  description = "Whether this revision is the latest one registered under the configuration name."
  value       = aws_apprunner_observability_configuration.this.latest
}
