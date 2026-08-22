resource "digitalocean_kubernetes_node_pool" "node_pool" {
  # The cluster reference arrives flattened: the owning DOKS cluster's UUID.
  # Changing it replaces the pool (provider ForceNew).
  cluster_id = var.spec.cluster

  name = var.spec.node_pool_name

  # Changing the size replaces the pool (provider ForceNew).
  size = var.spec.size

  # With auto_scale enabled this is the initial count; the provider then
  # suppresses diffs while the live count drifts between the bounds.
  node_count = var.spec.node_count

  auto_scale = var.spec.auto_scale
  min_nodes  = var.spec.auto_scale ? var.spec.min_nodes : null
  max_nodes  = var.spec.auto_scale ? var.spec.max_nodes : null

  # Kubernetes node labels: user labels over the standard Planton labels
  # (identical set in both provisioners).
  labels = local.node_labels

  # Pool-level Droplet tags: spec tags plus the standard Planton tags.
  tags = local.tags

  # AMD GPU partitioning; changing it replaces the pool (provider ForceNew).
  gpu_partition_mode = var.spec.gpu_partition_mode != "" ? var.spec.gpu_partition_mode : null

  dynamic "taint" {
    for_each = var.spec.taints
    content {
      key = taint.value.key
      # Kubernetes allows valueless taints; the provider requires the value
      # be sent, possibly empty.
      value  = taint.value.value
      effect = taint.value.effect
    }
  }
}
