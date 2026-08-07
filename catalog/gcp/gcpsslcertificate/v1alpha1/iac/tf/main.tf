# Enable the Compute Engine API so a fresh project can host the certificate.
# disable_on_destroy is false: tearing down one certificate must never disable
# the API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A self-managed Compute Engine SSL certificate — you bring the PEM chain and
# private key; the load balancer presents them to clients. Target HTTPS (and
# SSL) proxies reference the certificate's self_link exactly like a
# Google-managed certificate (the two share one API collection and name
# namespace), but nothing here renews itself: track expire_time and rotate.
#
# One kind, two provider resources: GCP models global and regional SSL
# certificates as separate API collections with an identical surface, so this
# module creates google_compute_ssl_certificate when spec.region is empty and
# google_compute_region_ssl_certificate when it is set. Exactly one of the
# two count guards is 1 — never both.
#
# EVERY argument is immutable (ForceNew): any change destroys and recreates
# the certificate. Because a certificate attached to a proxy cannot be
# deleted (GCP returns resourceInUseByAnotherResource), rotation is
# create-before-destroy: create the replacement under a new name, repoint the
# proxy's certificate list, then destroy this one.
#
# The private key is write-only in GCP: the API never returns it, the
# provider tracks it by hash, and this module never surfaces it in outputs.
resource "google_compute_ssl_certificate" "this" {
  count = local.is_regional ? 0 : 1

  name        = local.certificate_name
  project     = local.project_id
  description = var.spec.description != "" ? var.spec.description : null

  certificate = var.spec.certificate
  private_key = var.spec.private_key # secret material; never surfaced in outputs

  depends_on = [google_project_service.compute_api]
}

# The regional variant — identical surface, addressed under
# regions/<region>/sslCertificates. Regional proxies can only reference
# certificates in their own region.
resource "google_compute_region_ssl_certificate" "this" {
  count = local.is_regional ? 1 : 0

  name        = local.certificate_name
  project     = local.project_id
  region      = var.spec.region
  description = var.spec.description != "" ? var.spec.description : null

  certificate = var.spec.certificate
  private_key = var.spec.private_key # secret material; never surfaced in outputs

  depends_on = [google_project_service.compute_api]
}
