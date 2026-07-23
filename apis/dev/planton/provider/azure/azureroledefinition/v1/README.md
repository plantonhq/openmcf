# AzureRoleDefinition

## Overview

`AzureRoleDefinition` provisions a custom Azure RBAC role -- a named, reusable
bundle of permissions that role assignments then grant to principals. Azure
ships hundreds of built-in roles, but real organizations routinely need
permission sets the built-ins don't express; a custom role captures such a set
once, with a meaningful name, and every grant of it stays consistent as the
definition evolves.

## Why a First-Class Resource?

Custom roles are real infrastructure with their own lifecycle:

- **Defined once, granted many times** -- one definition backs any number of
  role assignments; updating its permissions updates what every existing
  assignment allows, without touching a single grant
- **Composition seam** -- the definition's fully-scoped ARM ID is exactly what
  an `AzureRoleAssignment`'s `role_definition_id` consumes, so custom roles and
  their grants are independent, referenceable nodes in an environment's graph
- **Governance artifact** -- a reviewed, version-controlled catalog of what
  each role in the organization actually permits, instead of portal-crafted
  one-offs

## Key Features

- **Full permission model** -- control-plane `actions`/`not_actions` and
  data-plane `data_actions`/`not_data_actions`, with ARM operation wildcards
- **Broad-grant-minus-carve-out** -- express "Contributor, except RBAC writes"
  style roles with `not_actions` trimming a wildcard grant
- **Any creation scope** -- management group, subscription, or resource group;
  the scope defaults to referencing an `AzureResourceGroup`, and subscription
  or management-group IDs pass as literals
- **Assignable-scope control** -- pre-authorize where the role may ever be
  granted; defaults to the definition's own scope
- **Pinnable GUID** -- preserve an externally-known role definition GUID when
  recreating a definition other tooling references

## When to Use

- Roles the built-ins don't express: "operate VMs but never create or delete",
  "read everything plus blob data", "manage a project's resources except RBAC"
- Establishing an organization-wide role catalog reviewed like code
- Pairing with `AzureRoleAssignment` to grant the custom role to workload
  identities, CI/CD principals, or teams

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | The role's tenant-unique display name |
| `scope` | StringValueOrRef | Yes | Creation scope (defaults to an AzureResourceGroup reference; subscription/management-group IDs as literals) |
| `description` | string | No | Intent text shown in the portal's role picker |
| `permissions` | list | No | Permission blocks: `actions`, `not_actions`, `data_actions`, `not_data_actions` |
| `assignable_scopes` | list of StringValueOrRef | No | Where the role may be assigned; defaults to `[scope]` |
| `role_definition_id` | string | No | Pinned GUID for the definition's ARM resource name |

Role definitions do not support ARM tags, so the spec carries none.

## Outputs

| Output | Description |
|--------|-------------|
| `role_definition_id` | Fully-scoped ARM ID of the definition (what a role assignment binds) |
| `role_definition_guid` | The definition's GUID resource name |
| `role_name` | The role's display name as deployed |
| `scope` | The scope the definition was created at |
| `assignable_scopes` | The assignable scopes Azure recorded |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleDefinition
metadata:
  name: vm-operator
  org: mycompany
  env: production
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

Grant it with an `AzureRoleAssignment` referencing the definition's output:

```yaml
spec:
  scope:
    valueFrom:
      name: platform-rg
  roleDefinitionId:
    # custom roles are bound by fully-scoped definition ID
  principalId:
    valueFrom:
      name: ops-identity
```

(The assignment's `role_definition_id` is a plain string field -- pass the
definition's `status.outputs.role_definition_id` through your chart's value
wiring.)

## Required Permissions

The deploying credential needs `Microsoft.Authorization/roleDefinitions/write`
at the target scope -- typically via Owner, User Access Administrator, or Role
Based Access Control Administrator. Contributor alone is NOT sufficient (it
manages resources, not authorization).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
