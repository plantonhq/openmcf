# Azure Role Assignment

Grants an Azure RBAC role to a principal at a scope -- the atomic unit of Azure authorization, modeled as a first-class composable resource. Use it to give a managed identity, user, group, or service principal exactly the access it needs on exactly the resources it touches.

## What Gets Created

When you deploy an AzureRoleAssignment resource, Planton provisions:

- **Role Assignment** — an `azurerm_role_assignment` binding the role (by built-in name or definition ID) to the principal at the target ARM scope

Role assignments carry no Azure tags -- the `Microsoft.Authorization` resource type does not support them.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **Authorization-plane rights**: the deploying credential needs `Microsoft.Authorization/roleAssignments/write` at the target scope (Owner, User Access Administrator, or Role Based Access Control Administrator -- Contributor alone is not sufficient)
- **The principal's OBJECT ID** (not the application/client ID) when passing it as a literal

## Quick Start

Create a file `roleassignment.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleAssignment
metadata:
  name: app-reader
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureRoleAssignment.app-reader
spec:
  scope:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg
  roleDefinitionName: Reader
  principalId:
    value: 11111111-1111-1111-1111-111111111111
```

Deploy:

```shell
planton apply -f roleassignment.yaml
```

This grants the principal read-only access to everything in `platform-rg`.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `scope` | `StringValueOrRef` | The ARM scope the role applies at: management group, subscription, resource group, or a single resource ID. Defaults to referencing an `AzureResourceGroup`'s ARM ID. | Required |
| `principalId` | `StringValueOrRef` | The Azure AD object ID of the grantee. Defaults to referencing an `AzureUserAssignedIdentity`'s principal ID output. | Required |
| `roleDefinitionName` / `roleDefinitionId` | `string` | Exactly one identifies the role: built-in role name (e.g. `Reader`) or fully-scoped role definition ID (custom roles). | Exactly one of the two |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `principalType` | `enum` | `SERVICE_PRINCIPAL`, `USER`, or `GROUP`. Azure infers the type when omitted; set it explicitly when the deploying credential is constrained by ABAC delegation rules. |
| `description` | `string` | Audit note shown in the portal's IAM blade -- record WHY the grant exists. |
| `condition` | `string` | Azure ABAC condition expression narrowing when the role applies. |
| `conditionVersion` | `string` | `1.0` or `2.0`; only with `condition` (Azure defaults to `2.0`). |
| `delegatedManagedIdentityResourceId` | `string` | Cross-tenant (Azure Lighthouse) delegation only. |
| `skipServicePrincipalAadCheck` | `bool` | Skip the Azure AD existence check for a freshly created service principal (replication-lag escape hatch). |
| `name` | `string` | Pinned GUID for the assignment's ARM resource name; generated when omitted. |

## Examples

### Grant a Workload Identity Pull Access to a Registry

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleAssignment
metadata:
  name: app-acr-pull
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRoleAssignment.app-acr-pull
spec:
  scope:
    valueFrom:
      kind: AzureContainerRegistry
      name: platform-acr
      fieldPath: status.outputs.registry_id
  roleDefinitionName: AcrPull
  principalId:
    valueFrom:
      name: app-identity
  skipServicePrincipalAadCheck: true
  description: Application workload identity pulls images from the platform registry
```

### ABAC-Conditioned Storage Grant

Read access limited to blobs tagged `Project=Cascade`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleAssignment
metadata:
  name: cascade-blob-reader
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRoleAssignment.cascade-blob-reader
spec:
  scope:
    valueFrom:
      kind: AzureStorageAccount
      name: analytics-storage
      fieldPath: status.outputs.storage_account_id
  roleDefinitionName: Storage Blob Data Reader
  principalId:
    value: 22222222-2222-2222-2222-222222222222
  condition: >-
    ((!(ActionMatches{'Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read'}))
    OR (@Resource[Microsoft.Storage/storageAccounts/blobServices/containers/blobs/tags:Project<$key_case_sensitive$>] StringEquals 'Cascade'))
  conditionVersion: "2.0"
  description: Analysts read only Cascade-tagged blobs
```

### Custom Role at Subscription Scope

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleAssignment
metadata:
  name: cost-auditor
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureRoleAssignment.cost-auditor
spec:
  scope:
    value: /subscriptions/00000000-0000-0000-0000-000000000000
  roleDefinitionId: /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
  principalId:
    value: 33333333-3333-3333-3333-333333333333
  principalType: GROUP
  description: FinOps group audits costs subscription-wide via the custom CostAuditor role
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `role_assignment_id` | `string` | Fully-scoped ARM ID of the assignment (`{scope}/providers/Microsoft.Authorization/roleAssignments/{guid}`) |
| `name` | `string` | The assignment's GUID resource name |
| `scope` | `string` | The scope the role was granted at |
| `role_definition_id` | `string` | The role definition ID Azure resolved (even when the spec referenced the role by name) |
| `principal_id` | `string` | The principal's Azure AD object ID |
| `principal_type` | `string` | The principal type Azure recorded |

## Related Components

- [AzureUserAssignedIdentity](/docs/catalog/azure/azureuserassignedidentity) — the managed identity most grants target
- [AzureResourceGroup](/docs/catalog/azure/azureresourcegroup) — the most common grant scope
- [AzureKeyVault](/docs/catalog/azure/azurekeyvault) — grant identities secret/key/certificate access in RBAC mode
- [AzureContainerRegistry](/docs/catalog/azure/azurecontainerregistry) — grant AcrPull/AcrPush to workload identities
- [AzureStorageAccount](/docs/catalog/azure/azurestorageaccount) — grant blob/queue/table data-plane roles
