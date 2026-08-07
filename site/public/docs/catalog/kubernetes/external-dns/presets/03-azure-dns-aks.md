---
title: "Azure DNS on AKS (Keyless via Workload Identity)"
description: "This preset installs ExternalDNS on an AKS cluster publishing to an Azure DNS zone, authenticating keylessly through Azure AD Workload Identity — no service-principal secret anywhere. The module..."
type: "preset"
rank: "03"
presetSlug: "03-azure-dns-aks"
componentSlug: "external-dns"
componentTitle: "External DNS"
provider: "kubernetes"
icon: "package"
order: 3
---

# Azure DNS on AKS (Keyless via Workload Identity)

This preset installs ExternalDNS on an AKS cluster publishing to an Azure
DNS zone, authenticating keylessly through Azure AD Workload Identity — no
service-principal secret anywhere. The module renders the controller's
`azure.json` from these fields (Workload Identity mode) and stamps the
`azure.workload.identity/use` pod label the identity webhook requires. This
is the standard production posture on Azure.

## When to Use

- AKS clusters (with Workload Identity enabled) whose Services/Ingresses
  publish records into an Azure DNS zone
- Zones dedicated to (or safely shareable with) this cluster
- Production deployments — keyless Workload Identity is the recommended
  authentication

## Key Configuration Choices

- **Keyless authentication** (`workloadIdentity.aks.clientId`) — a
  user-assigned managed identity with a federated credential for the
  controller ServiceAccount; no service-principal secret is created
- **Public zones** (`privateZones` unset) — set `privateZones: true` to
  manage Azure Private DNS zones instead
- **Zone scoping** (`zoneIdFilters` + `domainFilters`) — restricts
  management to one zone within the resource group
- **`policy: sync`** — full reconciliation; the TXT registry limits deletes
  to records tagged with this instance's `txtOwnerId`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<dns-zone-resource-group>` | Resource group containing the DNS zones | Azure portal |
| `<azure-subscription-id>` | Subscription owning the zones | Azure portal |
| `<azure-dns-zone-id>` | DNS zone resource ID | Azure portal or `AzureDnsZone` outputs |
| `<managed-identity-client-id>` | Client ID of the managed identity with "DNS Zone Contributor", carrying a federated credential for `system:serviceaccount:external-dns:my-external-dns-azure` | Azure portal or `AzureUserAssignedIdentity` outputs |
| `<cluster-name>` | Unique owner ID for this instance | Your cluster naming |
| `<example.com>` | Domain suffix this instance manages | Your zone's domain |

## Related Presets

- **01-aws-route53-eks-keyless** — the same posture on EKS + Route 53
- **02-google-cloud-dns-gke** — the same posture on GKE + Cloud DNS
