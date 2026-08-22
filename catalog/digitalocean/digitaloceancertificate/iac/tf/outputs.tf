output "certificate_id" {
  description = "The certificate's resource identifier -- the certificate NAME at the current provider pin (a Let's Encrypt certificate's UUID rotates on every auto-renewal, so DigitalOcean addresses certificates by name)."
  value       = digitalocean_certificate.certificate.id
}

output "expiry_rfc3339" {
  description = "The expiration timestamp of the certificate in RFC 3339 format."
  value       = digitalocean_certificate.certificate.not_after
}
