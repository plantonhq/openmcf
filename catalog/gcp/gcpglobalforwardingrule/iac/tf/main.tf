# Enable the Compute Engine API so a fresh project can host the rule.
# disable_on_destroy is false: tearing down one forwarding rule must never
# disable the API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A global Compute Engine forwarding rule — the VIP node where traffic
# enters a global load balancer (or, with the PSC form, where a VPC's
# private path to Google APIs / a producer service begins). It binds an IP
# address and port to a target proxy; everything behind it (proxy → URL map
# → backend service → backends) is wiring.
#
# target and labels update in place (setTarget is the zero-downtime frontend
# swap); every other field is immutable (ForceNew). The VIP itself survives
# recreation only when ip_address references a reserved static address —
# which is why production frontends reserve one.
resource "google_compute_global_forwarding_rule" "this" {
  name        = local.forwarding_rule_name
  project     = local.project_id
  description = local.description

  # Arrives resolved to a literal: a proxy self-link, a PSC bundle name
  # (all-apis / vpc-sc), or a service attachment URI.
  target = var.spec.target

  # Null → GCP assigns an ephemeral IP. A GcpGlobalAddress ref resolves to
  # its literal IP (the API reads back the IP number, so passing the number
  # keeps state drift-free).
  ip_address = local.ip_address

  ip_protocol = local.ip_protocol
  ip_version  = local.ip_version

  # The spec's NONE sentinel maps to the API's empty scheme (Private
  # Service Connect); null lets the provider apply its default (EXTERNAL).
  load_balancing_scheme = local.load_balancing_scheme

  port_range = local.port_range

  network    = local.network
  subnetwork = local.subnetwork

  # Global rules are PREMIUM-only (spec CEL enforces it).
  network_tier = local.network_tier

  # Traffic Director xDS scoping (INTERNAL_SELF_MANAGED only).
  dynamic "metadata_filters" {
    for_each = var.spec.metadata_filters
    content {
      filter_match_criteria = metadata_filters.value.filter_match_criteria

      dynamic "filter_labels" {
        for_each = metadata_filters.value.filter_labels
        content {
          name  = filter_labels.value.name
          value = filter_labels.value.value
        }
      }
    }
  }

  # Service Directory registration for PSC-for-Google-APIs frontends.
  dynamic "service_directory_registrations" {
    for_each = var.spec.service_directory_registration != null ? [var.spec.service_directory_registration] : []
    content {
      namespace                = service_directory_registrations.value.namespace != "" ? service_directory_registrations.value.namespace : null
      service_directory_region = service_directory_registrations.value.service_directory_region != "" ? service_directory_registrations.value.service_directory_region : null
    }
  }

  # Only meaningful for PSC; null keeps the API default (auto-create the
  # PSC DNS zone).
  no_automate_dns_zone = local.no_automate_dns_zone

  labels = local.labels

  # EXTERNAL → EXTERNAL_MANAGED backend-bucket canary migration.
  external_managed_backend_bucket_migration_state              = local.migration_state
  external_managed_backend_bucket_migration_testing_percentage = local.migration_testing_percentage

  # What destroy does to the frontend: DELETE (default), PREVENT (refuse),
  # or ABANDON (drop from state, keep serving traffic).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.compute_api]
}
