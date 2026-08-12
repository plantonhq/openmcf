# AzureMongoClusterUser Terraform Module

## Overview

Grants a Microsoft Entra ID principal access to an Azure Cosmos DB for MongoDB vCore cluster -- an access binding, not a password user.

## Resources Created

- `azurerm_mongo_cluster_user` -- the grant (principal binding plus database role grants)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMongoClusterUserSpec fields; the cluster and principal references arrive as resolved literals

## Outputs

- `mongo_cluster_user_id` -- the grant's ARM resource ID (`{cluster_id}/users/{object_id}`)
- `mongo_cluster_user_name` -- the grant's ARM name (the granted principal's object id)

## Behavior Notes

- **"MicrosoftEntraID" is the identity provider's only legal value today** -- deliberately not part of the spec; the module sends it explicitly.
- **Everything is create-only**: the resource has no update path -- changing anything replaces the grant (a harmless drop-and-re-add of an access binding; add-then-remove for zero-gap role changes).
- **Deploy-time contract**: the target cluster must list "MicrosoftEntraID" in its `authentication_methods` -- Azure rejects the grant otherwise.
- **Principal types are lowercase tokens** (`user`, `servicePrincipal`) -- Azure's own casing; managed identities bind through their service principal.
