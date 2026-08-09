# Enable the Cloud DNS API so a fresh project can host record sets.
# disable_on_destroy is false: tearing down one record must never disable
# the API for everything else in the project.
resource "google_project_service" "dns_api" {
  project = local.project_id
  service = "dns.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# One DNS record set: a (name, type) pair answered either with static
# rrdatas or with exactly one routing policy — never both (the provider
# enforces ExactlyOneOf; the spec enforces the same rule pre-deploy).
resource "google_dns_record_set" "record" {
  project      = local.project_id
  managed_zone = var.spec.managed_zone
  name         = var.spec.name
  type         = var.spec.type
  ttl          = local.ttl_seconds

  # Static values arm. An empty list is passed as null so the provider's
  # ExactlyOneOf sees the attribute as absent when routing_policy is used.
  rrdatas = length(var.spec.values) > 0 ? var.spec.values : null

  # What destroy does to the record set: DELETE (default), PREVENT
  # (refuse), or ABANDON (drop from state, keep answering queries).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  depends_on = [google_project_service.dns_api]

  dynamic "routing_policy" {
    for_each = local.routing_policy != null ? [local.routing_policy] : []
    content {
      # Geo fencing applies only to geolocation routing; the spec rejects
      # it for other styles, so passing it through unconditionally is safe.
      enable_geo_fencing = routing_policy.value.enable_geo_fencing
      health_check       = local.routing_health_check

      # Weighted round robin: traffic splits by weight ratio.
      dynamic "wrr" {
        for_each = routing_policy.value.wrr
        content {
          weight  = wrr.value.weight
          rrdatas = length(wrr.value.values) > 0 ? wrr.value.values : null

          dynamic "health_checked_targets" {
            for_each = wrr.value.health_checked_targets != null ? [wrr.value.health_checked_targets] : []
            content {
              external_endpoints = length(health_checked_targets.value.external_endpoints) > 0 ? health_checked_targets.value.external_endpoints : null

              dynamic "internal_load_balancers" {
                for_each = health_checked_targets.value.internal_load_balancers
                content {
                  ip_address         = internal_load_balancers.value.ip_address
                  ip_protocol        = internal_load_balancers.value.ip_protocol
                  load_balancer_type = internal_load_balancers.value.load_balancer_type != "" ? internal_load_balancers.value.load_balancer_type : null
                  network_url        = internal_load_balancers.value.network_url
                  port               = internal_load_balancers.value.port
                  project            = internal_load_balancers.value.project
                  region             = internal_load_balancers.value.region != "" ? internal_load_balancers.value.region : null
                }
              }
            }
          }
        }
      }

      # Geolocation routing: nearest-location answering.
      dynamic "geo" {
        for_each = routing_policy.value.geo
        content {
          location = geo.value.location
          rrdatas  = length(geo.value.values) > 0 ? geo.value.values : null

          dynamic "health_checked_targets" {
            for_each = geo.value.health_checked_targets != null ? [geo.value.health_checked_targets] : []
            content {
              external_endpoints = length(health_checked_targets.value.external_endpoints) > 0 ? health_checked_targets.value.external_endpoints : null

              dynamic "internal_load_balancers" {
                for_each = health_checked_targets.value.internal_load_balancers
                content {
                  ip_address         = internal_load_balancers.value.ip_address
                  ip_protocol        = internal_load_balancers.value.ip_protocol
                  load_balancer_type = internal_load_balancers.value.load_balancer_type != "" ? internal_load_balancers.value.load_balancer_type : null
                  network_url        = internal_load_balancers.value.network_url
                  port               = internal_load_balancers.value.port
                  project            = internal_load_balancers.value.project
                  region             = internal_load_balancers.value.region != "" ? internal_load_balancers.value.region : null
                }
              }
            }
          }
        }
      }

      # Failover: global primaries answered while healthy, then the
      # regional backup_geo policy takes over.
      dynamic "primary_backup" {
        for_each = routing_policy.value.primary_backup != null ? [routing_policy.value.primary_backup] : []
        content {
          trickle_ratio                  = primary_backup.value.trickle_ratio
          enable_geo_fencing_for_backups = primary_backup.value.enable_geo_fencing_for_backups

          primary {
            external_endpoints = length(primary_backup.value.primary.external_endpoints) > 0 ? primary_backup.value.primary.external_endpoints : null

            dynamic "internal_load_balancers" {
              for_each = primary_backup.value.primary.internal_load_balancers
              content {
                ip_address         = internal_load_balancers.value.ip_address
                ip_protocol        = internal_load_balancers.value.ip_protocol
                load_balancer_type = internal_load_balancers.value.load_balancer_type != "" ? internal_load_balancers.value.load_balancer_type : null
                network_url        = internal_load_balancers.value.network_url
                port               = internal_load_balancers.value.port
                project            = internal_load_balancers.value.project
                region             = internal_load_balancers.value.region != "" ? internal_load_balancers.value.region : null
              }
            }
          }

          dynamic "backup_geo" {
            for_each = primary_backup.value.backup_geo
            content {
              location = backup_geo.value.location
              rrdatas  = length(backup_geo.value.values) > 0 ? backup_geo.value.values : null

              dynamic "health_checked_targets" {
                for_each = backup_geo.value.health_checked_targets != null ? [backup_geo.value.health_checked_targets] : []
                content {
                  external_endpoints = length(health_checked_targets.value.external_endpoints) > 0 ? health_checked_targets.value.external_endpoints : null

                  dynamic "internal_load_balancers" {
                    for_each = health_checked_targets.value.internal_load_balancers
                    content {
                      ip_address         = internal_load_balancers.value.ip_address
                      ip_protocol        = internal_load_balancers.value.ip_protocol
                      load_balancer_type = internal_load_balancers.value.load_balancer_type != "" ? internal_load_balancers.value.load_balancer_type : null
                      network_url        = internal_load_balancers.value.network_url
                      port               = internal_load_balancers.value.port
                      project            = internal_load_balancers.value.project
                      region             = internal_load_balancers.value.region != "" ? internal_load_balancers.value.region : null
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
