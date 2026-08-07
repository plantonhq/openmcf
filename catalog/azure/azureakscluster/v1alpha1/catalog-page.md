# Azure AKS Cluster

Creates an Azure Kubernetes Service (AKS) managed cluster -- the control plane, its identity and access model, its network fabric, the mandatory default (system) node pool, and the Azure-managed add-ons -- at the full `azurerm` v4.80 surface.

## What Gets Created

When you deploy an AzureAksCluster resource, Planton provisions:

- **Managed Cluster** — an `azurerm_kubernetes_cluster` carrying the control plane, the default node pool, identity, networking, and add-ons

The cluster deliberately carries exactly ONE node pool -- the default (system) pool Azure requires at creation. Every additional pool is its own composable `AzureAksNodePool` resource referencing this cluster's `cluster_id` output, so pools scale, upgrade, and rotate without ever touching the cluster resource.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A resource group** to create the cluster in (an `AzureResourceGroup` in composed environments)
- **Container service write rights**: `Microsoft.ContainerService/managedClusters/write` (Azure Kubernetes Service Contributor, Contributor, or Owner)
- Optional: an `AzureSubnet` when bringing your own network (leave the subnet unset and AKS provisions managed networking)

## Quick Start

Create a file `cluster.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureAksCluster
metadata:
  name: prod-aks
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureAksCluster.prod-aks
spec:
  region: eastus
  resourceGroup:
    value: prod-rg
  name: prod-aks
  defaultNodePool:
    name: system
    vmSize: Standard_D4s_v5
    nodeCount: 3
```

Deploy:

```shell
planton apply -f cluster.yaml
```

This is Azure's own simplest honest cluster: a resource group, a default pool, and AKS-managed networking. No subnet, no version pin -- AKS provisions the latest recommended GA Kubernetes version into a network it manages itself.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the control plane and default pool. Changing it replaces the cluster. | Required |
| `resourceGroup` | `StringValueOrRef` | Resource group name. Defaults to referencing an `AzureResourceGroup`'s name output. | Required |
| `name` | `string` | Cluster name, unique within the resource group. | Required, 1-63 chars, Azure naming rules |
| `defaultNodePool` | `object` | The mandatory default (system) pool. See below. | Required |

### Default Node Pool

| Field | Type | Description |
|-------|------|-------------|
| `defaultNodePool.name` | `string` | Pool name: 1-12 lowercase letters and numbers, starting with a letter. Required. |
| `defaultNodePool.vmSize` | `string` | Azure VM size, e.g. `Standard_D4s_v5`. Required. |
| `defaultNodePool.nodeCount` | `int` | Fixed node count (1-1000) without autoscaling; initial count with it. |
| `defaultNodePool.autoScalingEnabled` | `bool` | Let the cluster autoscaler own the count between `minCount` and `maxCount`. |
| `defaultNodePool.minCount` / `maxCount` | `int` | Autoscaling bounds (1-1000; the default pool cannot scale to zero). |
| `defaultNodePool.zones` | `list(string)` | Availability zones (`"1"`, `"2"`, `"3"`). Empty = regional placement. |
| `defaultNodePool.vnetSubnetId` | `StringValueOrRef` | Optional BYO subnet (references an `AzureSubnet`). Unset = AKS-managed networking. |
| `defaultNodePool.podSubnetId` | `StringValueOrRef` | Separate pod subnet for non-overlay Azure CNI with dynamic pod IPs. |
| `defaultNodePool.onlyCriticalAddonsEnabled` | `bool` | Taint the pool `CriticalAddonsOnly` so only system pods schedule here -- the recommended posture once workload pools exist. |
| `defaultNodePool.*` | | Disk (size/type/kubelet disk/ultra SSD), OS SKU, orchestrator version, labels, FIPS, host encryption, node public IPs (+ prefix FK), GPU instance/driver, placement groups, scale-down mode, snapshot, workload runtime, `temporaryNameForRotation`, kubelet config, Linux OS/sysctl config, node network profile, upgrade settings, tags. |

### Cluster-Wide Optional Fields (grouped)

