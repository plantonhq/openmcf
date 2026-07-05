output "job_name" {
  description = "Name of the Cloud Run job as created in GCP"
  value       = google_cloud_run_v2_job.main.name
}

output "location" {
  description = "Region the job is deployed in"
  value       = google_cloud_run_v2_job.main.location
}

output "uid" {
  description = "Server-assigned unique identifier of the job"
  value       = google_cloud_run_v2_job.main.uid
}

output "latest_created_execution" {
  description = "Name of the most recently created execution, if any"
  value       = try(google_cloud_run_v2_job.main.latest_created_execution[0].name, "")
}
