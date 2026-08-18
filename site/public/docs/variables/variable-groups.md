---
title: "Variable Groups"
description: "Collect related variables into named sets for managing application configuration as a cohesive unit"
icon: package
order: 20
tags:
  - Variables
  - Variable Groups
  - Configuration
---

# Variable Groups

A variable group is a named collection of related variables. Instead of managing dozens of individual variables, you group them by purpose — so you can see, update, and assign an application's configuration as a cohesive unit.

## Why Variable Groups

Consider a service that needs database configuration:

Without groups, you manage five separate variables: `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_SSL_MODE`. Each is created, updated, and tracked independently.

With a variable group called `database-config`, all five variables live together. You can see the complete database configuration at a glance, update it as a unit, and assign it to services as a single reference.

Variable groups promote:

- **Configuration as a unit.** Related values stay together. When you review the `database-config` group, you see everything a service needs to connect to the database.
- **Reuse across services.** Multiple services can reference the same variable group. Update the group, and every service picks up the change.
- **Consistency.** When two services need the same database connection, they reference the same group — no risk of one having a stale hostname.

## Scoping

Like individual variables, variable groups are scoped to the organization or a specific environment:

- **Organization groups** — Shared configuration available across all environments (e.g., a `monitoring-config` group with endpoints that are the same everywhere).
- **Environment groups** — Environment-specific configuration (e.g., a `database-config` group with production-specific endpoints and credentials).

## Managing Variable Groups

### Using the Web Console

Variable groups are managed from the Variables section of the web console. You can create groups, add or remove entries, and edit values individually or in bulk.

<!-- SCREENSHOT: Variable group management
  Page: /orgs/{org}/variables (group view)
  Action: Show a variable group with multiple entries
  Focus: The group entries table showing names, values, and actions
  Alt: Variable group showing multiple configuration entries with names, descriptions, values, and action buttons
-->

### Using the CLI

```bash
# List all variables (shows group membership)
planton service variables

# Get a specific variable from a group
planton service variables get-value --group database-config --name DB_HOST

# Add or update an entry in a group
planton service variables upsert

# Delete an entry from a group
planton service variables delete --group database-config --name DB_HOST
```

## Entries with Dynamic References

Each entry in a variable group can independently be a literal value or a dynamic reference. This means you can mix static and dynamic values in the same group:

```yaml
group: database-config
entries:
  - name: DB_HOST
    description: "Resolved from deployed Postgres instance"
    value_from:
      kind: PostgresCluster
      env: production
      name: my-postgres
      field_path: "status.endpoint"
  - name: DB_PORT
    value: "5432"
  - name: DB_NAME
    value: "myapp"
  - name: DB_SSL_MODE
    value: "require"
```

The `DB_HOST` entry resolves dynamically from infrastructure. The other entries are static values. When the group is resolved, all entries — dynamic and static — are returned together.

## Related Documentation

- [Variables](/docs/variables) — Individual variables, scoping, and dynamic references
- [Secrets](/docs/secrets) — Sensitive value management with encryption
- [What is a Service?](/docs/ci-cd/what-is-a-service) — How services consume variables and secrets
