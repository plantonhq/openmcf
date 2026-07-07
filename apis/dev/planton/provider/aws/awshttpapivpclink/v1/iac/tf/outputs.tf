output "vpc_link_id" {
  description = "The VPC link ID -- the identifier private integrations set as their connection_id."
  value       = aws_apigatewayv2_vpc_link.this.id
}

output "vpc_link_arn" {
  description = "The ARN of the VPC link."
  value       = aws_apigatewayv2_vpc_link.this.arn
}
