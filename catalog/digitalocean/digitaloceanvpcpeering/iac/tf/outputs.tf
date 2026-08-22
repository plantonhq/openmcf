# Stack outputs — exactly the DigitalOceanVpcPeeringStackOutputs contract,
# identical across both provisioners.

output "peering_id" {
  description = "UUID of the VPC peering connection (its API identity and import id)"
  value       = digitalocean_vpc_peering.peering.id
}

output "status" {
  description = "Lifecycle status of the peering as reported by DigitalOcean at apply time (UPPERCASE, e.g. ACTIVE)"
  value       = digitalocean_vpc_peering.peering.status
}
