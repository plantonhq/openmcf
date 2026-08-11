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
# BGP surface at all — router_bgp matters only when the router will also
# terminate BGP sessions (Interconnect/VPN), where it controls the ASN and
# what routes are advertised to every peer.
resource "google_compute_router" "router" {
  name        = var.spec.router_name
  region      = var.spec.region
  network     = var.spec.vpc_self_link
  project     = local.project_id
  description = var.spec.router_description != "" ? var.spec.router_description : null

  # Dedicates the router to encrypted VLAN attachments (HA VPN over
  # Interconnect). Immutable, and an encrypted router cannot be converted.
  encrypted_interconnect_router = var.spec.encrypted_interconnect_router

  # Message presence IS the block: the spec's router_bgp carries a valid
  # ASN whenever it is set (CEL-enforced), so no derivation is needed.
  dynamic "bgp" {
    for_each = var.spec.router_bgp != null ? [var.spec.router_bgp] : []
    content {
      asn = bgp.value.asn
      # advertise_mode empty means DEFAULT; custom groups/ranges are legal
      # only in CUSTOM mode (spec CELs mirror the provider's constraint).
      advertise_mode     = bgp.value.advertise_mode != "" ? bgp.value.advertise_mode : null
      advertised_groups  = length(bgp.value.advertised_groups) > 0 ? bgp.value.advertised_groups : null
      keepalive_interval = bgp.value.keepalive_interval > 0 ? bgp.value.keepalive_interval : null
      # identifier_range is Optional+Computed: an empty string must stay out
      # of the payload or it would fight the API's computed value.
      identifier_range = bgp.value.identifier_range != "" ? bgp.value.identifier_range : null

      dynamic "advertised_ip_ranges" {
        for_each = bgp.value.advertised_ip_ranges
        content {
          range       = advertised_ip_ranges.value.range
          description = advertised_ip_ranges.value.description != "" ? advertised_ip_ranges.value.description : null
        }
      }
    }
  }

  # Create-time resource-manager tags (org policy / IAM conditions).
  dynamic "params" {
    for_each = length(var.spec.resource_manager_tags) > 0 ? [1] : []
    content {
      resource_manager_tags = var.spec.resource_manager_tags
    }
  }

  deletion_policy = local.deletion_policy

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

  # NAT64 (IPv6-to-IPv4 translation). Null means no NAT64; only one NAT
  # per region in the network may claim ALL_IPV6_SUBNETWORKS.
  source_subnetwork_ip_ranges_to_nat64 = local.source_subnetwork_ip_ranges_to_nat64

  dynamic "nat64_subnetwork" {
    for_each = var.spec.nat64_subnetworks
    content {
      name = nat64_subnetwork.value
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

  deletion_policy = local.deletion_policy
}
