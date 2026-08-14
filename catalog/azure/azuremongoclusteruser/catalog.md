# Azure Mongo Cluster User

Grants a Microsoft Entra ID principal access to an Azure Cosmos DB for MongoDB vCore cluster -- an access binding, not a password user. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Mongo cluster user grant** -- the Entra principal's binding on the cluster, with its database role grants

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **Mongo Cluster** -- reference an AzureMongoCluster's ID output or provide an existing cluster's ARM ID.
- **The principal** -- for workload identities, reference the AzureUserAssignedIdentity's principal ID output; for humans or app registrations, the object ID from Entra.

### Azure Subscription

- **The cluster must allow Entra authentication** -- "MicrosoftEntraID" must be in the cluster's `authenticationMethods`; Azure rejects the grant at deploy time otherwise.
- **Everything is create-only** -- the resource has no update path; changing anything replaces the grant (a harmless drop-and-re-add of an access binding).
- **The role vocabulary is "root" today** -- scoped to a database; grant it on "admin" for cluster-wide access. Azure will widen the vocabulary over time.
- **A grant is free** -- it is a data-plane user entry, not billable infrastructure.

## Deploy

### Console

Open the deployment store, find **Azure Mongo Cluster User**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **App Identity Access** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f mongo-cluster-user.yaml
```

## After Deploy

The application connects with MongoDB's `MONGODB-OIDC` authentication mechanism using its Entra token -- no password anywhere. Verify access with `mongosh` using `--authenticationMechanism MONGODB-OIDC` under the granted identity; revoke by deleting this resource.
