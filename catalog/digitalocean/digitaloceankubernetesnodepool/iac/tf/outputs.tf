output "node_pool_id" {
  description = "The unique identifier (UUID) of the created node pool"
  value       = digitalocean_kubernetes_node_pool.node_pool.id
}

output "node_ids" {
  description = "The DOKS node object UUIDs of the pool's current members"
  value       = digitalocean_kubernetes_node_pool.node_pool.nodes[*].id
}

output "cluster_id" {
  description = "The UUID of the cluster that owns this pool"
  value       = digitalocean_kubernetes_node_pool.node_pool.cluster_id
}

output "droplet_ids" {
  description = "The integer ids (as strings) of the Droplets backing the pool's nodes"
  value       = digitalocean_kubernetes_node_pool.node_pool.nodes[*].droplet_id
}
