# Output keys and shapes must match outputs.proto — the outputs
# transformer maps them onto the proto by name (list entries by
# dot-indexed keys).

output "domain" {
  description = "The mapped domain (the mapping's name in GCP)."
  value       = google_cloud_run_domain_mapping.this.name
}

output "region" {
  description = "GCP region the mapping lives in."
  value       = google_cloud_run_domain_mapping.this.location
}

output "resource_records" {
  description = "DNS records the domain's zone must publish for the mapping to serve (A/AAAA for a root domain, one CNAME for a subdomain)."
  value = [
    for record in try(google_cloud_run_domain_mapping.this.status[0].resource_records, []) : {
      record_type = record.type
      record_name = record.name
      rrdata      = record.rrdata
    }
  ]
}

output "mapped_route_name" {
  description = "The Cloud Run route (service) the mapping currently points to, as reported by GCP."
  value       = try(google_cloud_run_domain_mapping.this.status[0].mapped_route_name, "")
}
