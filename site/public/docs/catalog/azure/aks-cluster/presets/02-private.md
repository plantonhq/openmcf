---
title: "Private AKS Cluster with Workload Identity"
description: "This preset deploys a private AKS cluster with no public API server endpoint and the workload-identity loop enabled. The Kubernetes API is accessible only from within the VNet or via peered networks..."
type: "preset"
rank: "02"
presetSlug: "02-private"
componentSlug: "aks-cluster"
componentTitle: "AKS Cluster"
provider: "azure"
icon: "package"
order: 2
---

# Private AKS Cluster with Workload Identity

This preset deploys a private AKS cluster with no public API server endpoint and the workload-identity loop enabled. The Kubernetes API is accessible only from within the VNet or via peered networks (VPN, ExpressRoute); nodes deploy into your own subnet.

## When to Use

- Regulated or security-sensitive environments that prohibit public Kubernetes API endpoints
- Clusters accessed exclusively via VPN, ExpressRoute, or Azure Bastion
- Pods that authenticate to Azure services keylessly through federated credentials
- Enterprise environments with strict network perimeter policies

## Key Configuration Choices

- **Private cluster** (`privateClusterEnabled: true`) -- API server has no public IP; AKS manages the private DNS zone (reference an `AzurePrivateDnsZone` via `privateDnsZoneId` to bring your own)
- **Standard tier** (`skuTier: STANDARD`) -- financially-backed 99.95% uptime SLA
- **BYO subnet** (`defaultNodePool.vnetSubnetId`) -- private clusters place nodes in your network; the cluster identity needs Network Contributor on the subnet
- **Workload identity** (`oidcIssuerEnabled: true`, `workloadIdentityEnabled: true`) -- the cluster's `oidc_issuer_url` output feeds `AzureFederatedIdentityCredential` for secret-less pod auth
- **System pool isolated** (`onlyCriticalAddonsEnabled: true`, autoscaling 3-5 across 3 zones) -- workload pools are separate `AzureAksNodePool` resources

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (e.g., `eastus`, `westeurope`) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<nodes-subnet-id>` | ARM ID of the subnet for cluster nodes | Azure portal or an `AzureSubnet`'s `status.outputs.subnet_id` |

## Related Presets

- **01-standard** -- Use instead when a public API endpoint is acceptable
- **03-hardened-enterprise** -- Use instead when you also need AAD RBAC, Defender, and host encryption
