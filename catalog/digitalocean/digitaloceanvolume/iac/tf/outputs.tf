output "volume_id" {
  description = "The unique identifier (UUID) of the created DigitalOcean volume."
  value       = digitalocean_volume.this.id
}

output "urn" {
  description = "The uniform resource name (URN) of the volume."
  value       = digitalocean_volume.this.urn
}
