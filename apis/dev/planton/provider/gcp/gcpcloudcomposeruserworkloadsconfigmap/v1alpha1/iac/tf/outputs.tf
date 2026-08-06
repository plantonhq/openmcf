# Fully qualified resource name
# (projects/{p}/locations/{r}/environments/{e}/userWorkloadsConfigMaps/{n}).
output "name" {
  description = "Fully qualified user workloads ConfigMap resource name"
  value       = google_composer_user_workloads_config_map.config_map.id
}

# The Kubernetes ConfigMap name — what KubernetesPodOperator tasks
# reference.
output "config_map_name" {
  description = "The Kubernetes ConfigMap name"
  value       = google_composer_user_workloads_config_map.config_map.name
}
