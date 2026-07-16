# AzureRoleAssignment Terraform Module

## Overview

This Terraform module provisions an Azure RBAC role assignment using the
`azurerm` provider. It creates a single `azurerm_role_assignment` binding a
role (by built-in name or definition ID) to a principal at an ARM scope.

Role assignments are immutable in Azure: any change replaces the assignment
(delete + create). They carry no ARM tags -- the `Microsoft.Authorization`
resource type does not support them.

## Resources Created

- `azurerm_role_assignment.main` -- the Azure role assignment

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Role assignment specification (scope, role, principal, optional refinements) |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `scope` | yes | ARM scope of the grant (management group / subscription / resource group / resource ID) |
| `role_definition_name` | one-of | Built-in role name, e.g. "Reader" |
| `role_definition_id` | one-of | Fully-scoped role definition ID (custom roles) |
| `principal_id` | yes | Azure AD OBJECT ID of the grantee (not the client ID) |
| `principal_type` | no | "SERVICE_PRINCIPAL", "USER", or "GROUP"; inferred by Azure when omitted |
| `condition` / `condition_version` | no | Azure ABAC condition and its syntax version |
| `delegated_managed_identity_resource_id` | no | Cross-tenant (Lighthouse) delegation only |
| `skip_service_principal_aad_check` | no | Replication-lag escape hatch for freshly created service principals |
| `name` | no | Pinned GUID for the assignment's ARM resource name |

## Outputs

| Output | Description |
|--------|-------------|
| `role_assignment_id` | Fully-scoped ARM ID of the assignment |
| `name` | The assignment's GUID resource name |
| `scope` | The scope the role was granted at |
| `role_definition_id` | The role definition ID Azure resolved (even when a role name was given) |
| `principal_id` | The principal's Azure AD object ID |
| `principal_type` | The principal type Azure recorded |

## Usage

```hcl
module "role_assignment" {
  source = "./iac/tf"

  metadata = {
    name = "ci-acr-pull"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    scope                = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg"
    role_definition_name = "Reader"
    principal_id         = "11111111-1111-1111-1111-111111111111"
  }
}
```

## Required Permissions

The deploying credential needs `Microsoft.Authorization/roleAssignments/write`
at the target scope -- typically via Owner, User Access Administrator, or Role
Based Access Control Administrator. Contributor alone is NOT sufficient.
