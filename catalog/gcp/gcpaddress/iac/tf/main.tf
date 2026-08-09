# Enable the Compute Engine API so a fresh project can host the address.
# disable_on_destroy is false: tearing down one address must never disable
# the API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A regional Compute Engine address reservation — either a public static IP
# (EXTERNAL, for Cloud NAT, regional LBs, or VM instances) or a private IP
# or range (INTERNAL, for GCE endpoints, internal LB VIPs, VPC peering, or
# IPsec interconnect).
#
# Every field except labels is immutable (ForceNew): any change destroys and
# recreates the reservation — and a recreated EXTERNAL address is a NEW IP,
# so DNS pointing at the old one breaks. Reserve once, reference everywhere.
resource "google_compute_address" "this" {
  project     = local.project_id
  name        = local.address_name
  region      = var.spec.region
  description = local.description

  address_type = var.spec.address_type
  ip_version   = var.spec.ip_version

  # Empty lets GCP assign the IP / range start.
  address = local.address

  # INTERNAL-only wiring (spec CEL enforces the coherence pre-deploy).
  network       = local.network
  subnetwork    = local.subnetwork
  prefix_length = var.spec.prefix_length
  purpose       = local.purpose

  # EXTERNAL-only; spec CEL rejects network_tier on INTERNAL.
  network_tier       = local.network_tier
  ipv6_endpoint_type = local.ipv6_endpoint_type

  # BYOIP: reserve out of a customer-owned PublicDelegatedPrefix.
  ip_collection = local.ip_collection

  labels = local.final_labels

  # What destroy does to the reservation: DELETE (default), PREVENT
  # (refuse), or ABANDON (drop from state, keep the IP reserved).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.compute_api]
}
