# Enable the Network Connectivity API — the control plane that owns
# service connection policies. disable_on_destroy is false: tearing down
# one policy must never disable the API for everything else in the
# project.
resource "google_project_service" "networkconnectivity_api" {
  project = local.project_id
  service = "networkconnectivity.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# Enable the Compute Engine API — the network and subnets the policy
# points at live in Compute, and the automation's forwarding rules are
# Compute-side objects.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The service connection policy: the per-network authorization that lets
# Google's service connectivity automation create PSC endpoints in the
# listed subnets on a producer's behalf. PSC-first managed services
# (Memorystore for Valkey, Redis Cluster) refuse to create instances on
# a network until a policy for their service class exists in the
# instance's region — this resource is that prerequisite.
#
# Cardinality is one policy per (network, service class, region) triple;
# GCP rejects a second. location, network, service_class, and the policy
# name are all immutable (ForceNew) — only the psc_config contents,
# description, and labels update in place. Keep the policy alive as long
# as any instance depends on it: deleting it strands existing endpoints
# and blocks new ones.
resource "google_network_connectivity_service_connection_policy" "this" {
  name          = local.policy_name
  project       = local.project_id
  location      = var.spec.location
  network       = local.network
  service_class = var.spec.service_class
  description   = local.description
  labels        = local.final_labels

  dynamic "psc_config" {
    for_each = var.spec.psc_config != null ? [var.spec.psc_config] : []
    content {
      subnetworks                = local.subnetworks
      limit                      = local.psc_limit
      producer_instance_location = local.producer_instance_location

      allowed_google_producers_resource_hierarchy_level = local.allowed_hierarchy_levels
    }
  }

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [
    google_project_service.networkconnectivity_api,
    google_project_service.compute_api,
  ]
}
