# Enable the Compute Engine API so a fresh project can host the address.
# disable_on_destroy is false: tearing down one address must never disable
# the API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A global Compute Engine address reservation — either a public static IP
# (EXTERNAL, the frontend VIP a global forwarding rule binds) or a private
# range (INTERNAL, for VPC peering with managed services or Private Service
# Connect).
#
# Every field except labels is immutable (ForceNew): any change destroys and
# recreates the reservation — and a recreated EXTERNAL address is a NEW IP,
# so DNS pointing at the old one breaks. Reserve once, reference everywhere.
resource "google_compute_global_address" "this" {
  project     = local.project_id
  name        = local.address_name
  description = local.description

  address_type = var.spec.address_type
  ip_version   = var.spec.ip_version

  # Empty lets GCP assign the IP / range start.
  address = local.address

  # INTERNAL-only wiring (spec CEL enforces the coherence pre-deploy).
  network       = local.network
  prefix_length = var.spec.prefix_length
  purpose       = local.purpose

  labels = local.final_labels

  depends_on = [google_project_service.compute_api]
}
