# Azure Role Definition

Deploys a custom Azure RBAC role: a named, reusable bundle of permissions that principals can then be granted through role assignments. Azure ships hundreds of built-in roles, but real organizations routinely need permission sets the built-ins don't express — "Contributor, except role assignments and policy writes", "can restart VMs but not create or delete them", "read-only plus blob data access". A custom role captures such a set once, with a meaningful name, and every grant of it stays consistent as the definition evolves. The component integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring to resource-group scopes.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Custom Role Definition** -- one tenant-visible RBAC role at the chosen scope (management group, subscription, or resource group), with its permission blocks and assignable scopes

A definition grants nothing by itself -- permissions only take effect when an **AzureRoleAssignment** binds the role to a principal at a scope. The definition's `role_definition_id` output is exactly what an assignment's `roleDefinitionId` field consumes.

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### Azure Subscription

- **Authorization-plane rights on the target scope** -- creating role definitions requires `Microsoft.Authorization/roleDefinitions/write`, held via Owner, User Access Administrator, or Role Based Access Control Administrator. A Contributor-level deploying credential fails with "AuthorizationFailed".

## Deploy

### Console

Open the deployment store, find **Azure Role Definition**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields -- with quick-add templates for the classic permission-block shapes. Start from the **VM Operator Role** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleDefinition
metadata:
  name: acme-vm-operator
  org: acme-corp
  env: prod
spec:
  name: acme-vm-operator
  scope:
    value: "/subscriptions/00000000-0000-0000-0000-000000000000"
  description: "Operate existing VMs: start/stop/restart, no create or delete"
  permissions:
    - actions:
        - Microsoft.Compute/virtualMachines/read
        - Microsoft.Compute/virtualMachines/start/action
        - Microsoft.Compute/virtualMachines/restart/action
        - Microsoft.Compute/virtualMachines/powerOff/action
        - Microsoft.Compute/virtualMachines/deallocate/action
```

```shell
planton apply -f role.yaml
```

Assignable scopes are omitted, so Azure defaults them to the definition's own scope (the subscription).

### InfraChart

Compose the definition with the assignments that grant it -- the assignment consumes the definition's output:

```yaml
# The custom role
spec:
  name: acme-blob-reader
  scope:
    valueFrom:
      kind: AzureResourceGroup
      name: data-rg
      fieldPath: status.outputs.resource_group_id
  permissions:
    - dataActions:
        - Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read
---
# The grant (separate AzureRoleAssignment resource)
spec:
  scope:
    valueFrom:
      kind: AzureResourceGroup
      name: data-rg
      fieldPath: status.outputs.resource_group_id
  roleDefinitionId: "<from the definition's role_definition_id output>"
  principalId:
    valueFrom:
      kind: AzureUserAssignedIdentity
      name: analytics-reader
      fieldPath: status.outputs.principal_id
```

## Key Configuration

These are the most important decisions when configuring a custom role. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Control plane vs data plane** -- THE distinction to get right. `actions` govern managing resources through ARM ("restart a VM"); `dataActions` govern the data inside them ("read this blob", "read a secret's value in RBAC mode"). A role with `Microsoft.Storage/storageAccounts/*` actions can manage the account yet cannot read one byte of blob data.

**Carve-outs are trims, not denies** -- `notActions`/`notDataActions` subtract from THIS role's own grant; they cannot take away permissions another assignment gives the same principal (that is Azure deny assignments' job). The classic shape: actions `["*"]`, notActions `["Microsoft.Authorization/*/write"]` -- everything except changing RBAC.

**Scope and assignability** -- The definition's scope anchors its ARM ID (ForceNew) and defaults where it can be assigned. Choose the highest scope the role will ever be assigned at; narrow `assignableScopes` below it to pre-authorize where the role may ever be granted. At most one management group may appear in the list.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AzureResourceGroup** | `scope`, `assignableScopes` | `status.outputs.resource_group_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `role_definition_id` | Fully-scoped ARM ID of the definition | AzureRoleAssignment's `roleDefinitionId` field |
| `role_definition_guid` | The definition's GUID resource name | Authorization API automation |
| `role_name` | The display name as deployed | Portal cross-reference |
| `scope` | The creation scope as resolved | Auditing |
| `assignable_scopes` | The scopes Azure recorded (incl. the provider-defaulted own scope) | Auditing where grants may exist |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Operational roles** -- Verb-scoped control-plane roles ("operate, don't create"): the VM operator pattern. Start from the **VM Operator Role** preset.

**Data-access roles** -- Pure data-plane grants with zero management surface: the blob-data-reader pattern. Start from the **Blob Data Reader Role** preset.

**Admin-minus-carve-out** -- Broad grants with the RBAC-write carve-out, keeping permission changes in the hands of the few. Start from the **Project Admin Carve-Out Role** preset.

## Works With

- [**Azure Role Assignment**](/cloud-catalog/azure-role-assignment) -- binds this role to a principal at a scope, consuming `role_definition_id`
- [**Azure Resource Group**](/cloud-catalog/azure-resource-group) -- a common definition scope and assignable-scope boundary
- [**Azure User Assigned Identity**](/cloud-catalog/azure-user-assigned-identity) -- the principal most grants of custom roles target
