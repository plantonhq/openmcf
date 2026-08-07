# AzureAksCluster Terraform Module

## Overview

This Terraform module provisions an Azure Kubernetes Service managed
cluster using the `azurerm` provider (`~> 4.0`). It creates a single
`azurerm_kubernetes_cluster` carrying the control plane, the mandatory
default (system) node pool, the cluster identity, the network fabric, and
the Azure-managed add-ons -- the full azurerm v4.80 surface.

The cluster carries exactly ONE node pool. Every additional pool is a
standalone `AzureAksNodePool` resource referencing this cluster's
`cluster_id` output, so pools have independent lifecycles.

The network fabric (plugin/mode/CIDRs/outbound type) mostly replaces the
cluster when changed; the module writes the modern AKS default (Azure CNI
overlay) explicitly when the spec leaves it unset, because azurerm's
implicit fallback (kubenet) is deprecated. Many default-pool shape changes
rotate the pool -- `temporary_name_for_rotation` lets AKS stand up a
replacement first.

## Resources Created

- `azurerm_kubernetes_cluster.main` -- the managed cluster (control plane +
  default node pool + add-ons)

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Cluster specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region for the control plane and default pool |
| `resource_group` | yes | Resource group name |
| `name` | yes | Cluster name, unique within the resource group |
| `default_node_pool` | yes | The mandatory system pool: `name`, `vm_size`, count/autoscaling, zones, optional subnet FKs, disks, kubelet/Linux OS config, upgrade settings, tags |
| `kubernetes_version` | no | Unset provisions the latest AKS-recommended GA version; pin for production |
| `sku_tier` / `support_plan` | no | Control-plane tier (`FREE`/`STANDARD`/`PREMIUM`) and support plan |
| `identity` / `kubelet_identity` | no | System-assigned (default) or user-assigned managed identity |
| `azure_active_directory_role_based_access_control` | no | AAD RBAC: admin group IDs, tenant, Azure-RBAC-for-Kubernetes |
| `api_server_access_profile` / `private_cluster_enabled` / `private_dns_zone_id` | no | API-server exposure: authorized ranges or a private cluster |
| `network_profile` | no | Plugin/mode/policy/data plane, CIDRs, outbound type, LB/NAT profiles, ACNS |
| `oidc_issuer_enabled` / `workload_identity_enabled` | no | Workload-identity trust anchor (issuer defaults ON) |
| `auto_scaler_profile` | no | Cluster-wide autoscaler tuning |
| `maintenance_window*` | no | Legacy, auto-upgrade, and node-OS maintenance windows |
| `oms_agent`, `key_vault_secrets_provider`, `azure_policy_enabled`, `microsoft_defender`, `monitor_metrics`, `ingress_application_gateway`, `aci_connector_linux`, `confidential_computing`, `web_app_routing` | no | Add-ons |
| `service_mesh_profile`, `storage_profile`, `workload_autoscaler_profile`, `key_management_service`, `http_proxy_config`, `linux_profile`, `windows_profile`, `bootstrap_profile`, `node_provisioning_profile`, `upgrade_override` | no | Platform profiles |
| `disk_encryption_set_id`, `edge_zone`, `node_resource_group`, `custom_ca_trust_certificates_base64`, `tags` | no | Misc |

## Outputs

| Output | Description |
|--------|-------------|
| `cluster_id` | Full ARM ID of the managed cluster (the node-pool parent seam) |
| `cluster_name` | The cluster's name as deployed |
| `fqdn` / `private_fqdn` / `portal_fqdn` | API-server FQDNs |
| `oidc_issuer_url` | OIDC issuer URL (workload-identity trust anchor) |
| `node_resource_group` / `node_resource_group_id` | The Azure-managed node resource group |
| `cluster_kubeconfig` | Base64-encoded kubeconfig (sensitive) |
| `cluster_identity_principal_id` | Cluster managed-identity principal ID |
| `kubelet_identity_object_id` / `kubelet_identity_client_id` | Kubelet identity (grant AcrPull on registries) |
| `current_kubernetes_version` | Version the control plane is actually running |

## Usage

```hcl
module "aks_cluster" {
  source = "./iac/tf"

  metadata = {
    name = "prod-aks"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    region         = "eastus"
    resource_group = "prod-rg"
    name           = "prod-aks"
    default_node_pool = {
      name       = "system"
      vm_size    = "Standard_D4s_v5"
      node_count = 3
    }
  }
}
```

## Required Permissions

The deploying credential needs
`Microsoft.ContainerService/managedClusters/write` on the resource group --
held via Azure Kubernetes Service Contributor, Contributor, or Owner. A
BYO subnet additionally requires the cluster identity to hold Network
Contributor on the subnet; a BYO private DNS zone requires Private DNS
Zone Contributor.
