---
title: "Relationship Tuple"
description: "Relationship Tuple deployment documentation"
icon: "package"
order: 100
componentName: "openfgarelationshiptuple"
---

# OpenFGA Relationship Tuple

Deploys a relationship tuple into an existing OpenFGA store. A relationship tuple is the fundamental unit of authorization data, representing that a specific user (or userset) has a particular relation to an object. Together with an authorization model, tuples determine access decisions at check time. All tuple fields are immutable -- changing any field replaces the tuple. Integrates with Planton's Provider Connections for OpenFGA credential management and ValueFromRef for store and model dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Relationship Tuple** -- an `openfga_relationship_tuple` resource that writes a single authorization tuple (user, relation, object) into the specified OpenFGA store

## Before You Deploy

### Planton Setup

- **OpenFGA Provider Connection** -- an active connection in the Connect module with the OpenFGA API URL and authentication credentials (API token or client credentials). Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline authentication.

### OpenFGA Server

- **A running OpenFGA instance** -- self-hosted or cloud-hosted, reachable from the Planton Runner or provisioner environment.
- **An existing OpenFGA store** with an authorization model that defines the types and relations used in the tuple.

## Deploy

### Console

Open the deployment store, find **OpenFGA Relationship Tuple**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **User-Document Access** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openfga.planton.dev/v1
kind: OpenFgaRelationshipTuple
metadata:
  name: anne-views-budget
  org: acme-corp
  env: prod
spec:
  storeId:
    value: "01HXYZ..."
  user:
    type: user
    id:
      value: anne
  relation: viewer
  object:
    type: document
    id:
      value: budget-2024
```

```shell
planton apply -f openfga-tuple.yaml
```

This creates a relationship tuple granting user `anne` the `viewer` relation on `document:budget-2024`. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the tuple to a store and model deployed in the same InfraPipeline:

```yaml
spec:
  storeId:
    valueFrom:
      kind: OpenFgaStore
      name: prod-authz
      fieldPath: status.outputs.id
  authorizationModelId:
    valueFrom:
      kind: OpenFgaAuthorizationModel
      name: rbac-model
      fieldPath: status.outputs.id
```

The InfraPipeline resolves the dependency graph, deploys the store and model first, then writes the relationship tuple.

## Key Configuration

These are the most important decisions when configuring an OpenFGA relationship tuple. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**User structure** -- The `user` object contains `type`, `id`, and an optional `relation` for usersets. Simple users resolve to `type:id` (e.g., `user:anne`). Adding a `relation` field creates a userset reference like `group:engineering#member`, granting access to all members of that group.

**Relation** -- The `relation` field must match a relation defined in the authorization model for the object type. Common values include `viewer`, `editor`, `owner`, `member`, and `admin`.

**Object structure** -- The `object` contains `type` and `id`, resolving to `type:id` (e.g., `document:budget-2024`). The type must be defined in the authorization model.

**Authorization model pinning** -- The optional `authorizationModelId` field pins the tuple to a specific model version for validation. When omitted, the tuple is associated with the latest model in the store. Pinning is useful in production to ensure tuples are validated against a known model version.

**Conditions** -- The optional `condition` block adds dynamic access rules evaluated at check time. Specify a `name` matching a condition defined in the authorization model and an optional `contextJson` with partial context that is merged with runtime context during evaluation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenFgaStore** | `storeId` | `status.outputs.id` |
| **OpenFgaAuthorizationModel** | `authorizationModelId` | `status.outputs.id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that confirm the tuple was created:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `user` | The subject of the tuple in `type:id` or `type:id#relation` format | Audit logs, permission verification |
| `relation` | The relationship type that was created | Audit logs, permission verification |
| `object` | The resource the tuple grants access to in `type:id` format | Audit logs, permission verification |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Direct user access** -- Grants a specific user a relation on a specific resource (e.g., `user:anne` is a `viewer` of `document:budget-2024`). This is the most fundamental authorization pattern in OpenFGA. Start from the **User-Document Access** preset.

**Group membership** -- Adds a user to a group by creating a tuple like `user:anne` is a `member` of `group:engineering`. Other tuples can then grant access to `group:engineering#member`, providing inherited permissions to all group members. Start from the **Group Membership** preset.

**Conditional access** -- Attaches a condition to a tuple so that the relation is only active when runtime context satisfies the condition (e.g., IP range restrictions, time-of-day checks). The condition must be defined in the authorization model.

## Works With

- [**OpenFGA Store**](/cloud-catalog/openfga-store) -- provides the store where relationship tuples are written
- [**OpenFGA Authorization Model**](/cloud-catalog/openfga-authorization-model) -- defines the types, relations, and conditions that govern how tuples are evaluated
