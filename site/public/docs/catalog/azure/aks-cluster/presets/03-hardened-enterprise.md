---
title: "Hardened Enterprise AKS Cluster"
description: "This preset deploys a compliance-posture AKS cluster: Azure AD RBAC with local accounts disabled, API-server access restricted to authorized ranges, Cilium-compatible Azure network policy, host..."
type: "preset"
rank: "03"
presetSlug: "03-hardened-enterprise"
componentSlug: "aks-cluster"
componentTitle: "AKS Cluster"
provider: "azure"
icon: "package"
order: 3
---

# Hardened Enterprise AKS Cluster

This preset deploys a compliance-posture AKS cluster: Azure AD RBAC with local accounts disabled, API-server access restricted to authorized ranges, Cilium-compatible Azure network policy, host encryption on the system pool, Microsoft Defender for Containers, cost analysis, and a pinned Kubernetes version.

## When to Use

- Enterprises with compliance requirements (SOC 2, PCI, HIPAA-adjacent postures)
- Organizations standardizing on Azure AD as the single Kubernetes auth source
- Environments that must restrict API access to known networks without going fully private
- Teams that need security-event streaming and cost attribution from day one

## Key Configuration Choices

- **AAD RBAC, local accounts off** (`azureActiveDirectoryRoleBasedAccessControl.azureRbacEnabled: true`, `localAccountDisabled: true`) -- Azure AD is the only way in; Kubernetes RBAC is managed through Azure role assignments
- **Authorized IP ranges** (`apiServerAccessProfile.authorizedIpRanges`) -- public endpoint restricted to your networks
- **Pinned version** (`kubernetesVersion`) -- upgrades happen when you choose, not when you redeploy
- **Network policy** (`networkProfile.networkPolicy: AZURE`) -- pod-level traffic rules enforced
- **Host encryption** (`defaultNodePool.hostEncryptionEnabled: true`) -- temp disks and caches encrypted at rest (requires the EncryptionAtHost feature on the subscription)
- **Defender + OMS** -- security events and logs stream to your Log Analytics workspace
- **Cost analysis** (`costAnalysisEnabled: true`) -- namespace-level cost attribution (requires Standard tier)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<pinned-k8s-version>` | Kubernetes version to pin (e.g., `1.35`) | `az aks get-versions --location <region>` |
| `203.0.113.0/24` | CIDR block allowed to reach the API server | Your network team |
| `<nodes-subnet-id>` | ARM ID of the subnet for cluster nodes | Azure portal or an `AzureSubnet`'s `status.outputs.subnet_id` |
| `<log-analytics-workspace-id>` | ARM ID of the Log Analytics workspace | Azure portal or `AzureLogAnalyticsWorkspace` status outputs |

## Related Presets

- **01-standard** -- Use instead for standard production workloads without the compliance hardening
- **02-private** -- Use instead (or combine) when the API server must be fully private
