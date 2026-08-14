# AzureMongoCluster Terraform Module

## Overview

Creates an Azure Cosmos DB for MongoDB vCore cluster -- a real MongoDB engine on dedicated vCore tiers -- plus one firewall-rule resource per named client-IP range.

## Resources Created

- `azurerm_mongo_cluster` -- the cluster (mode, sizing, auth, identity, CMK, network posture)
- `azurerm_mongo_cluster_firewall_rule` -- one per `firewall_rules` entry via `for_each` keyed by the rule's name (renames replace only that rule)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMongoClusterSpec fields; the resource group, password, CMK, identity, and source-cluster references arrive as resolved literals

## Outputs

- `mongo_cluster_id` -- the cluster's ARM resource ID (the target an AzureMongoClusterUser's `mongo_cluster_id` references)
- `mongo_cluster_name` -- the cluster's name (the first label of its hostname)
- `connection_string` (sensitive) -- the primary MongoDB URI, administrator credentials substituted in
- `connection_strings` (sensitive) -- every published connection string, keyed by Azure's name for it

## Behavior Notes

- **The provider owns the mode machinery** this module deliberately does not reimplement: Default mode stages the Data API in a separate post-create update, upgrades away from Free/M25 stage a tier-first update, and a `create_mode` change is forced to a replacement (Azure never returns the mode on reads).
- **Identity add/remove is a replacement** -- Azure rejects the in-place transition between no-identity and user-assigned; changing the SET of identities updates in place. User-assigned is the only flavor the service supports.
- **Sizing fields are sent only when set** -- replica and restore modes inherit them from the source; the storage type rides the storage size (platform default "PremiumSSD" when the size is set).
- **`data_api_mode_enabled` is sent only when the manifest carries it** -- the provider errors when the raw config sets it (even false) on a non-Default-mode cluster.
- **`authentication_methods` is sent only when set** -- Azure defaults an unset list to ["NativeAuth"] server-side.
- **Platform defaults always sent**: create_mode "Default" (`coalesce`), public network access enabled.
