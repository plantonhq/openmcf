output "workgroup_name" {
  description = "The workgroup name -- the handle the Redshift Serverless APIs and the credentials API address the workgroup by."
  value       = aws_redshiftserverless_workgroup.this.workgroup_name
}

output "workgroup_id" {
  description = "The unique identifier AWS assigns to the workgroup."
  value       = aws_redshiftserverless_workgroup.this.workgroup_id
}

output "arn" {
  description = "The workgroup's Amazon Resource Name, for IAM policies and usage limits."
  value       = aws_redshiftserverless_workgroup.this.arn
}

output "endpoint_address" {
  description = "The DNS hostname SQL clients connect to."
  value       = try(aws_redshiftserverless_workgroup.this.endpoint[0].address, "")
}

output "port" {
  description = "The port the workgroup accepts connections on."
  value       = aws_redshiftserverless_workgroup.this.port
}
