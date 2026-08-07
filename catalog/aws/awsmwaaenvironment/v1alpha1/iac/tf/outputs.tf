output "environment_arn" {
  description = "ARN of the MWAA environment"
  value       = aws_mwaa_environment.this.arn
}

output "environment_name" {
  description = "Name of the MWAA environment"
  value       = aws_mwaa_environment.this.name
}

output "webserver_url" {
  description = "URL of the Airflow web UI"
  value       = aws_mwaa_environment.this.webserver_url
}

output "airflow_version" {
  description = "Effective Apache Airflow version"
  value       = aws_mwaa_environment.this.airflow_version
}

output "service_role_arn" {
  description = "ARN of the AWS service role for the environment"
  value       = aws_mwaa_environment.this.service_role_arn
}

output "environment_class" {
  description = "Effective environment class"
  value       = aws_mwaa_environment.this.environment_class
}

output "status" {
  description = "Current status of the MWAA environment"
  value       = aws_mwaa_environment.this.status
}

output "created_at" {
  description = "Timestamp when the environment was created"
  value       = aws_mwaa_environment.this.created_at
}

output "database_vpc_endpoint_service" {
  description = "VPC endpoint service name for the metadata database (used with endpoint_management = CUSTOMER)"
  value       = aws_mwaa_environment.this.database_vpc_endpoint_service
}

output "webserver_vpc_endpoint_service" {
  description = "VPC endpoint service name for the webserver (used with endpoint_management = CUSTOMER; empty when PUBLIC_ONLY)"
  value       = aws_mwaa_environment.this.webserver_vpc_endpoint_service
}
