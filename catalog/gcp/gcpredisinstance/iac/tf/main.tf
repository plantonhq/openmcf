# Enable the Memorystore for Redis API so a fresh project can host the
# instance. disable_on_destroy is false: tearing down one instance must never
# disable the API for everything else in the project.
resource "google_project_service" "redis_api" {
  project = local.project_id
  service = "redis.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A Memorystore for Redis instance — the classic VPC-peered managed Redis.
# One resource is one instance: a standalone BASIC node, or a STANDARD_HA
# primary + failover replica (optionally with read replicas).
#
# Lifecycle notes the API enforces:
#   - name, tier, connect_mode, transit_encryption_mode, authorized_network,
#     reserved_ip_range, location_id, alternative_location_id, and the CMEK
#     key are immutable — changing any of them replaces the instance.
#   - memory_size_gb resizes in place; redis_version upgrades in place but a
#     downgrade replaces the instance.
#   - read_replicas_mode is set at creation; replica_count and
#     secondary_ip_range are the in-place scale-out levers afterwards.
#   - with connect_mode PRIVATE_SERVICE_ACCESS, the network must already
#     carry a service networking connection — the API rejects the create
#     otherwise.
resource "google_redis_instance" "this" {
  name           = var.spec.instance_name
  project        = local.project_id
  region         = var.spec.region
  tier           = var.spec.tier
  memory_size_gb = var.spec.memory_size_gb

  redis_version = local.redis_version
  display_name  = local.display_name

  # Zone placement: primary zone, and (STANDARD_HA only) the replica zone.
  # Left unset, GCP spreads the nodes across zones automatically.
  location_id             = local.location_id
  alternative_location_id = local.alternative_location_id

  # Connectivity is fixed at creation: the peered network, how it peers
  # (direct peering vs the network's private services access connection),
  # and which internal range the nodes occupy.
  authorized_network = local.authorized_network
  connect_mode       = local.connect_mode
  reserved_ip_range  = local.reserved_ip_range

  # The scale-out range: adding read replicas to an existing instance needs
  # more address space than the original /29 provides.
  secondary_ip_range = local.secondary_ip_range

  auth_enabled            = var.spec.auth_enabled
  transit_encryption_mode = local.transit_encryption_mode

  redis_configs = length(var.spec.redis_configs) > 0 ? var.spec.redis_configs : null

  # Self-service maintenance: setting a newer available version applies the
  # update now instead of waiting for GCP's rollout.
  maintenance_version = local.maintenance_version

  read_replicas_mode = local.read_replicas_mode
  replica_count      = local.replica_count

  customer_managed_key = local.customer_managed_key

  labels = local.final_labels

  # Destroy guard, sent explicitly from the spec (default true) so destroy
  # behavior is identical on both engines: omitting it would leave the
  # decision to whatever happens to be in Terraform state.
  deletion_protection = var.spec.deletion_protection

  # Client-side destroy behavior: DELETE (default), PREVENT, or ABANDON.
  # Sent only when set so the provider default stays in charge otherwise.
  # Evaluated only after deletion_protection allows the destroy at all.
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  dynamic "maintenance_policy" {
    for_each = var.spec.maintenance_window != null ? [var.spec.maintenance_window] : []
    content {
      description = maintenance_policy.value.description != "" ? maintenance_policy.value.description : null
      weekly_maintenance_window {
        day = maintenance_policy.value.day
        start_time {
          hours   = maintenance_policy.value.hour
          minutes = maintenance_policy.value.minute
        }
      }
    }
  }

  dynamic "persistence_config" {
    for_each = var.spec.persistence_config != null ? [var.spec.persistence_config] : []
    content {
      persistence_mode        = persistence_config.value.persistence_mode
      rdb_snapshot_period     = persistence_config.value.rdb_snapshot_period != "" ? persistence_config.value.rdb_snapshot_period : null
      rdb_snapshot_start_time = persistence_config.value.rdb_snapshot_start_time != "" ? persistence_config.value.rdb_snapshot_start_time : null
    }
  }

  depends_on = [google_project_service.redis_api]
}
