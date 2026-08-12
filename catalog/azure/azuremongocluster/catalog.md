# Azure Mongo Cluster (Cosmos DB for MongoDB vCore)

Deploys an Azure Cosmos DB for MongoDB vCore cluster -- a real MongoDB engine on dedicated vCore tiers with sharding, zone-redundant HA, geo replicas, and point-in-time restore. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Mongo vCore cluster** -- the cluster itself, in one of three creation modes (fresh, geo replica, or point-in-time clone)
- **Firewall rules** -- one per named client-IP range in the spec

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Resource Group** -- reference an AzureResourceGroup or provide an existing group's name.
- **For CMK encryption (optional)** -- an AzureKeyVaultKey (referenced by its VERSIONLESS ID) and an AzureUserAssignedIdentity with unwrap/wrap permissions on the vault, granted BEFORE the cluster is created.

### Azure Subscription

- **The cluster name is a global hostname** ({name}.mongocluster.cosmos.azure.com) -- a taken name fails at deploy time; prefix with your org.
- **One Free-tier cluster per subscription** -- and Free/M25 tiers cannot shard beyond one shard or use zone-redundant HA.
- **create_mode is invisible to Azure reads** -- changing it always replaces the cluster; pick the mode deliberately.
- **Identity transitions replace the cluster** -- adding the first user-assigned identity or removing the last one is a replacement, not an update.
- **Continuous backup is automatic** (35 days) -- point-in-time restore needs no configuration on the source.

## Deploy

### Console

Open the deployment store, find **Azure Mongo Cluster**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production Cluster** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f mongo-cluster.yaml
```

## After Deploy

The sensitive `connection_string` output carries the primary MongoDB URI with the administrator credentials substituted in -- wire applications to it by reference, never by copy. For Entra-based application access, create AzureMongoClusterUser grants against the cluster's ID output (the cluster must list "MicrosoftEntraID" in `authenticationMethods`). Watch storage growth on the **Metrics** blade; storage grows in place, but never shrinks.
