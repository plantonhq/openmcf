output "certificate-id" {
  description = "The ID of the created certificate"
  value = local.is_managed ? (
    length(google_certificate_manager_certificate.cert) > 0 ? google_certificate_manager_certificate.cert[0].id : ""
  ) : (
    length(google_compute_managed_ssl_certificate.lb_cert) > 0 ? google_compute_managed_ssl_certificate.lb_cert[0].id : ""
  )
}

output "certificate-name" {
  description = "The name of the created certificate"
  value = local.is_managed ? (
    length(google_certificate_manager_certificate.cert) > 0 ? google_certificate_manager_certificate.cert[0].name : ""
  ) : (
    length(google_compute_managed_ssl_certificate.lb_cert) > 0 ? google_compute_managed_ssl_certificate.lb_cert[0].name : ""
  )
}

output "certificate-domain-name" {
  description = "The primary domain name of the certificate"
  value       = var.spec.primary_domain_name
}

output "certificate-status" {
  description = "The status of the certificate"
  value       = "PROVISIONING"
}

