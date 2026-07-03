#############################################
# Enable Compute API
#############################################

resource "google_project_service" "compute_api" {
  project                    = local.project_id
  service                    = "compute.googleapis.com"
  disable_dependent_services = true
  disable_on_destroy         = false
}

#############################################
# GCP VPC Network
#############################################

resource "google_compute_network" "vpc" {
  name                    = var.spec.network_name
  project                 = local.project_id
  auto_create_subnetworks = var.spec.auto_create_subnetworks
  routing_mode            = local.routing_mode

  description = var.spec.description != "" ? var.spec.description : null

  mtu = var.spec.mtu != null ? var.spec.mtu : null

  enable_ula_internal_ipv6 = var.spec.enable_ula_internal_ipv6 ? true : null
  internal_ipv6_range      = var.spec.internal_ipv6_range != "" ? var.spec.internal_ipv6_range : null

  network_firewall_policy_enforcement_order = var.spec.network_firewall_policy_enforcement_order != "" ? var.spec.network_firewall_policy_enforcement_order : null

  network_profile = var.spec.network_profile != "" ? var.spec.network_profile : null

  delete_default_routes_on_create = var.spec.delete_default_routes_on_create ? true : null

  bgp_best_path_selection_mode = var.spec.bgp_best_path_selection.mode != "" ? var.spec.bgp_best_path_selection.mode : null
  bgp_always_compare_med       = var.spec.bgp_best_path_selection.always_compare_med ? true : null
  bgp_inter_region_cost        = var.spec.bgp_best_path_selection.inter_region_cost != "" ? var.spec.bgp_best_path_selection.inter_region_cost : null

  depends_on = [google_project_service.compute_api]
}
