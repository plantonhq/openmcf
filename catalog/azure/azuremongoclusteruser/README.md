# Overview

The **AzureMongoClusterUser** component grants a Microsoft Entra ID principal access to an Azure Cosmos DB for MongoDB vCore cluster (AzureMongoCluster). It is an ACCESS GRANT, not a password user: the principal -- a person, a service principal, or a managed identity -- authenticates to MongoDB with its Entra token and receives the granted database roles. Native username/password administration lives on the cluster itself.

## Purpose

- **Passwordless application access**: an app's managed identity connects to MongoDB with a token; nothing to rotate, nothing to leak.
- **Self-service onboarding**: many grants share one cluster with independent lifecycles -- app teams bind their own principals without touching the cluster.
- **Clean revocation**: deleting the grant severs access; the cluster and its data are untouched.

## Key Features

- Full azurerm v5 surface: the principal's object ID, its type (user or service principal), and the role grants, modeled exactly.
- Azure pins the identity provider to "MicrosoftEntraID" (the only value the service accepts today) -- deliberately not part of the spec; both engines send it explicitly.
- Chart-ready: `mongo_cluster_id` defaults its reference to AzureMongoCluster's ID output; `object_id` takes any principal (reference a managed identity's principal ID or pass a literal UUID).

## Use Cases

- **App-to-database wiring**: grant the workload identity root on its database; the app connects with `MONGODB-OIDC` authentication.
- **Human break-glass access**: grant an operator's Entra user for incident work, delete the grant after.
- **Tenant data-plane isolation**: one grant per tenant service principal, each scoped to its own database.

## Future Enhancements

- Azure's role vocabulary is "root" today; finer-grained built-in roles widen the `roles` field when the service ships them.
