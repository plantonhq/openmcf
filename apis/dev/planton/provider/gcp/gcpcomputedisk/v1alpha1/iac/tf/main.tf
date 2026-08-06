# Enable the Compute Engine API — the control plane that owns disks.
# disable_on_destroy is false: tearing down one disk must never disable
# the API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The persistent disk. Sharp edges, all taught by the API rather than
# invented here:
#
#   - name, zone, type, sources, encryption, and architecture are
#     ForceNew — changing them replaces the disk and its data. size grows
#     in place but never shrinks.
#
#   - At most one source (image / snapshot / source_disk); none creates
#     an empty disk — the common case for data volumes.
#
#   - Deleting a disk still attached to a running instance fails; detach
#     first (or delete the instance).
#
#   - create_snapshot_before_destroy takes a final snapshot during
#     destroy — a last-resort recovery net for precious volumes (CMEK
#     disks reuse their key for the snapshot).
resource "google_compute_disk" "this" {
  name    = local.disk_name
  project = local.project_id
  zone    = var.spec.zone
  labels  = local.final_labels

  description = var.spec.description != "" ? var.spec.description : null
  type        = local.type
  size        = var.spec.size_gb

  image        = local.image
  snapshot     = local.snapshot
  source_disk  = local.source_disk
  access_mode  = local.access_mode
  architecture = local.architecture
  storage_pool = local.storage_pool

  provisioned_iops       = var.spec.provisioned_iops
  provisioned_throughput = var.spec.provisioned_throughput

  enable_confidential_compute = var.spec.enable_confidential_compute ? true : null
  physical_block_size_bytes   = var.spec.physical_block_size_bytes

  create_snapshot_before_destroy        = var.spec.create_snapshot_before_destroy
  create_snapshot_before_destroy_prefix = var.spec.snapshot_before_destroy_prefix != "" ? var.spec.snapshot_before_destroy_prefix : null

  # CMEK: the Compute Engine service agent must hold
  # roles/cloudkms.cryptoKeyEncrypterDecrypter on this key before create.
  dynamic "disk_encryption_key" {
    for_each = var.spec.kms_key != "" ? [var.spec.kms_key] : []
    content {
      kms_key_self_link = disk_encryption_key.value
    }
  }

  # Resource Manager tags bind at create time only.
  dynamic "params" {
    for_each = length(var.spec.resource_manager_tags) > 0 ? [var.spec.resource_manager_tags] : []
    content {
      resource_manager_tags = params.value
    }
  }

  depends_on = [
    google_project_service.compute_api,
  ]
}
