# AzureRoleDefinition: Research & Design Documentation

## 1. What Is an Azure Role Definition?

An Azure role definition is a named collection of permissions -- the "WHAT is
permitted" half of Azure RBAC. Azure ships several hundred built-in
definitions (Reader, Contributor, Storage Blob Data Reader, ...), all defined
at tenant level and immutable. A CUSTOM role definition is the same object
created by an organization: a tenant-unique display name, one or more
permission blocks, and a set of scopes where it may be assigned.

A definition grants nothing by itself. Permissions only take effect when a
role ASSIGNMENT binds the definition to a principal at a scope. The two kinds
compose: `AzureRoleDefinition` produces a fully-scoped definition ID, and
`AzureRoleAssignment.role_definition_id` consumes it.

### Anatomy

- **name** -- the display name shown in the portal's role picker. Unique per
  Azure AD tenant; renames are in-place (assignments track the GUID).
- **scope** -- where the definition is created: a management group,
  subscription, or resource group. Anchors the definition's ARM resource ID:
  `{scope}/providers/Microsoft.Authorization/roleDefinitions/{guid}`.
- **permissions** -- a LIST of blocks, each with four operation lists:
  `actions` / `not_actions` (control plane) and `data_actions` /
  `not_data_actions` (data plane). Azure evaluates blocks as a union; one
  block is the norm.
- **assignable_scopes** -- the scopes at which assignments of this role may
  be created. Defaults to `[scope]`. At most ONE management group may appear.
- **role_definition_id** -- the GUID resource name; generated when not
  pinned.

### Control plane vs data plane

The single most important distinction when authoring permissions:

- `actions` govern MANAGING resources through ARM -- create a storage
  account, restart a VM, read a vault's properties.
- `data_actions` govern the DATA inside resources -- read a blob's bytes,
  receive a queue message, read a secret's value (RBAC mode).

A role with `Microsoft.Storage/storageAccounts/*` in actions but no
data_actions can delete the whole account yet cannot read one byte of blob
data. The operation catalog per provider is queryable:
`az provider operation show --namespace Microsoft.Storage`.

### not_actions are carve-outs, not denies

A `not_actions` entry subtracts from THIS role's `actions` grant; it does not
(and cannot) take away permissions another assignment gives the same
principal. Azure's actual deny mechanism is deny assignments, a separate
system created by Blueprints/managed apps. The classic carve-out pattern:
`actions: ["*"]`, `not_actions: ["Microsoft.Authorization/*/write"]` --
everything except changing RBAC.

## 2. Authorization Model Nuances

### Who can create definitions

Creating, updating, or deleting a custom role requires
`Microsoft.Authorization/roleDefinitions/write` at the target scope -- held
via Owner, User Access Administrator, or Role Based Access Control
Administrator. Contributor manages resources, not authorization, and is NOT
sufficient. This is the same authorization-plane boundary as role
assignments.

### Tenant limits and uniqueness

- A tenant holds at most **5,000 custom roles** (500 in Azure Government /
  Azure China). Definitions are metadata; they cost nothing.
- Display names are tenant-unique. Azure rejects a duplicate name at create
  time, so org-prefixed names ("acme-vm-operator") avoid collisions with
  other teams' roles.

### Scope and assignable-scope mechanics

- Valid creation scopes: management group, subscription, resource group.
  azurerm validates the prefix (`/subscriptions/` or
  `/providers/Microsoft.Management/managementGroups/`), which admits resource
  groups (their IDs start with `/subscriptions/`); azurerm's own acceptance
  suite exercises resource-group-scoped creation.
- `assignable_scopes` entries may be management groups, subscriptions,
  resource groups, or individual resources. An assignment whose scope is not
  at or under one of them is rejected. When the list is omitted, Azure
  defaults it to the creation scope -- both engines inherit this server-side
  defaulting identically.
- Deleting a definition that still has assignments fails; assignments must go
  first. In a composed environment the dependency graph's reverse destroy
  order does this naturally (assignments reference the definition).

### Eventual consistency

"Updating" a role definition actually creates a new record that Azure
consolidates over the following seconds. azurerm deliberately polls until the
created/updated timestamps settle (roughly 3+ minutes), and deletion likewise
waits for consecutive not-found responses (~4 minutes). Both engines inherit
this via the shared azurerm logic (pulumi-azure bridges it). The deployed
permissions are correct as soon as the apply returns; the waits exist so
subsequent reads are stable.

