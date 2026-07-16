output "authorization_id" {
  description = "Fully-qualified resource ID (projects/{project}/locations/{location}/dnsAuthorizations/{name}) — the value a certificate's dns_authorizations list consumes"
  value       = google_certificate_manager_dns_authorization.authorization.id
}

output "authorization_name" {
  description = "Name of the authorization as it exists in GCP"
  value       = google_certificate_manager_dns_authorization.authorization.name
}

output "domain" {
  description = "The authorized domain (covers the domain and its wildcard)"
  value       = google_certificate_manager_dns_authorization.authorization.domain
}

output "dns_record_name" {
  description = "Fully-qualified name of the DNS validation record to create"
  value       = google_certificate_manager_dns_authorization.authorization.dns_resource_record[0].name
}

output "dns_record_type" {
  description = "Type of the DNS validation record (CNAME)"
  value       = google_certificate_manager_dns_authorization.authorization.dns_resource_record[0].type
}

output "dns_record_data" {
  description = "Data of the DNS validation record — the value the CNAME points at"
  value       = google_certificate_manager_dns_authorization.authorization.dns_resource_record[0].data
}
