output "function_id" {
  description = "Fully qualified resource name (projects/{project}/locations/{region}/functions/{name})"
  value       = google_cloudfunctions2_function.function.id
}

output "function_url" {
  description = "HTTPS URL of the function (the cloudfunctions.net / run.app endpoint)"
  value       = google_cloudfunctions2_function.function.url
}

output "service_account_email" {
  description = "Email of the service account the function runs as"
  value       = try(google_cloudfunctions2_function.function.service_config[0].service_account_email, "")
}

output "state" {
  description = "Current state of the function (ACTIVE, OFFLINE, DEPLOY_IN_PROGRESS, ...)"
  value       = google_cloudfunctions2_function.function.state
}

output "cloud_run_service_id" {
  description = "The underlying Cloud Run service serving this Gen 2 function"
  value       = try(google_cloudfunctions2_function.function.service_config[0].service, "")
}

output "eventarc_trigger_id" {
  description = "Eventarc trigger resource name; empty for HTTP functions"
  value       = try(google_cloudfunctions2_function.function.event_trigger[0].trigger, "")
}

output "name" {
  description = "Bare function name — the handle serverless NEGs and gcloud reference"
  value       = google_cloudfunctions2_function.function.name
}

output "uri" {
  description = "URI of the underlying Cloud Run service (*.run.app)"
  value       = try(google_cloudfunctions2_function.function.service_config[0].uri, "")
}

output "environment" {
  description = "The environment the function runs in (e.g. GEN_2)"
  value       = google_cloudfunctions2_function.function.environment
}

output "update_time" {
  description = "Timestamp of the last update to the function"
  value       = google_cloudfunctions2_function.function.update_time
}
