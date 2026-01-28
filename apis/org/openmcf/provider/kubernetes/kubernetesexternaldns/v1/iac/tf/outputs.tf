###########################
# outputs.tf
###########################

output "namespace" {
  description = "Namespace where ExternalDNS is deployed"
  value       = local.namespace_name
}

output "release_name" {
  description = "Helm release name for ExternalDNS"
  value       = helm_release.external_dns.name
}

output "service_account_name" {
  description = "Kubernetes service account name for ExternalDNS"
  value       = kubernetes_service_account.external_dns.metadata[0].name
}

output "provider_type" {
  description = "DNS provider type (google, aws, azure, cloudflare)"
  value       = local.provider_type
}

output "gke_service_account_email" {
  description = "Google Service Account email for Workload Identity (GKE only)"
  value       = local.is_gke ? local.gke_gsa_email : null
}

output "cloudflare_secret_name" {
  description = "Kubernetes secret name for Cloudflare API token (Cloudflare only)"
  value       = local.is_cloudflare ? local.cloudflare_api_token_secret_name : null
}

