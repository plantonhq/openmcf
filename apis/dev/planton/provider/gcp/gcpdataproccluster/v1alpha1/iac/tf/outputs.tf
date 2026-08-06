# Fully qualified cluster resource name
# (projects/{project}/regions/{region}/clusters/{cluster}) — the
# composition handle downstream resources reference, including another
# cluster's spark_history_server_config which consumes this exact format.
output "cluster_id" {
  description = "Fully qualified cluster resource name"
  value       = google_dataproc_cluster.cluster.id
}

output "cluster_name" {
  description = "Short name of the cluster"
  value       = google_dataproc_cluster.cluster.name
}

# The staging bucket actually in use: the user-supplied bucket when one
# was referenced, otherwise the bucket GCP auto-created. The virtual arm
# reports its own staging bucket the same way.
output "staging_bucket" {
  description = "Cloud Storage bucket used for staging job dependencies"
  value = try(
    google_dataproc_cluster.cluster.cluster_config[0].bucket,
    try(google_dataproc_cluster.cluster.virtual_cluster_config[0].staging_bucket, "")
  )
}
