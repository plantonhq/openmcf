output "attachment_id" {
  description = "The Transit Gateway attachment ID -- the join key referenced by route table associations, propagations, and static routes."
  value       = aws_ec2_transit_gateway_vpc_attachment.this.id
}

output "attachment_arn" {
  description = "The ARN of the attachment, for IAM policies and resource-level permissions."
  value       = aws_ec2_transit_gateway_vpc_attachment.this.arn
}

output "vpc_owner_id" {
  description = "The AWS account ID that owns the attached VPC (differs from the gateway owner in cross-account topologies)."
  value       = aws_ec2_transit_gateway_vpc_attachment.this.vpc_owner_id
}
