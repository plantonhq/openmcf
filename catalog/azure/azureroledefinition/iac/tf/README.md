# AzureRoleDefinition Terraform Module

## Overview

This Terraform module provisions a custom Azure RBAC role definition using the
`azurerm` provider. It creates a single `azurerm_role_definition`: a named,
reusable bundle of control-plane and data-plane permissions that role
assignments can then grant to principals.

Name, description, permissions, and assignable scopes update in place; the
creation scope and the pinned GUID are the definition's ARM identity, so
changing either replaces it. Updates and deletes are eventually consistent --
azurerm polls until Azure settles, so those applies take a few minutes. Role
definitions carry no ARM tags (the `Microsoft.Authorization` resource type
does not support them).

## Resources Created

- `azurerm_role_definition.main` -- the custom role definition

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Role definition specification (name, scope, permissions, optional refinements) |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | Tenant-unique display name of the role |
| `scope` | yes | Creation scope (management group / subscription / resource group ID) |
| `description` | no | Intent text shown in the portal's role picker |
| `permissions` | no | Permission blocks: `actions`, `not_actions`, `data_actions`, `not_data_actions` |
| `assignable_scopes` | no | Where the role may be assigned; defaults to `[scope]` |
| `role_definition_id` | no | Pinned GUID for the definition's ARM resource name |

## Outputs

| Output | Description |
|--------|-------------|
| `role_definition_id` | Fully-scoped ARM ID of the definition (what a role assignment binds) |
| `role_definition_guid` | The definition's GUID resource name |
| `role_name` | The role's display name as deployed |
| `scope` | The scope the definition was created at |
| `assignable_scopes` | The assignable scopes Azure recorded |

## Usage

```hcl
module "role_definition" {
  source = "./iac/tf"

  metadata = {
    name = "vm-operator"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    name        = "acme-vm-operator"
    scope       = "/subscriptions/00000000-0000-0000-0000-000000000000"
    description = "Operate existing VMs: start/stop/restart, no create or delete"
    permissions = [{
      actions = [
        "Microsoft.Compute/virtualMachines/read",
        "Microsoft.Compute/virtualMachines/start/action",
        "Microsoft.Compute/virtualMachines/restart/action",
        "Microsoft.Compute/virtualMachines/deallocate/action",
      ]
    }]
  }
}
```

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege action set the deploying principal needs.
