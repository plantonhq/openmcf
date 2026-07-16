# Semantic stack outputs, matching AzureAksClusterStackOutputs field for
# field. cluster_id is the parent seam every standalone AzureAksNodePool
# consumes; oidc_issuer_url is the trust anchor an
# AzureFederatedIdentityCredential binds to for workload identity.

output "cluster_id" {
  description = "The Azure Resource Manager ID of the managed cluster."
  value       = azurerm_kubernetes_cluster.main.id
}

output "cluster_name" {
  description = "The name of the managed cluster."
  value       = azurerm_kubernetes_cluster.main.name
}

output "fqdn" {
  description = "The public FQDN of the Kubernetes API server (empty for private clusters without a public FQDN)."
  value       = azurerm_kubernetes_cluster.main.fqdn
}

output "private_fqdn" {
  description = "The private FQDN of the API server, populated for private clusters."
  value       = azurerm_kubernetes_cluster.main.private_fqdn
}

output "portal_fqdn" {
  description = "The FQDN used by the Azure Portal to reach the cluster."
  value       = azurerm_kubernetes_cluster.main.portal_fqdn
}

output "oidc_issuer_url" {
  description = "The cluster's OIDC issuer URL -- consumed by AzureFederatedIdentityCredential as its issuer."
  value       = azurerm_kubernetes_cluster.main.oidc_issuer_url
}

output "node_resource_group" {
  description = "The name of the Azure-managed node resource group."
  value       = azurerm_kubernetes_cluster.main.node_resource_group
}

output "node_resource_group_id" {
  description = "The ARM ID of the node resource group."
  value       = azurerm_kubernetes_cluster.main.node_resource_group_id
}

output "cluster_kubeconfig" {
  description = "Base64-encoded kubeconfig for the cluster (user credential). Treat as a secret."
  value       = base64encode(azurerm_kubernetes_cluster.main.kube_config_raw)
  sensitive   = true
}

output "cluster_identity_principal_id" {
  description = "The principal (object) ID of the cluster's managed identity."
  value       = try(azurerm_kubernetes_cluster.main.identity[0].principal_id, "")
}

output "kubelet_identity_object_id" {
  description = "The object ID of the kubelet identity -- grant it AcrPull on registries."
  value       = try(azurerm_kubernetes_cluster.main.kubelet_identity[0].object_id, "")
}

output "kubelet_identity_client_id" {
  description = "The client ID of the kubelet identity."
  value       = try(azurerm_kubernetes_cluster.main.kubelet_identity[0].client_id, "")
}

output "current_kubernetes_version" {
  description = "The Kubernetes version the control plane is actually running."
  value       = azurerm_kubernetes_cluster.main.current_kubernetes_version
}
