output "vpc_endpoint_id" {
  description = "The endpoint's id (vpce-...), the handle every AWS API and route inspection uses."
  value       = aws_vpc_endpoint.this.id
}

output "arn" {
  description = "The endpoint's Amazon Resource Name."
  value       = aws_vpc_endpoint.this.arn
}

output "state" {
  description = "The endpoint's lifecycle state -- available on a successful create (pendingAcceptance when the PrivateLink service requires manual acceptance)."
  value       = aws_vpc_endpoint.this.state
}

output "prefix_list_id" {
  description = "The service's prefix list (gateway endpoints only) -- reference it in security-group or route rules scoped to the service's address ranges. Empty for ENI-based endpoint types."
  value       = try(coalesce(aws_vpc_endpoint.this.prefix_list_id, ""), "")
}

output "dns_name" {
  description = "The endpoint's primary private DNS name (interface endpoints only) -- what clients use when private DNS is off and what Route53 aliases target. Empty for gateway endpoints."
  value       = try(aws_vpc_endpoint.this.dns_entry[0].dns_name, "")
}

output "hosted_zone_id" {
  description = "The Route53 hosted zone of dns_name, needed alongside it for alias records (interface endpoints only)."
  value       = try(aws_vpc_endpoint.this.dns_entry[0].hosted_zone_id, "")
}

output "network_interface_ids" {
  description = "The endpoint's ENIs, one per attached subnet -- what flow logs, firewall rules, and IP lookups operate on. Empty for gateway endpoints."
  value       = aws_vpc_endpoint.this.network_interface_ids
}
