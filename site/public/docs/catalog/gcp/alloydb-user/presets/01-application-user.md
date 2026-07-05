---
title: "Application User (ALLOYDB_BUILT_IN)"
description: "This preset creates a classic username/password AlloyDB user for an application service."
type: "preset"
rank: "01"
presetSlug: "01-application-user"
componentSlug: "alloydb-user"
componentTitle: "AlloyDB User"
provider: "gcp"
icon: "package"
order: 1
---

# Application User (ALLOYDB_BUILT_IN)

This preset creates a classic username/password AlloyDB user for an application service.

## When to Use

- Applications that connect with a stored credential (outside IAM proxy flows)
- One user per service with its own rotatable password

## Key Configuration Choices

- **ALLOYDB_BUILT_IN (default)** — password-authenticated database role
- **databaseRoles: [alloydbiamuser]** — standard application role; adjust for your privilege model

## Related Presets

- **02-iam-user** — passwordless IAM-authenticated user

## Related Components

- [GcpAlloydbCluster](/docs/catalog/gcp/gcpalloydbcluster) — the cluster this user lives on
