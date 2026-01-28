output "vpc_id" {
  description = "The ID of the VPC"
  value       = aws_vpc.this.id
}

output "internet_gateway_id" {
  description = "The ID of the Internet Gateway"
  value       = aws_internet_gateway.this.id
}

output "public_subnet_ids" {
  description = "Map of all public subnets created keyed by their for_each keys"
  value       = { for k, v in aws_subnet.public : k => v.id }
}

output "private_subnet_ids" {
  description = "Map of all private subnets created keyed by their for_each keys"
  value       = { for k, v in aws_subnet.private : k => v.id }
}

output "public_route_table_id" {
  description = "The ID of the public route table"
  value       = aws_route_table.public.id
}

output "private_route_table_ids" {
  description = "Map of private route table IDs keyed by their corresponding private subnet keys"
  value       = { for k, v in aws_route_table.private : k => v.id }
}

output "nat_gateway_ids" {
  description = "Map of NAT gateway IDs keyed by AZ (only if NAT is enabled)"
  value       = { for k, v in aws_nat_gateway.this : k => v.id }
}

output "nat_gateway_eip_addresses" {
  description = "Map of NAT Gateway EIP addresses keyed by AZ (only if NAT is enabled)"
  value       = { for k, v in aws_eip.nat : k => v.public_ip }
}
