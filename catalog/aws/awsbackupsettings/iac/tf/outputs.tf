output "account_id" {
  description = "The AWS account ID the global arm manages (the global settings resource's identity at AWS)"
  value       = data.aws_caller_identity.this.account_id
}

output "region" {
  description = "The region the region arm manages (the region settings resource's identity at AWS)"
  value       = data.aws_region.this.region
}
