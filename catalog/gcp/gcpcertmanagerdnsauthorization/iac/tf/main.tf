# Enable the Certificate Manager API so a fresh project can host DNS
# authorizations. disable_on_destroy is false: tearing down one
# authorization must never disable the API for everything else in the
# project.
resource "google_project_service" "certificatemanager_api" {
  project = local.project_id
  service = "certificatemanager.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# One DNS authorization: the proof-of-domain-control a Google-managed
# certificate consumes. It covers the domain AND its wildcard, and exports
# the CNAME validation record a GcpDnsRecord composes into the zone —
# validation can complete BEFORE any certificate exists, which is what
# makes zero-downtime certificate migration possible.
resource "google_certificate_manager_dns_authorization" "authorization" {
  name        = local.authorization_name
  project     = local.project_id
  domain      = var.spec.domain
  description = var.spec.description != "" ? var.spec.description : null
  location    = var.spec.location != "" ? var.spec.location : null
  type        = var.spec.type != "" ? var.spec.type : null
  labels      = length(local.final_labels) > 0 ? local.final_labels : null

  depends_on = [google_project_service.certificatemanager_api]
}
