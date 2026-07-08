output "configuration_set_arn" {
  description = "The ARN of the configuration set -- the target for IAM policies scoping who may send under it."
  value       = aws_sesv2_configuration_set.this.arn
}

output "configuration_set_name" {
  description = "The configuration set's name (derived from metadata.name) -- what email identities reference and SendEmail calls name."
  value       = aws_sesv2_configuration_set.this.configuration_set_name
}
