output "vpc_connector_arn" {
  description = "The ARN of the VPC connector -- the identifier App Runner services set as their egress vpc_connector_arn."
  value       = aws_apprunner_vpc_connector.this.arn
}

output "vpc_connector_revision" {
  description = "The revision of this connector under its name."
  value       = aws_apprunner_vpc_connector.this.vpc_connector_revision
}

output "status" {
  description = "The connector's lifecycle status at the end of the deployment (ACTIVE when ready for services to attach)."
  value       = aws_apprunner_vpc_connector.this.status
}
