# Fully qualified resource name
# (projects/{p}/locations/{r}/environments/{e}/userWorkloadsSecrets/{n}).
# The Secret's data is deliberately never exported.
output "name" {
  description = "Fully qualified user workloads Secret resource name"
  value       = google_composer_user_workloads_secret.secret.id
}

# The Kubernetes Secret name — what KubernetesPodOperator tasks and
# Airflow connections reference.
output "secret_name" {
  description = "The Kubernetes Secret name"
  value       = google_composer_user_workloads_secret.secret.name
}
