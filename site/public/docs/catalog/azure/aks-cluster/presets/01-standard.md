---
title: "Standard Production AKS Cluster"
description: "This preset deploys a production-ready AKS cluster with a public API endpoint, Azure CNI Overlay networking, a 3-zone autoscaling system pool tainted for system pods only, and the recommended add-ons..."
type: "preset"
rank: "01"
presetSlug: "01-standard"
componentSlug: "aks-cluster"
componentTitle: "AKS Cluster"
provider: "azure"
icon: "package"
order: 1
---

# Standard Production AKS Cluster

This preset deploys a production-ready AKS cluster with a public API endpoint, Azure CNI Overlay networking, a 3-zone autoscaling system pool tainted for system pods only, and the recommended add-ons (Container Insights via OMS agent, Key Vault secrets provider with rotation, Azure Policy, workload identity). Application workloads run in separately-deployed `AzureAksNodePool` resources.

## When to Use

- Production Kubernetes clusters that need a public API endpoint
- Standard web, API, and microservice workloads on Azure CNI Overlay
- Teams that want Azure-managed add-ons out of the box
- Clusters requiring the 99.95% uptime SLA with availability-zone distribution

## Key Configuration Choices

- **Standard tier** (`skuTier: STANDARD`) -- financially-backed 99.95% uptime SLA with availability zones
- **Azure CNI Overlay** (`networkProfile.networkPlugin: AZURE_CNI`, `networkPluginMode: OVERLAY`) -- pods get private-CIDR IPs, avoiding VNet IP exhaustion while retaining VNet integration
- **Public endpoint** (`privateClusterEnabled: false`) -- restrict with `apiServerAccessProfile.authorizedIpRanges` if needed
- **System pool isolated** (`defaultNodePool` with `onlyCriticalAddonsEnabled: true`, autoscaling 3-5 across 3 zones) -- CoreDNS and metrics-server never compete with app pods; workload pools are separate `AzureAksNodePool` resources
- **Managed networking** -- no subnet reference; AKS provisions its own network (reference an `AzureSubnet` on the default pool for BYO networking)
- **Monitoring + governance on** -- OMS agent with managed-identity auth, Key Vault CSI with secret rotation, Azure Policy, workload identity

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (e.g., `eastus`, `westeurope`) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<log-analytics-workspace-id>` | ARM ID of the Log Analytics workspace | Azure portal or `AzureLogAnalyticsWorkspace` status outputs |

## Related Presets

- **02-private** -- Use instead when the API server must not be publicly accessible
- **03-hardened-enterprise** -- Use instead for compliance environments needing AAD RBAC, authorized ranges, Defender, and host encryption