## 3. Provider Surface (the completeness floor)

`azurerm_role_definition` (terraform-provider-azurerm v4.80,
`internal/services/authorization/role_definition_resource.go`):

| azurerm field | Spec field | Notes |
|---|---|---|
| `name` (required) | `name` | tenant-unique display name; in-place update |
| `scope` (required, ForceNew) | `scope` (StringValueOrRef) | prefix-validated; anchors the ARM ID |
| `description` | `description` | in-place update |
| `permissions` (list) | `permissions` (repeated message) | 4 operation lists per block |
| `assignable_scopes` (optional, computed) | `assignable_scopes` (repeated StringValueOrRef) | defaults to `[scope]` server-side |
| `role_definition_id` (optional, computed, ForceNew) | `role_definition_id` | pinned GUID; UUID-validated |
| `role_definition_resource_id` (attribute) | output `role_definition_id` | the fully-scoped ARM ID |

Parity verified against `pulumi-azure/sdk/v6 v6.28.0`
`authorization.RoleDefinition` -- field-for-field identical (both engines
wrap the same azurerm logic). Every spec field is honored by both modules.

Nothing from the azurerm surface is skipped. Role definitions carry no ARM
tags -- `Microsoft.Authorization` resources do not support them.

## 4. Design Decisions

### First-class kind completing the custom-role story

`AzureRoleAssignment` binds roles by name (built-ins) or by definition ID
(custom roles). Without this kind, the definition-ID path pointed at
hand-crafted roles outside the platform's graph. With it, an organization's
role catalog is version-controlled, reviewable infrastructure, and
assignments reference definitions as first-class nodes.

### Output naming: the fully-scoped ID is `role_definition_id`

azurerm names the bare GUID `role_definition_id` and the fully-scoped ARM ID
`role_definition_resource_id`. Planton's Azure surface consistently uses
`role_definition_id` for the FULLY-SCOPED form -- that is how the assignment
kind documents and consumes it -- so this component's primary output carries
that name and format, making the composition seam zero-translation. The bare
GUID is exported separately as `role_definition_guid`.

### Reference defaults chosen for composed environments

`scope` (and each `assignable_scopes` element) defaults to referencing an
`AzureResourceGroup`'s ARM ID -- the grain of composed environments, where a
team-scoped custom role lives inside the environment's resource group.
Subscription- and management-group-scoped definitions (the common org-wide
case) pass the scope as a literal value; the reference machinery imposes no
ceremony on literals.

### Permissions mirror ARM's list-of-blocks shape

ARM models `permissions` as an array of blocks; azurerm and pulumi-azure
both surface it that way; this spec does too. Flattening to a single block
would be friendlier for the 90% case but would diverge from the provider
contract and forbid the legal multi-block form -- the spec stays faithful
and the presets teach the one-block norm.

### Empty permissions are legal

azurerm explicitly supports a definition with no permission blocks (its
acceptance suite covers it). The spec allows it (a placeholder role granting
nothing) rather than inventing a stricter-than-Azure rule.

## 5. Operational Notes

- **Slow updates and deletes are normal** -- eventual-consistency polling
  (see §2). Expect minutes, not seconds.
- **Renaming is safe** -- assignments track the GUID, so renaming a role does
  not affect existing grants.
- **Changing scope or the pinned GUID replaces the definition** -- both are
  part of its ARM identity. Replacement fails while assignments still
  reference the old definition; destroy those first.
- **Wildcards compound** -- `actions: ["*"]` includes every future operation
  ARM adds. Prefer explicit operation lists for least-privilege roles;
  reserve wildcards for broad-minus-carve-out roles.

## 6. Verification

- Spec tests: `go test ./apis/dev/planton/provider/azure/azureroledefinition/v1/`
- Offline module proof: `tofu init/validate/plan` against `iac/hack/manifest.yaml`;
  `go build` of the Pulumi module and release-equivalent entrypoint build
- Outputs conformance: `planton validate-outputs --kind AzureRoleDefinition`
  plus the `pkg/outputs` case
- Live E2E: `TestAzureRoleDefinition_Pulumi` / `TestAzureRoleDefinition_Terraform`
  in `e2e/azure` -- creates a custom role at an ephemeral resource-group
  scope, verifies via the authorization API (GetByID), destroys, and verifies
  absence, on both engines
