# Enable the Compute Engine API so a fresh project can host the router and
# NAT. disable_on_destroy is false: tearing down one gateway must never
# disable the API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Cloud Router carrying the NAT configuration. NAT-only routers need no
# BGP surface at all — the asn/keepalive knobs matter only when the router
# will also terminate BGP sessions (Interconnect/VPN).
resource "google_compute_router" "router" {
  name    = var.spec.router_name
  region  = var.spec.region
  network = var.spec.vpc_self_link
  project = local.project_id

  dynamic "bgp" {
    for_each = local.router_asn != null || local.router_keepalive_interval != null ? [1] : []
    content {
      asn                = local.router_asn
      keepalive_interval = local.router_keepalive_interval
    }
  }

  depends_on = [google_project_service.compute_api]
}

# The NAT gateway. Manual NAT IPs are REFERENCED GcpAddress reservations,
# never created here — the reservation is its own composable node with its
# own lifecycle, and its literal IP is read from that node's outputs.
#
# Everything below the allocation block updates in place: subnetwork
# scoping, port tuning, timeouts, rules, and logging are all zero-downtime
# levers on a live gateway.
resource "google_compute_router_nat" "nat" {
  name    = var.spec.nat_name
  router  = google_compute_router.router.name
  region  = var.spec.region
  project = local.project_id

  # PUBLIC (internet egress) or PRIVATE (NCC spoke-to-spoke); omitted means
  # PUBLIC. Private NAT draws addresses from subnetwork ranges, so the
  # manual/auto IP machinery below stays empty for it (spec CEL enforces).
  type = local.type

  # Allocation is derived from the nat_ips list (see locals.tf). Drained
  # IPs keep serving established connections but accept no new ones — the
  # zero-downtime path for rotating an egress IP out of service.
  # Empty lists become null: nat_ips/drain_nat_ips are Optional+Computed in
  # the provider, and an explicitly-set empty collection would fight the
  # API's computed value on every plan.
  nat_ip_allocate_option = local.type == "PRIVATE" ? null : local.nat_ip_allocate_option
  nat_ips                = length(var.spec.nat_ips) > 0 ? var.spec.nat_ips : null
  drain_nat_ips          = length(var.spec.drain_nat_ips) > 0 ? var.spec.drain_nat_ips : null
  auto_network_tier      = local.auto_network_tier

  # Which subnetworks (and which of their ranges) route through this NAT.
  source_subnetwork_ip_ranges_to_nat = local.source_subnetwork_ip_ranges_to_nat

  dynamic "subnetwork" {
    for_each = var.spec.subnetworks
    content {
      name = subnetwork.value.subnetwork
      # An empty list means everything: primary + all secondary ranges.
      source_ip_ranges_to_nat  = length(subnetwork.value.source_ip_ranges_to_nat) > 0 ? subnetwork.value.source_ip_ranges_to_nat : ["ALL_IP_RANGES"]
      secondary_ip_range_names = length(subnetwork.value.secondary_ip_range_names) > 0 ? subnetwork.value.secondary_ip_range_names : null
    }
  }

  # Port allocation: nulls defer to GCP's defaults (64 static / 32 dynamic
  # min ports). Dynamic allocation grows a busy VM's ports toward the max;
  # it cannot coexist with endpoint-independent mapping (spec CEL enforces
  # the conflict pre-deploy, matching the API).
  min_ports_per_vm                    = local.min_ports_per_vm
  max_ports_per_vm                    = local.max_ports_per_vm
  enable_dynamic_port_allocation      = var.spec.enable_dynamic_port_allocation
  enable_endpoint_independent_mapping = var.spec.enable_endpoint_independent_mapping

  # Which resource class this NAT serves (VM instances by default).
  endpoint_types = length(var.spec.endpoint_types) > 0 ? var.spec.endpoint_types : null

  # Connection timeouts: nulls defer to GCP's defaults (30s UDP/ICMP/
  # transitory-TCP, 1200s established-TCP, 120s TIME_WAIT).
  udp_idle_timeout_sec             = local.udp_idle_timeout_sec
  icmp_idle_timeout_sec            = local.icmp_idle_timeout_sec
  tcp_established_idle_timeout_sec = local.tcp_established_idle_timeout_sec
  tcp_transitory_idle_timeout_sec  = local.tcp_transitory_idle_timeout_sec
  tcp_time_wait_timeout_sec        = local.tcp_time_wait_timeout_sec

  # NAT rules: dedicated NAT IPs/ranges for traffic matching a CEL
  # expression — e.g. a stable, separately allowlistable source IP for one
  # partner's endpoints. Lower rule numbers win.
  dynamic "rules" {
    for_each = var.spec.rules
    content {
      rule_number = rules.value.rule_number
      match       = rules.value.match
      description = rules.value.description != "" ? rules.value.description : null

      dynamic "action" {
        for_each = rules.value.action != null ? [rules.value.action] : []
        content {
          source_nat_active_ips    = length(action.value.source_nat_active_ips) > 0 ? action.value.source_nat_active_ips : null
          source_nat_drain_ips     = length(action.value.source_nat_drain_ips) > 0 ? action.value.source_nat_drain_ips : null
          source_nat_active_ranges = length(action.value.source_nat_active_ranges) > 0 ? action.value.source_nat_active_ranges : null
          source_nat_drain_ranges  = length(action.value.source_nat_drain_ranges) > 0 ? action.value.source_nat_drain_ranges : null
        }
      }
    }
  }

  log_config {
    enable = local.enable_logging
    filter = local.log_filter
  }
}
