# OpenFGA Authorization Model

Deploys an authorization model into an existing OpenFGA store. The model defines types, relations, and access rules that govern how permission checks are evaluated. Models can be specified in DSL format (recommended) or JSON format, and are immutable -- each update creates a new model version with a new ID. Integrates with Planton's Provider Connections for OpenFGA credential management and ValueFromRef for store dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Authorization Model Document** -- when `modelDsl` is provided, a data source converts the human-readable DSL to JSON for the OpenFGA API
- **Authorization Model** -- an `openfga_authorization_model` resource containing the type definitions, relations, and conditions for fine-grained access control

## Before You Deploy

### Planton Setup

- **OpenFGA Provider Connection** -- an active connection in the Connect module with the OpenFGA API URL and authentication credentials (API token or client credentials). Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline authentication.

### OpenFGA Server

- **A running OpenFGA instance** -- self-hosted or cloud-hosted, reachable from the Planton Runner or provisioner environment.
- **An existing OpenFGA store** -- provide the store ID directly or reference an OpenFgaStore Cloud Resource via ValueFromRef.

## Deploy

### Console

Open the deployment store, find **OpenFGA Authorization Model**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **RBAC Model (DSL)** preset in the [Presets](#presets) tab for a common role-based access control pattern.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openfga.planton.dev/v1
kind: OpenFgaAuthorizationModel
metadata:
  name: rbac-model
  org: acme-corp
  env: prod
spec:
  storeId:
    value: "01HXYZ..."
  modelDsl: |
    model
      schema 1.1

    type user

    type document
      relations
        define viewer: [user]
        define editor: [user]
        define owner: [user]
```

```shell
planton apply -f openfga-authz-model.yaml
```

This creates an authorization model with a `user` type and a `document` type that has viewer, editor, and owner relations. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the model to a store deployed in the same InfraPipeline:

```yaml
spec:
  storeId:
    valueFrom:
      kind: OpenFgaStore
      name: prod-authz
      fieldPath: status.outputs.id
```

The InfraPipeline resolves the dependency graph, deploys the store first, then creates the authorization model in it.

## Key Configuration

These are the most important decisions when configuring an OpenFGA authorization model. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Model format** -- Choose between `modelDsl` (recommended) and `modelJson`. The DSL format is more human-readable and is the format used in OpenFGA documentation. The IaC module automatically converts DSL to JSON during deployment. Exactly one of the two must be specified.

**Store ID** -- The `storeId` field identifies which OpenFGA store holds this model. Provide a direct value or reference an OpenFgaStore via ValueFromRef. The store ID is immutable -- changing it requires replacing the model.

**Model immutability** -- Authorization models are immutable in OpenFGA. Each change creates a new model version with a new ID. Existing relationship tuples remain valid across model versions as long as the types and relations they reference are still defined.

**Type design** -- Define types that map to your application's entities (e.g., `user`, `group`, `document`, `folder`). Each type can have relations that control access (e.g., `viewer`, `editor`, `owner`). Use usersets like `group#member` to enable group-based inheritance.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenFgaStore** | `storeId` | `status.outputs.id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `id` | Unique identifier of the authorization model version | Referenced by OpenFGA Relationship Tuple components via `authorizationModelId` |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**RBAC model** -- Defines user, group, and document types with viewer/editor/owner roles and group-based access inheritance. Groups use a `member` relation, and documents grant access to `group#member` usersets. This is the most common starting pattern for applications that need role-based access control. Start from the **RBAC Model (DSL)** preset.

**Hierarchical document access** -- Adds folder and document types with parent relationships so that viewer and editor permissions cascade from folders to their contents. Suitable for file-management and content-management systems with nested permission hierarchies. Start from the **Hierarchical Document Access Model (DSL)** preset.

## Works With

- [**OpenFGA Store**](/cloud-catalog/openfga-store) -- provides the store where authorization models are created
