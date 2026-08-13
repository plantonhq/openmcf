# Enable the Service Networking API — the producer-side control plane this
# connection talks to. disable_on_destroy is false: tearing down one
# connection must never disable the API for everything else in the project.
resource "google_project_service" "servicenetworking_api" {
  project = local.project_id
  service = "servicenetworking.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# Enable the Compute Engine API — the network and the reserved peering
# ranges live in Compute, and the peering the connection creates is a
# Compute-side object.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The private services access connection: a VPC peering between this network
# and the service producer's network, carved from the reserved VPC_PEERING
# address ranges. This single resource is what turns "Cloud SQL with private
# IP" from an error into a working deployment — producers allocate service
# subnets out of the reserved ranges and route them over the peering.
#
# Cardinality is one connection per (network, service) pair — GCP rejects a
# second. Capacity grows by appending range names to reserved_peering_ranges
# on THIS resource (an in-place update that never disturbs subnets the
# producer already provisioned), never by adding another connection.
#
# network and service are immutable (ForceNew): changing either destroys and
# recreates the connection, severing private connectivity for every producer
# resource on it. Teardown ordering: GCP refuses to delete the connection
# while the producer still holds subnets — destroy the private-IP service
# instances (Cloud SQL, AlloyDB, Memorystore, ...) before this resource.
resource "google_service_networking_connection" "this" {
  network                 = var.spec.network
  service                 = local.service
  reserved_peering_ranges = var.spec.reserved_peering_ranges

  # Adopts a pre-existing connection for the same pair instead of failing
  # with "Cannot modify allocated ranges" (see the spec comment).
  update_on_creation_fail = var.spec.update_on_creation_fail

  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [
    google_project_service.servicenetworking_api,
    google_project_service.compute_api,
  ]
}
