output "zone_name" {
  description = "The domain name of the DNS zone"
  value       = digitalocean_domain.dns_zone.name
}

output "zone_id" {
  description = "The zone's resource identifier — DigitalOcean addresses domains by name, so this is the domain name itself"
  value       = digitalocean_domain.dns_zone.id
}

output "name_servers" {
  description = "DigitalOcean's authoritative name servers (a fixed platform-wide set the API does not return per zone); set these at the registrar to delegate"
  value = [
    "ns1.digitalocean.com",
    "ns2.digitalocean.com",
    "ns3.digitalocean.com"
  ]
}

output "urn" {
  description = "The uniform resource name of the domain (e.g. do:domain:example.com)"
  value       = digitalocean_domain.dns_zone.urn
}
