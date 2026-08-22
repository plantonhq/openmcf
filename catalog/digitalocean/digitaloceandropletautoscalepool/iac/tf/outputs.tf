# Stack outputs — exactly the DigitalOceanDropletAutoscalePoolStackOutputs
# contract, identical across both provisioners.

output "pool_id" {
  description = "UUID of the autoscale pool (its API identity and import id)"
  value       = digitalocean_droplet_autoscale.pool.id
}

output "status" {
  description = "Health status of the pool as reported by DigitalOcean at apply time"
  value       = digitalocean_droplet_autoscale.pool.status
}
