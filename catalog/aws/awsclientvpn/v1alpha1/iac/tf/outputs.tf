output "client_vpn_endpoint_id" {
  description = "The AWS-assigned identifier for the Client VPN endpoint -- used for AWS CLI/API operations, including exporting the client configuration."
  value       = aws_ec2_client_vpn_endpoint.this.id
}

output "client_vpn_endpoint_arn" {
  description = "The Amazon Resource Name of the endpoint -- used in IAM policies and cross-service permissions."
  value       = aws_ec2_client_vpn_endpoint.this.arn
}

output "endpoint_dns_name" {
  description = "The DNS name clients connect to (the exported client configuration prepends the required random subdomain label automatically)."
  value       = aws_ec2_client_vpn_endpoint.this.dns_name
}

output "self_service_portal_url" {
  description = "The self-service portal URL where federated users download their own client configuration -- empty when the portal is disabled."
  value       = aws_ec2_client_vpn_endpoint.this.self_service_portal_url
}

output "subnet_association_ids" {
  description = "Map of subnet ID to target network association ID for each associated subnet."
  value       = { for subnet_id, assoc in aws_ec2_client_vpn_network_association.this : subnet_id => assoc.association_id }
}

output "transit_gateway_attachment_id" {
  description = "The transit gateway attachment created for a transit-gateway-associated endpoint -- empty for VPC-attached endpoints."
  value       = try(aws_ec2_client_vpn_endpoint.this.transit_gateway_configuration[0].transit_gateway_attachment_id, "")
}
