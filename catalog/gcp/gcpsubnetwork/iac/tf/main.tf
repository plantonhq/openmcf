# Enable the Compute Engine API so a fresh project can host the subnetwork.
# disable_on_destroy is false: tearing down one subnet must never disable the
# API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A subnetwork in a custom-mode VPC — the regional address space workloads
# live in: primary IPv4 range for VM interfaces, secondary ranges for alias
# IPs (GKE pods/services), optional IPv6, and VPC Flow Logs.
#
# name, project, region, network, and description are immutable (ForceNew):
# changing any of them destroys and recreates the subnet — an outage for
# everything addressed in it. The primary range is the one asymmetric knob:
# EXPANDING ip_cidr_range updates in place, shrinking recreates (the provider
# guards this with a shrink-only ForceNew diff).
#
# purpose creates the special-role subnets other features depend on:
# REGIONAL_MANAGED_PROXY reserves Envoy address space for the region's
# regional load balancers (required before a regional ALB exists in the VPC);
# PRIVATE_SERVICE_CONNECT backs published PSC services.
resource "google_compute_subnetwork" "main" {
  # google-beta: allow_subnet_cidr_routes_overlap is beta-only on the released
  # 6.x line (everything else here is GA and identical in beta). The Pulumi
  # provider is beta-bridged by construction, so both engines expose the same
  # surface.
  provider = google-beta

  name          = var.spec.subnetwork_name
  project       = local.project_id
  region        = var.spec.region
  network       = var.spec.vpc_self_link
  description   = var.spec.description
  ip_cidr_range = local.ip_cidr_range

  purpose = local.purpose
  role    = local.role

  private_ip_google_access   = local.private_ip_google_access
  private_ipv6_google_access = local.private_ipv6_google_access

  # Dual-stack wiring: stack_type opts into IPv6; ipv6_access_type decides
  # whether the assigned prefix is internet-routable (EXTERNAL GUAs) or
  # VPC-internal (INTERNAL ULAs — requires the VPC's ULA range enabled).
  stack_type           = local.stack_type
  ipv6_access_type     = local.ipv6_access_type
  external_ipv6_prefix = local.external_ipv6_prefix

  # Deliberate address-space reclaims only: subnet routes still win over the
  # overlapping peer/on-prem routes this permits.
  allow_subnet_cidr_routes_overlap = coalesce(var.spec.allow_subnet_cidr_routes_overlap, false)

  # Safety latch: by default an empty secondary-range list is NOT sent on
  # update, so a partial manifest cannot silently wipe GKE pod ranges.
  send_secondary_ip_range_if_empty = coalesce(var.spec.send_secondary_ip_range_if_empty, false)

  # Secondary (alias) ranges — the mechanism GKE uses for pod/service IPs.
  # Consumers select a range by its name (e.g. ip_allocation_policy).
  dynamic "secondary_ip_range" {
    for_each = local.secondary_ip_ranges
    content {
      range_name    = secondary_ip_range.value.range_name
      ip_cidr_range = secondary_ip_range.value.ip_cidr_range
    }
  }

  # VPC Flow Logs: presence of the block enables logging. Defaults mirror
  # the API's own (5s aggregation, 50% sampling, all metadata) so an empty
  # spec object behaves sanely; metadata_fields only accompanies
  # CUSTOM_METADATA (spec-enforced).
  dynamic "log_config" {
    for_each = local.log_config != null ? [local.log_config] : []
    content {
      aggregation_interval = log_config.value.aggregation_interval
      flow_sampling        = log_config.value.flow_sampling
      metadata             = log_config.value.metadata
      metadata_fields      = log_config.value.metadata == "CUSTOM_METADATA" ? log_config.value.metadata_fields : null
      filter_expr          = log_config.value.filter_expr
    }
  }

  # Depend on API enablement so a fresh project works first try.
  depends_on = [google_project_service.compute_api]
}
