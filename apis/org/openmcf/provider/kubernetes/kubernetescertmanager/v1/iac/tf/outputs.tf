# Outputs for Cert-Manager deployment
# These outputs match the KubernetesCertManagerStackOutputs protobuf message

output "namespace" {
  description = "Kubernetes namespace where cert-manager was deployed"
  value       = local.namespace_name
}

output "release_name" {
  description = "Helm release name (useful for upgrades)"
  value       = local.release_name
}

output "solver_identity" {
  description = "The service account email/ARN/ClientID used for DNS-01 solver"
  value       = local.solver_identity
}

output "cloudflare_secret_name" {
  description = "The name of the Kubernetes Secret containing the Cloudflare API token (only set for Cloudflare configs)"
  value       = local.cloudflare_secret_name
}

output "cluster_issuer_names" {
  description = "List of ClusterIssuer names (one per domain)"
  value       = [for issuer in local.cluster_issuers : issuer.domain]
}

