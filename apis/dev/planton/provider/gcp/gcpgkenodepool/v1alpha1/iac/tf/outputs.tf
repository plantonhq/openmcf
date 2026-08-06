output "node_pool_name" {
  description = "Name of the node pool as created in GKE"
  value       = google_container_node_pool.this.name
}

output "instance_group_urls" {
  description = "Resource URLs of the managed instance groups backing this pool (one per zone for regional clusters)"
  value       = google_container_node_pool.this.instance_group_urls
}

output "min_nodes" {
  description = "Effective minimum size of the pool (autoscaling minimum, or the fixed node_count)"
  value       = local.effective_min
}

output "max_nodes" {
  description = "Effective maximum size of the pool (autoscaling maximum, or the fixed node_count)"
  value       = local.effective_max
}

output "current_node_count" {
  description = "Nodes per zone at the last deploy (a snapshot — the autoscaler moves it at runtime)"
  value       = google_container_node_pool.this.node_count
}

output "node_pool_id" {
  description = "Fully qualified node pool resource ID (projects/{project}/locations/{location}/clusters/{cluster}/nodePools/{name})"
  value       = google_container_node_pool.this.id
}

output "location" {
  description = "The pool's location (the parent cluster's region or zone), exactly as provided in the spec"
  value       = var.spec.location
}

output "version" {
  description = "The Kubernetes version running on the pool's nodes"
  value       = google_container_node_pool.this.version
}
