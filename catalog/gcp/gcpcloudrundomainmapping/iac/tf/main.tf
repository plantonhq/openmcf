# Enables the Cloud Run Admin API, which serves domain mappings.
# disable_on_destroy stays false: tearing down one mapping must never
# disable Cloud Run for everything else in the project.
resource "google_project_service" "run_api" {
  project = local.project_id
  service = "run.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The provider's resolved default project — read ONLY when both
# spec.namespace and spec.project_id are empty (the mapping's required
# namespace has nothing else to fall back to), so ordinary project-named
# deployments plan credential-free offline.
data "google_client_config" "default" {
  count = local.namespace_needs_client_config ? 1 : 0
}

# The domain mapping: points a verified custom domain at a Cloud Run
# service and (in AUTOMATIC mode) has Cloud Run provision and renew the
# TLS certificate.
#
# The resource is fully IMMUTABLE — every argument below is create-only,
# so any change replaces the mapping (cheap: the object is free and
# re-creates in seconds, with a brief serving gap while the certificate
# re-issues). The domain must already be verified by the provisioning
# identity; GCP rejects the create otherwise.
resource "google_cloud_run_domain_mapping" "this" {
  # The mapping's name IS the domain being mapped.
  name     = var.spec.domain
  location = var.spec.region
  project  = local.project_id

  # The v1 API requires exactly one metadata block whose namespace equals
  # the project ID or project number.
  metadata {
    namespace = local.namespace
    labels    = local.final_labels
    # Non-authoritative: the Cloud Run API adds server-side annotations of
    # its own; the provider manages only the entries declared here.
    annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null
  }

  spec {
    # The Cloud Run service this domain routes to. Must exist, in this
    # same region and project, before the mapping is created.
    route_name = var.spec.route

    # AUTOMATIC (managed certificate) is the provider default; sent only
    # when the spec sets a value so the default stays provider-owned.
    certificate_mode = var.spec.certificate_mode != "" ? var.spec.certificate_mode : null

    # Sent only when true: unset preserves GCP's safe conflict error when
    # the domain is already mapped elsewhere.
    force_override = var.spec.force_override ? true : null
  }

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.run_api]
}