| Group | Fields |
|-------|--------|
| Versioning & tier | `kubernetesVersion` (unset = latest AKS-recommended GA; PIN for production), `skuTier` (`FREE`/`STANDARD`/`PREMIUM`), `supportPlan`, `automaticUpgradeChannel`, `nodeOsUpgradeChannel`, `upgradeOverride` |
| DNS | `dnsPrefix` (unset = derived from name), `dnsPrefixPrivateCluster` (private clusters; mutually exclusive) |
| Identity & access | `identity` (system- or user-assigned, `identityIds` reference `AzureUserAssignedIdentity`), `kubeletIdentity`, `azureActiveDirectoryRoleBasedAccessControl` (AAD RBAC; `adminGroupObjectIds`/`tenantId`/`azureRbacEnabled`), `roleBasedAccessControlEnabled`, `localAccountDisabled` |
| API server | `apiServerAccessProfile.authorizedIpRanges`, `privateClusterEnabled`, `privateDnsZoneId` (references `AzurePrivateDnsZone`), `privateClusterPublicFqdnEnabled` |
| Networking | `networkProfile`: plugin/mode/policy/data plane, pod & service CIDRs (dual-stack), outbound type, load-balancer + NAT-gateway profiles, `advancedNetworking` (ACNS observability/security -- requires Cilium) |
| Workload identity | `oidcIssuerEnabled` (default ON -- publishes `oidc_issuer_url`), `workloadIdentityEnabled` |
| Autoscaler | `autoScalerProfile` (cluster-wide autoscaler tuning) |
| Maintenance | `maintenanceWindow`, `maintenanceWindowAutoUpgrade`, `maintenanceWindowNodeOs` |
| Add-ons | `omsAgent` (references `AzureLogAnalyticsWorkspace`), `keyVaultSecretsProvider`, `azurePolicyEnabled`, `ingressApplicationGateway` (references `AzureApplicationGateway`), `microsoftDefender`, `monitorMetrics`, `aciConnectorLinux`, `confidentialComputing`, `webAppRouting` |
| Platform | `serviceMeshProfile` (Istio), `storageProfile` (CSI drivers), `workloadAutoscalerProfile` (KEDA/VPA), `keyManagementService` (etcd CMK), `httpProxyConfig`, `linuxProfile`, `windowsProfile` (prerequisite for Windows pools), `imageCleaner*`, `costAnalysisEnabled`, `runCommandEnabled`, `bootstrapProfile` (references `AzureContainerRegistry`), `nodeProvisioningProfile` (NAP/Karpenter), `aiToolchainOperatorEnabled` |
| Misc | `diskEncryptionSetId`, `edgeZone`, `nodeResourceGroup`, `customCaTrustCertificatesBase64`, `tags` |

Secret-by-default: `windowsProfile.adminPassword` and `httpProxyConfig.trustedCa` are `sensitive` -- Planton forces them to managed-secret references.

## Examples

### Private Cluster with Workload Identity

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureAksCluster
metadata:
  name: private-aks
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureAksCluster.private-aks
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: private-aks
  skuTier: STANDARD
  privateClusterEnabled: true
  oidcIssuerEnabled: true
  workloadIdentityEnabled: true
  defaultNodePool:
    name: system
    vmSize: Standard_D4s_v5
    autoScalingEnabled: true
    minCount: 3
    maxCount: 5
    vnetSubnetId:
      valueFrom:
        name: nodes-subnet
    onlyCriticalAddonsEnabled: true
```

### Monitored Production Cluster with Azure CNI Overlay

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureAksCluster
metadata:
  name: prod-aks
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureAksCluster.prod-aks
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: prod-aks
  kubernetesVersion: "1.35"
  skuTier: STANDARD
  networkProfile:
    networkPlugin: AZURE_CNI
    networkPluginMode: OVERLAY
  defaultNodePool:
    name: system
    vmSize: Standard_D4s_v5
    autoScalingEnabled: true
    minCount: 3
    maxCount: 5
    zones: ["1", "2", "3"]
    onlyCriticalAddonsEnabled: true
  omsAgent:
    logAnalyticsWorkspaceId:
      valueFrom:
        name: prod-logs
    msiAuthForMonitoringEnabled: true
  keyVaultSecretsProvider:
    secretRotationEnabled: true
  azurePolicyEnabled: true
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `cluster_id` | `string` | Full ARM ID of the managed cluster -- the parent seam every `AzureAksNodePool` consumes |
| `cluster_name` | `string` | The cluster's name as deployed |
| `fqdn` | `string` | Public API-server FQDN (empty for private-only clusters) |
| `private_fqdn` | `string` | Private API-server FQDN, populated for private clusters |
| `portal_fqdn` | `string` | FQDN the Azure Portal uses to reach the cluster |
| `oidc_issuer_url` | `string` | OIDC issuer URL -- consumed by `AzureFederatedIdentityCredential` for workload identity |
| `node_resource_group` | `string` | Name of the Azure-managed node resource group |
| `node_resource_group_id` | `string` | ARM ID of the node resource group |
| `cluster_kubeconfig` | `string` | Base64-encoded kubeconfig (treat as a secret) |
| `cluster_identity_principal_id` | `string` | Principal ID of the cluster's managed identity |
| `kubelet_identity_object_id` | `string` | Kubelet identity object ID -- grant it AcrPull on registries |
| `kubelet_identity_client_id` | `string` | Kubelet identity client ID |
| `current_kubernetes_version` | `string` | Kubernetes version the control plane is actually running |

## Related Components

- [AzureAksNodePool](/docs/catalog/azure/aks-node-pool) — additional workload-shaped pools (general, spot, GPU, Windows) with independent lifecycles
- [AzureSubnet](/docs/catalog/azure/subnet) — optional BYO-network placement for the default pool
- [AzureFederatedIdentityCredential](/docs/catalog/azure/federated-identity-credential) — binds workload identity to the cluster's OIDC issuer
- [AzurePrivateDnsZone](/docs/catalog/azure/private-dns-zone) — BYO private DNS zone for private clusters
