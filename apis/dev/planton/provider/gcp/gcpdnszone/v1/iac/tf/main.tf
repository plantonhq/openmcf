# Enable the Cloud DNS API so a fresh project can host managed zones.
# disable_on_destroy is false: tearing down one zone must never disable
# the API for everything else in the project.
resource "google_project_service" "dns_api" {
  project = local.project_id
  service = "dns.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Cloud DNS managed zone — the zone shell only. Individual DNS records
# belong in the separate GcpDnsRecord kind.
resource "google_dns_managed_zone" "managed_zone" {
  name        = local.managed_zone_name
  project     = local.project_id
  description = local.description
  dns_name    = local.zone_dns_name
  visibility  = local.visibility
  labels      = length(local.labels) > 0 ? local.labels : null

  force_destroy = var.spec.force_destroy

  depends_on = [google_project_service.dns_api]

  dynamic "dnssec_config" {
    for_each = var.spec.dnssec_config != null ? [var.spec.dnssec_config] : []
    content {
      state         = dnssec_config.value.state != "" ? dnssec_config.value.state : null
      non_existence = dnssec_config.value.non_existence != "" ? dnssec_config.value.non_existence : null

      dynamic "default_key_specs" {
        for_each = dnssec_config.value.default_key_specs
        content {
          algorithm  = default_key_specs.value.algorithm != "" ? default_key_specs.value.algorithm : null
          key_length = default_key_specs.value.key_length > 0 ? default_key_specs.value.key_length : null
          key_type   = default_key_specs.value.key_type != "" ? default_key_specs.value.key_type : null
        }
      }
    }
  }

  dynamic "private_visibility_config" {
    for_each = var.spec.private_visibility_config != null ? [var.spec.private_visibility_config] : []
    content {
      dynamic "networks" {
        for_each = private_visibility_config.value.networks
        content {
          network_url = networks.value.network_url
        }
      }

      dynamic "gke_clusters" {
        for_each = private_visibility_config.value.gke_clusters
        content {
          gke_cluster_name = gke_clusters.value.gke_cluster_name
        }
      }
    }
  }

  dynamic "forwarding_config" {
    # try() guards the null object: HCL's && does not short-circuit.
    for_each = (
      length(try(var.spec.forwarding_config.target_name_servers, [])) > 0
    ) ? [var.spec.forwarding_config] : []
    content {
      dynamic "target_name_servers" {
        for_each = forwarding_config.value.target_name_servers
        content {
          ipv4_address    = target_name_servers.value.ipv4_address != "" ? target_name_servers.value.ipv4_address : null
          domain_name     = target_name_servers.value.domain_name != "" ? target_name_servers.value.domain_name : null
          forwarding_path = target_name_servers.value.forwarding_path != "" ? target_name_servers.value.forwarding_path : null
        }
      }
    }
  }

  dynamic "peering_config" {
    for_each = var.spec.peering_config != null ? [var.spec.peering_config] : []
    content {
      target_network {
        network_url = peering_config.value.target_network
      }
    }
  }

  dynamic "cloud_logging_config" {
    for_each = var.spec.cloud_logging_config != null ? [var.spec.cloud_logging_config] : []
    content {
      enable_logging = cloud_logging_config.value.enable_logging
    }
  }
}
