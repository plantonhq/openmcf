# Enable the Filestore API — the control plane that owns instances.
# disable_on_destroy is false: tearing down one instance must never
# disable the API for everything else in the project.
resource "google_project_service" "file_api" {
  project = local.project_id
  service = "file.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Filestore instance. Sharp edges, all taught by the API rather than
# invented here:
#
#   - name, location, tier, protocol, network attachment, KMS key, and
#     replication are immutable — changing any of them replaces the
#     instance (and its data). file share capacity grows in place but
#     never shrinks.
#
#   - deletion_protection_enabled must be flipped false before a
#     protected instance can be destroyed.
#
#   - connect_mode defaults to DIRECT_PEERING; PRIVATE_SERVICE_ACCESS is
#     required for Shared VPC consumers and rides an existing
#     service-networking connection on the VPC.
resource "google_filestore_instance" "this" {
  name     = local.instance_name
  project  = local.project_id
  location = local.location
  tier     = local.tier

  description  = local.description
  protocol     = local.protocol
  kms_key_name = local.kms_key_name

  deletion_protection_enabled = var.spec.deletion_protection_enabled
  deletion_protection_reason  = var.spec.deletion_protection_reason != "" ? var.spec.deletion_protection_reason : null

  labels = local.final_labels
  tags   = length(var.spec.tags) > 0 ? var.spec.tags : null

  file_shares {
    name        = var.spec.file_share.name
    capacity_gb = var.spec.file_share.capacity_gb

    # Restore-from-backup is create-time only; the share's capacity must
    # cover the backup's source capacity.
    source_backup = var.spec.file_share.source_backup != "" ? var.spec.file_share.source_backup : null

    dynamic "nfs_export_options" {
      for_each = var.spec.file_share.nfs_export_options
      content {
        ip_ranges   = length(nfs_export_options.value.ip_ranges) > 0 ? nfs_export_options.value.ip_ranges : null
        access_mode = nfs_export_options.value.access_mode != "" ? nfs_export_options.value.access_mode : null
        squash_mode = nfs_export_options.value.squash_mode != "" ? nfs_export_options.value.squash_mode : null
        anon_uid    = nfs_export_options.value.anon_uid
        anon_gid    = nfs_export_options.value.anon_gid
      }
    }
  }

  networks {
    network           = local.network
    modes             = local.modes
    connect_mode      = local.connect_mode
    reserved_ip_range = local.reserved_ip_range
  }

  dynamic "performance_config" {
    for_each = var.spec.performance_config != null ? [var.spec.performance_config] : []
    content {
      dynamic "fixed_iops" {
        for_each = performance_config.value.fixed_iops != null ? [performance_config.value.fixed_iops] : []
        content {
          max_iops = fixed_iops.value.max_iops
        }
      }
      dynamic "iops_per_tb" {
        for_each = performance_config.value.iops_per_tb != null ? [performance_config.value.iops_per_tb] : []
        content {
          max_iops_per_tb = iops_per_tb.value.max_iops_per_tb
        }
      }
    }
  }

  # Create-time replication: this instance joins as ACTIVE source or
  # STANDBY replica of the referenced peers. Backups cannot be taken from
  # a STANDBY replica.
  dynamic "initial_replication" {
    for_each = var.spec.initial_replication != null ? [var.spec.initial_replication] : []
    content {
      role = initial_replication.value.role != "" ? initial_replication.value.role : null

      dynamic "replicas" {
        for_each = initial_replication.value.peer_instances
        content {
          peer_instance = replicas.value
        }
      }
    }
  }

  depends_on = [
    google_project_service.file_api,
  ]
}
