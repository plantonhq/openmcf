# Enable the Certificate Manager API so a fresh project can host
# certificates. disable_on_destroy is false: tearing down one certificate
# must never disable the API for everything else in the project.
resource "google_project_service" "certificatemanager_api" {
  project = local.project_id
  service = "certificatemanager.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# One Certificate Manager certificate. Exactly one arm is configured
# (enforced pre-deploy):
#   - managed: Google provisions and renews automatically. Domain control
#     is proven via referenced DNS authorizations, a private-PKI issuance
#     config, or load-balancer authorization when neither is set.
#   - self_managed: uploaded PEM chain + private key; rotation is an
#     in-place update with new material.
resource "google_certificate_manager_certificate" "certificate" {
  name        = local.cert_name
  project     = local.project_id
  description = var.spec.description != "" ? var.spec.description : null
  location    = var.spec.location != "" ? var.spec.location : null
  scope       = var.spec.scope != "" ? var.spec.scope : null
  labels      = length(local.final_labels) > 0 ? local.final_labels : null

  depends_on = [google_project_service.certificatemanager_api]

  dynamic "managed" {
    for_each = var.spec.managed != null ? [var.spec.managed] : []
    content {
      domains            = managed.value.domains
      dns_authorizations = length(managed.value.dns_authorizations) > 0 ? managed.value.dns_authorizations : null
      issuance_config    = managed.value.issuance_config != "" ? managed.value.issuance_config : null
    }
  }

  dynamic "self_managed" {
    for_each = var.spec.self_managed != null ? [var.spec.self_managed] : []
    content {
      pem_certificate = self_managed.value.pem_certificate
      pem_private_key = self_managed.value.pem_private_key
    }
  }
}
