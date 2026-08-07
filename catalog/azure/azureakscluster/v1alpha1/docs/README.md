# AzureAksCluster -- Design Research

## The Resource

An AKS managed cluster (`Microsoft.ContainerService/managedClusters`) is
Azure's managed Kubernetes control plane plus the one node pool ARM requires
at creation. The component maps 1:1 onto `azurerm_kubernetes_cluster`
(azurerm v4.80, `internal/services/containers/kubernetes_cluster_resource.go`),
parity-verified against pulumi-azure v6.38
(`containerservice.KubernetesCluster`).

## Shape Decisions (the ones that define the component)

- **Exactly one node pool lives on the cluster: the mandatory
  `default_node_pool`.** azurerm itself carries only the default pool on
  the cluster resource; additional pools are standalone
  `azurerm_kubernetes_cluster_node_pool` resources. The spec mirrors that
  grain exactly: the previous invented `system_node_pool`/`user_node_pools`
  split is gone, and every additional pool is an `AzureAksNodePool`
  referencing this cluster's `cluster_id` output. Pools get independent
  lifecycles; cluster updates are never forced by pool changes.
- **The default pool's sub-message shapes are converged with the standalone
  kind** (kubelet config, Linux OS/sysctl config, node network profile,
  upgrade settings), minus the fields azurerm itself excludes on the
  default pool (priority/spot, os_type, mode, taints). Moving a workload
  pool out to its own resource is a mechanical copy.
- **Subnet placement lives on the pool and is optional.** In azurerm,
  `vnet_subnet_id`/`pod_subnet_id` are node-pool fields, and unset means
  AKS provisions and manages its own network -- Azure's actual default.
  The old cluster-level required `vnet_subnet_id` was an invented shape
  that forced every cluster to bring a network. Registry prerequisites
  shrink to `[AzureResourceGroup]`; BYO-subnet is the composed advanced
  path (the cluster identity then needs Network Contributor on it).
- **`kubernetes_version` has no baked-in default.** A hardcoded
  `recommended_default` goes stale by design (verified live: the test
  region's AKS default had moved two minors past the old pin). Unset means
  azurerm's own contract -- AKS provisions the latest recommended GA
  version; the field comment teaches explicit pinning for production.
- **`oidc_issuer_enabled` defaults ON** (deliberately above Azure's
  provisioning default): it is the trust anchor for workload identity
  federation (`AzureFederatedIdentityCredential` consumes the
  `oidc_issuer_url` output), costs nothing, and enabling it later forces
  no replacement while disabling it after use does.
- **The network fabric is written explicitly even when unset.** azurerm's
  implicit fallback is kubenet -- deprecated, retiring 2028. Both modules
  make the modern AKS default (Azure CNI overlay) the actual default.

## Deliberately Not Modeled (recorded reasons)

- **`service_principal` block.** Managed identity is Azure's own stated
  direction; a client secret in cluster config is exactly the credential
  class the platform exists to eliminate. The `identity` block (system- or
  user-assigned) is the only auth model.
- **`http_application_routing_enabled`.** Retired by Azure in favor of
  web-app routing, which IS modeled (`web_app_routing`).
- **Deprecated 4.x kubelet/OS fields** (`container_log_max_line`,
  `transparent_huge_page_enabled`): azurerm carries them only for
  compatibility; the replacement fields (`container_log_max_size_mb`,
  `transparent_huge_page`) are modeled.

## Design Decisions

- **`key_management_service.key_vault_key_id` is a plain string**, not a
  StringValueOrRef: no `AzureKeyVaultKey` kind exists yet (committed for
  the data wave). This field is the FK seam when it lands.
- **DNS prefix defaulting:** ARM requires exactly one prefix flavor. When
  the spec sets neither, the modules derive the public prefix from the
  cluster name -- unless the cluster is private and carries its own
  private prefix.
- **Secret-by-default:** `windows_profile.admin_password` and
  `http_proxy_config.trusted_ca` are `(sensitive)`. The kubeconfig output
  follows the catalog's outputs grain (no output-side sensitivity
  annotations exist anywhere in the catalog); the Pulumi module wraps it
  in `pulumi.ToSecret`, the TF output is marked `sensitive`.

## Operational Behavior Worth Knowing

- The network fabric (plugin/mode/CIDRs/outbound type) mostly replaces the
  cluster when changed -- the network model is a day-0 decision.
- Many default-pool shape changes (vm_size, os_disk_type, fips, host
  encryption) rotate the pool; `temporary_name_for_rotation` lets AKS
  stand up a replacement pool first. Set it proactively in production.
- Upgrades move one minor at a time; node pools may lag the control plane
  by up to two minors. `automatic_upgrade_channel` + the maintenance
  windows govern when AKS acts on its own.
- Private clusters resolve through a private DNS zone; BYO zone
  (`private_dns_zone_id`) requires the cluster identity to hold Private
  DNS Zone Contributor + Network Contributor before creation.

## Composition

- `resource_group` → `AzureResourceGroup.status.outputs.resource_group_name`
- `default_node_pool.vnet_subnet_id` / `pod_subnet_id` → `AzureSubnet.subnet_id`
- `private_dns_zone_id` → `AzurePrivateDnsZone.zone_id`
- `identity.identity_ids` / `kubelet_identity.user_assigned_identity_id` →
  `AzureUserAssignedIdentity`
- `oms_agent` / `microsoft_defender` / `monitor_metrics` workspace ids →
  `AzureLogAnalyticsWorkspace.workspace_id`
- `bootstrap_profile.container_registry_id` → `AzureContainerRegistry`
- `ingress_application_gateway.gateway_id` → `AzureApplicationGateway`
- Consumed by: `AzureAksNodePool.kubernetes_cluster_id` (→
  `status.outputs.cluster_id`) and
  `AzureFederatedIdentityCredential.issuer` (→
  `status.outputs.oidc_issuer_url`)
