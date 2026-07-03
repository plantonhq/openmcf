# Azure AKS Cluster and Node Pool Depth Rework

**Date**: July 4, 2026
**Type**: Feature + Breaking Change
**Components**: Azure Provider, API Definitions, IAC Modules, E2E

## Summary

The compute wave opens with a paired rework of `AzureAksCluster` and
`AzureAksNodePool` to the full `azurerm` v4.80 surface. The cluster sheds inline
user node pools and optionalizes subnet placement on the default pool; the
standalone node pool converges to ~49 fields with a single `kubernetes_cluster_id`
parent foreign key. Both kinds migrate off engine outliers: Terraform to
`azurerm ~> 4.0`, Pulumi to classic `pulumi-azure/sdk/v6` via the shared
`pulumiazureprovider.Get` builder.

## Breaking Changes

- **`AzureAksCluster`**: `system_node_pool` → `default_node_pool`; `user_node_pools`
  removed; cluster-level `vnet_subnet_id` removed (optional on default pool);
  `control_plane_sku` → `sku_tier`; outputs renamed/rebuilt (`cluster_id`,
  `oidc_issuer_url`, etc.).
- **`AzureAksNodePool`**: `cluster_name` + `resource_group` → `kubernetes_cluster_id`;
  `initial_node_count` → `node_count`; `spot_enabled` → `priority`/`eviction_policy`/
  `spot_max_price`; outputs renamed (`node_pool_id`, `node_image_version`).

## Validation

Live dual-engine E2E passed for both kinds (Pulumi + Terraform). Offline gate:
spec tests, secret-coverage, validate-refs, outputs conformance, tofu plan.

## Impact

Users provisioning AKS through Planton get the full Azure-authentic configuration
model, managed-networking-by-default clusters, and composable node pools with
independent lifecycles. Workload identity composition closes via `oidc_issuer_url`.
