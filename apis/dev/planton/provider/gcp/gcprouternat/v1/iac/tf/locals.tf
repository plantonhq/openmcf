locals {
  # Honor the spec contract: an empty project_id falls back to the provider's
  # default project. Passing null (instead of "") lets the google provider
  # resolve its own project from configuration or the GOOGLE_PROJECT /
  # GOOGLE_CLOUD_PROJECT environment chain; an empty string would be sent
  # verbatim and rejected by the API.
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  # NAT IP allocation is derived, not declared: referencing reservations
  # selects MANUAL_ONLY, an empty list selects AUTO_ONLY. The two cannot be
  # combined, so a separate spec field would only create contradictions.
  nat_ip_allocate_option = length(var.spec.nat_ips) > 0 ? "MANUAL_ONLY" : "AUTO_ONLY"

  # Subnetwork scoping mode: an explicit value wins; otherwise listing
  # subnetworks implies LIST_OF_SUBNETWORKS and an empty list means
  # everything in the region (primary + secondary ranges).
  source_subnetwork_ip_ranges_to_nat = (
    var.spec.source_subnetwork_ip_ranges_to_nat != ""
    ? var.spec.source_subnetwork_ip_ranges_to_nat
    : (length(var.spec.subnetworks) > 0 ? "LIST_OF_SUBNETWORKS" : "ALL_SUBNETWORKS_ALL_IP_RANGES")
  )

  # Empty optional values become null so the provider omits them from the
  # API payload and GCP's own defaults apply (documented per field in the
  # spec) instead of the module hardcoding its own.
  type              = var.spec.type != "" ? var.spec.type : null
  auto_network_tier = var.spec.auto_network_tier != "" ? var.spec.auto_network_tier : null

  min_ports_per_vm = var.spec.min_ports_per_vm > 0 ? var.spec.min_ports_per_vm : null
  max_ports_per_vm = var.spec.max_ports_per_vm > 0 ? var.spec.max_ports_per_vm : null

  udp_idle_timeout_sec             = var.spec.udp_idle_timeout_sec > 0 ? var.spec.udp_idle_timeout_sec : null
  icmp_idle_timeout_sec            = var.spec.icmp_idle_timeout_sec > 0 ? var.spec.icmp_idle_timeout_sec : null
  tcp_established_idle_timeout_sec = var.spec.tcp_established_idle_timeout_sec > 0 ? var.spec.tcp_established_idle_timeout_sec : null
  tcp_transitory_idle_timeout_sec  = var.spec.tcp_transitory_idle_timeout_sec > 0 ? var.spec.tcp_transitory_idle_timeout_sec : null
  tcp_time_wait_timeout_sec        = var.spec.tcp_time_wait_timeout_sec > 0 ? var.spec.tcp_time_wait_timeout_sec : null

  router_asn                = var.spec.router_asn > 0 ? var.spec.router_asn : null
  router_keepalive_interval = var.spec.router_keepalive_interval > 0 ? var.spec.router_keepalive_interval : null

  # Logging: the DISABLED sentinel turns logging off; every other filter
  # value enables it. The filter must still be a valid value when disabled,
  # so ERRORS_ONLY is sent as a placeholder the API ignores.
  enable_logging = var.spec.log_filter != "DISABLED"
  log_filter     = var.spec.log_filter != "DISABLED" ? var.spec.log_filter : "ERRORS_ONLY"
}
