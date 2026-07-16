# Enable the Compute Engine API so a fresh project can host the certificate.
# disable_on_destroy is false: tearing down one certificate must never disable
# the API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Google-managed SSL certificate — a global Compute Engine SSL certificate
# whose private key and issuance are handled entirely by Google. Attach its
# self_link to a target HTTPS proxy to terminate TLS at a global external
# Application Load Balancer.
#
# The whole resource is immutable (name and domains are ForceNew): any change
# destroys and recreates the certificate. Because a cert attached to a proxy
# cannot be deleted, rotate by creating the replacement first and swapping the
# proxy's ssl_certificates reference before destroying the old one
# (create-before-destroy) — otherwise the destroy fails and TLS drops during
# the gap.
#
# Provisioning is asynchronous and DNS-gated: creation returns immediately, but
# the certificate stays PROVISIONING until each domain's DNS points at the load
# balancer's IP. expire_time stays empty until provisioning completes.
resource "google_compute_managed_ssl_certificate" "this" {
  name        = local.certificate_name
  project     = local.project_id
  description = var.spec.description != "" ? var.spec.description : null

  managed {
    domains = var.spec.domains
  }

  depends_on = [google_project_service.compute_api]
}
