output "vpc_link_id" {
  description = "The VPC link ID (what integrations set as their connection)"
  value       = aws_api_gateway_vpc_link.this.id
}

output "vpc_link_arn" {
  description = "The VPC link ARN"
  value       = aws_api_gateway_vpc_link.this.arn
}
