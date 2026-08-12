# Overview

The **AzureMongoCluster** component deploys an Azure Cosmos DB for MongoDB vCore cluster -- Azure's modern managed MongoDB offering. It runs a real MongoDB engine (community drivers connect unchanged) on dedicated vCore-based compute, with vertical tiers from a free sandbox to M200, optional sharding, zone-redundant high availability, geo replicas, and point-in-time restore.

## Purpose

- **MongoDB without the ops**: Azure owns patching, backups (continuous, 35 days), and failover; you own the sizing dials.
- **Three ways to exist**: `create_mode` builds a fresh cluster ("Default"), a cross-region read replica ("GeoReplica"), or a point-in-time clone ("PointInTimeRestore").
- **Two auth worlds**: native username/password administration lives here; Entra-principal access grants are separate AzureMongoClusterUser resources so app teams onboard themselves.

## Key Features

- Full azurerm v5 surface: the creation-mode matrix, compute tiers (Free through M200), storage size/type, sharding, high availability, authentication methods, user-assigned identity, customer-managed-key encryption, preview features, Data API, network posture, and composed firewall rules.
- The provider's whole creation-mode contract front-loaded as validation: Default mode demands its six sizing fields, GeoReplica its source coordinates, PointInTimeRestore its restore block -- rejected at manifest time, not at deploy time.
- Chart-ready: `resource_group` defaults to AzureResourceGroup, the CMK key to AzureKeyVaultKey's versionless ID, identities to AzureUserAssignedIdentity, and replica/restore sources to AzureMongoCluster itself; the sensitive `connection_string` output is the app-wiring edge.

## Use Cases

- **Application database**: a Default-mode M30+ cluster with zone-redundant HA and Entra-based app access.
- **Read scale-out across regions**: a GeoReplica cluster following a source that has the GeoReplicas preview feature enabled.
- **Incident recovery**: a PointInTimeRestore clone from minutes before the bad deploy, cut over after verification.

## Future Enhancements

- Entra access grants live in AzureMongoClusterUser -- point its `mongo_cluster_id` at this component's ID output.
