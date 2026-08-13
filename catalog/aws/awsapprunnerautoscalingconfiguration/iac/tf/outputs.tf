output "configuration_arn" {
  description = "The revision-carrying ARN of this auto scaling configuration -- the identifier App Runner services reference."
  value       = aws_apprunner_auto_scaling_configuration_version.this.arn
}

output "configuration_revision" {
  description = "The revision this deployment registered under the configuration name."
  value       = aws_apprunner_auto_scaling_configuration_version.this.auto_scaling_configuration_revision
}

output "latest" {
  description = "Whether this revision is the latest one registered under the configuration name."
  value       = aws_apprunner_auto_scaling_configuration_version.this.latest
}

output "is_default" {
  description = "Whether this configuration holds the account/region default designation at the end of the deployment. Derived from the designation resource: its successful apply IS the claim."
  value       = length(aws_apprunner_default_auto_scaling_configuration_version.this) > 0
}
