# Terraform outputs for Kubernetes ConfigMap

output "configmap_name" {
  description = "The name of the created Kubernetes ConfigMap"
  value       = kubernetes_config_map_v1.configmap.metadata[0].name
}

output "namespace" {
  description = "The namespace where the configmap was created"
  value       = kubernetes_config_map_v1.configmap.metadata[0].namespace
}
