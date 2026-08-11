# Enable the Bigtable Admin API — the control plane instance and table
# management run through. disable_on_destroy is false: tearing down one
# instance must never disable the API for everything else in the project.
resource "google_project_service" "bigtableadmin_api" {
  project = local.project_id
  service = "bigtableadmin.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Bigtable instance: the logical container, served by one or more
# clusters (physical replicas, each in its own zone). Multi-cluster
# instances replicate automatically; the client library routes and fails
# over transparently.
#
# Per-cluster immutability is enforced server-side: zone, storage_type,
# kms_key_name, and node_scaling_factor cannot change on an existing
# cluster_id. num_nodes and autoscaling bounds resize in place.
# deletion_protection (spec default TRUE) is a Terraform-side guard —
# destroy fails until it is explicitly set false; force_destroy clears
# backups that would otherwise block deletion.
resource "google_bigtable_instance" "this" {
  name                = var.spec.instance_name
  project             = local.project_id
  labels              = local.final_labels
  deletion_protection = var.spec.deletion_protection
  force_destroy       = var.spec.force_destroy

  display_name = var.spec.display_name != "" ? var.spec.display_name : null

  # Edition gates feature availability (ENTERPRISE_PLUS unlocks
  # multi-location automated-backup placement on tables). Unset lets the
  # provider apply its ENTERPRISE default; upgrades apply in place, there
  # is no downgrade path.
  edition = var.spec.edition != "" ? var.spec.edition : null

  # Resource Manager tags for org-policy and IAM conditions. Create-time
  # only (ForceNew): a tag change replaces the instance.
  tags = length(var.spec.resource_manager_tags) > 0 ? var.spec.resource_manager_tags : null

  # What a PERMITTED destroy does once deletion_protection allows one:
  # DELETE (default), PREVENT (destroy fails), or ABANDON (drop from
  # state, keep the instance running in GCP).
  deletion_policy = var.spec.deletion_policy != "" ? var.spec.deletion_policy : null

  dynamic "cluster" {
    for_each = var.spec.clusters
    content {
      cluster_id = cluster.value.cluster_id
      zone       = cluster.value.zone
      # Fixed node count only applies without autoscaling (mutually
      # exclusive, enforced by spec CEL); 0 means "let Bigtable
      # auto-allocate for the data footprint".
      num_nodes           = cluster.value.autoscaling_config == null ? (cluster.value.num_nodes > 0 ? cluster.value.num_nodes : null) : null
      storage_type        = cluster.value.storage_type
      kms_key_name        = cluster.value.kms_key_name != "" ? cluster.value.kms_key_name : null
      node_scaling_factor = cluster.value.node_scaling_factor != "" ? cluster.value.node_scaling_factor : null

      dynamic "autoscaling_config" {
        for_each = cluster.value.autoscaling_config != null ? [cluster.value.autoscaling_config] : []
        content {
          min_nodes  = autoscaling_config.value.min_nodes
          max_nodes  = autoscaling_config.value.max_nodes
          cpu_target = autoscaling_config.value.cpu_target
          # 0 leaves the storage target to the per-storage-type server
          # default (2560 GB SSD / 8192 GB HDD).
          storage_target = autoscaling_config.value.storage_target > 0 ? autoscaling_config.value.storage_target : null
        }
      }
    }
  }

  depends_on = [google_project_service.bigtableadmin_api]
}
