output "transit_gateway_id" {
  description = "The Transit Gateway ID -- the join key referenced by VPC attachments, route tables, subnet routes, VPN connections, and Direct Connect gateways."
  value       = aws_ec2_transit_gateway.this.id
}

output "transit_gateway_arn" {
  description = "The ARN of the Transit Gateway, for IAM policies and AWS RAM sharing."
  value       = aws_ec2_transit_gateway.this.arn
}

output "owner_id" {
  description = "The AWS account ID that owns the Transit Gateway."
  value       = aws_ec2_transit_gateway.this.owner_id
}

output "association_default_route_table_id" {
  description = "The ID of the default association route table (empty when default association is disabled)."
  value       = aws_ec2_transit_gateway.this.association_default_route_table_id
}

output "propagation_default_route_table_id" {
  description = "The ID of the default propagation route table (empty when default propagation is disabled)."
  value       = aws_ec2_transit_gateway.this.propagation_default_route_table_id
}
