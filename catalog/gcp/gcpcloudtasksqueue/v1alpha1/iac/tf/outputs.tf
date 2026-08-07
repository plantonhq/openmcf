output "queue_id" {
  description = "The fully qualified queue ID (projects/{project}/locations/{location}/queues/{name})"
  value       = google_cloud_tasks_queue.this.id
}

output "queue_name" {
  description = "The short queue name"
  value       = google_cloud_tasks_queue.this.name
}

# GCP computes max_burst_size from max_dispatches_per_second; it is not
# configurable, so the effective value is surfaced as an output. The
# rate_limits block is Optional+Computed, so the API populates it even when
# the manifest omits rate limits entirely.
output "max_burst_size" {
  description = "The effective max burst size computed by GCP from the queue's dispatch rate"
  value       = google_cloud_tasks_queue.this.rate_limits[0].max_burst_size
}
