output "route_table_id" {
  description = "The Transit Gateway route table ID (tgw-rtb-...)."
  value       = aws_ec2_transit_gateway_route_table.this.id
}

output "route_table_arn" {
  description = "The ARN of the route table, for IAM policies and resource-level permissions."
  value       = aws_ec2_transit_gateway_route_table.this.arn
}
