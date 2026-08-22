output "firewall_id" {
  description = "The unique ID of the firewall (a UUID assigned by DigitalOcean)"
  value       = digitalocean_firewall.firewall.id
}
