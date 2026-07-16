---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cosmos DB SQL Role Assignment"
type: "preset-list"
componentSlug: "cosmos-db-sql-role-assignment"
componentTitle: "Cosmos DB SQL Role Assignment"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-workload-data-contributor"
    rank: "01"
    title: "Workload Data Contributor"
    excerpt: "This preset is the standard keyless application grant: the built-in Data Contributor role bound to a workload's managed identity across the whole account. The app authenticates as its identity (no..."
  - slug: "02-container-scoped-reader"
    rank: "02"
    title: "Container-Scoped Reader"
    excerpt: "This preset is the least-privilege read grant: the built-in Data Reader bound to a principal on exactly ONE container. An analytics job, a dashboard, or a debugging operator sees the documents it..."
  - slug: "03-custom-role-grant"
    rank: "03"
    title: "Custom Role Grant"
    excerpt: "This preset completes the custom-RBAC composition: an AzureCosmosdbSqlRoleDefinition's fully-scoped ID flows into the grant by reference, and the definition, the grant, the identity, and the account..."
---

# Cosmos DB SQL Role Assignment Presets

Ready-to-deploy configuration presets for Cosmos DB SQL Role Assignment. Each preset is a complete manifest you can copy, customize, and deploy.
