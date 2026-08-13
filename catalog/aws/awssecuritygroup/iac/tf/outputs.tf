output "security_group_id" {
  description = "The ID of the security group (sg-...) -- the join key other resources reference to attach it."
  value       = aws_security_group.this.id
}

output "security_group_arn" {
  description = "The ARN of the security group -- the form IAM policy conditions and resource-level permissions expect."
  value       = aws_security_group.this.arn
}

output "owner_id" {
  description = "The AWS account ID that owns the security group -- needed for cross-account rule references (<owner_id>/<group_id>)."
  value       = aws_security_group.this.owner_id
}

output "additional_vpc_association_ids" {
  description = "Association IDs of the group's shares into other VPCs, keyed by VPC id (import id form <group_id>,<vpc_id>); empty when the group lives in one VPC."
  value       = { for vpc_id, assoc in aws_vpc_security_group_vpc_association.additional : vpc_id => "${assoc.security_group_id},${assoc.vpc_id}" }
}
