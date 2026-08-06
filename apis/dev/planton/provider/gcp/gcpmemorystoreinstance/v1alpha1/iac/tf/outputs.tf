# IP address of the discovery endpoint — the address clients connect to
# for topology discovery and command routing, extracted from the PSC
# auto-created endpoint connections.
output "discovery_address" {
  description = "IP address of the instance's discovery endpoint"
  value = try(
    google_memorystore_instance.this.endpoints[0].connections[0].psc_auto_connection[0].ip_address,
    ""
  )
}

# Port of the discovery endpoint (typically 6379).
output "discovery_port" {
  description = "Port of the instance's discovery endpoint"
  value = try(
    google_memorystore_instance.this.endpoints[0].connections[0].psc_auto_connection[0].port,
    0
  )
}

output "instance_uid" {
  description = "Server-generated unique identifier for the instance"
  value       = google_memorystore_instance.this.uid
}

# Memory per node is a consequence of node_type — reported for capacity
# planning rather than configured.
output "node_size_gb" {
  description = "Memory size per node in GB"
  value = try(
    google_memorystore_instance.this.node_config[0].size_gb,
    0
  )
}

# Full resource path — the composition key for cross-instance
# replication (a SECONDARY's primary_instance reference resolves to it).
output "name" {
  description = "Full resource path of the instance"
  value       = google_memorystore_instance.this.name
}

# Where managed backups for this instance live once automated backups
# are configured — the source of managed_backup_source paths for seeding
# new instances.
output "backup_collection" {
  description = "Full resource path of the instance's backup collection"
  value       = google_memorystore_instance.this.backup_collection
}
