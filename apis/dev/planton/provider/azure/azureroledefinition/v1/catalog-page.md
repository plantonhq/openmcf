# Azure Role Definition

Creates a custom Azure RBAC role -- a named, reusable bundle of permissions that role assignments then grant to principals. Use it when the built-in roles don't express what a team or workload should be allowed to do, and manage your organization's role catalog as version-controlled infrastructure.

## What Gets Created

When you deploy an AzureRoleDefinition resource, Planton provisions:

- **Role Definition** — an `azurerm_role_definition` holding the role's display name, permission blocks (control-plane and data-plane), and assignable scopes at the target ARM scope

Role definitions carry no Azure tags -- the `Microsoft.Authorization` resource type does not support them.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **Authorization-plane rights**: the deploying credential needs `Microsoft.Authorization/roleDefinitions/write` at the target scope (Owner, User Access Administrator, or Role Based Access Control Administrator -- Contributor alone is not sufficient)
- **A tenant-unique role name** (Azure rejects duplicates; prefer org-prefixed names)

## Quick Start

Create a file `roledefinition.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleDefinition
metadata:
  name: vm-operator
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRoleDefinition.vm-operator
spec:
  name: acme-vm-operator
  scope:
    value: /subscriptions/00000000-0000-0000-0000-000000000000
  description: Operate existing VMs (start/stop/restart) without create or delete rights
  permissions:
    - actions:
        - Microsoft.Compute/virtualMachines/read
        - Microsoft.Compute/virtualMachines/start/action
        - Microsoft.Compute/virtualMachines/restart/action
        - Microsoft.Compute/virtualMachines/deallocate/action
```

Deploy:

```shell
planton apply -f roledefinition.yaml
```

This creates a subscription-scoped custom role that any `AzureRoleAssignment` can then grant by its definition ID.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `name` | `string` | The role's display name, unique within the Azure AD tenant. | Required, non-empty |
| `scope` | `StringValueOrRef` | The ARM scope the definition is created at: management group, subscription, or resource group. Defaults to referencing an `AzureResourceGroup`'s ARM ID; subscription/management-group IDs pass as literals. | Required |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | `string` | Intent text shown beside the role in the portal's role picker. |
| `permissions` | `list` | Permission blocks, each with `actions`, `notActions`, `dataActions`, `notDataActions` (ARM operation patterns, `*` wildcards allowed). Azure unions multiple blocks; one block is the norm. |
| `assignableScopes` | `list of StringValueOrRef` | Where the role may be assigned. Defaults to the definition's own scope. At most one management group may appear. |
| `roleDefinitionId` | `string` | Pinned GUID for the definition's ARM resource name; generated when omitted. |

## Examples

### Broad Grant with an RBAC Carve-Out

"Contributor, except changing RBAC" -- the classic wildcard-minus-carve-out:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleDefinition
metadata:
  name: project-admin
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRoleDefinition.project-admin
spec:
  name: acme-project-admin
  scope:
    value: /subscriptions/00000000-0000-0000-0000-000000000000
  description: Manage all project resources but never RBAC or policy
  permissions:
    - actions:
        - "*"
      notActions:
        - Microsoft.Authorization/*/write
        - Microsoft.Authorization/*/delete
```

### Data-Plane Role for Blob Access

Control-plane actions manage resources; data-plane actions read the data inside them. A role that reads blob contents needs `dataActions`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleDefinition
metadata:
  name: blob-auditor
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRoleDefinition.blob-auditor
spec:
  name: acme-blob-auditor
  scope:
    value: /subscriptions/00000000-0000-0000-0000-000000000000
  description: Read blob data across the analytics accounts, nothing else
  permissions:
    - actions:
        - Microsoft.Storage/storageAccounts/read
      dataActions:
        - Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read
```

### Resource-Group-Scoped Role in a Composed Environment

Scope the definition to a deployed resource group by reference -- the default kind resolves it to the group's ARM ID:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleDefinition
metadata:
  name: env-operator
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRoleDefinition.env-operator
spec:
  name: acme-env-operator
  scope:
    valueFrom:
      name: platform-rg
  description: Restart and observe workloads in the platform environment
  permissions:
    - actions:
        - Microsoft.Compute/virtualMachines/read
        - Microsoft.Compute/virtualMachines/restart/action
        - Microsoft.Insights/metrics/read
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `role_definition_id` | `string` | Fully-scoped ARM ID of the definition (`{scope}/providers/Microsoft.Authorization/roleDefinitions/{guid}`) -- what an `AzureRoleAssignment`'s `role_definition_id` consumes |
| `role_definition_guid` | `string` | The definition's GUID resource name |
| `role_name` | `string` | The role's display name as deployed |
| `scope` | `string` | The scope the definition was created at |
| `assignable_scopes` | `list` | The assignable scopes Azure recorded (the definition's own scope when omitted) |

## Related Components

- [AzureRoleAssignment](/docs/catalog/azure/role-assignment) — grants this custom role to a principal at a scope
- [AzureUserAssignedIdentity](/docs/catalog/azure/azureuserassignedidentity) — the workload identity custom roles are most often granted to
- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — the default creation scope in composed environments
