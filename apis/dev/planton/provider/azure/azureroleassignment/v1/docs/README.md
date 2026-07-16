# AzureRoleAssignment: Research & Design Documentation

## 1. What Is an Azure Role Assignment?

An Azure role assignment is the mechanism by which Azure RBAC grants access: it
binds a **role definition** (a set of permitted actions) to a **principal** (a
user, group, service principal, or managed identity) at a **scope** (a slice of
the ARM resource hierarchy). Azure evaluates every management- and data-plane
request against the union of the caller's role assignments.

### The three coordinates

- **Scope** -- any node in the ARM hierarchy:
  management group → subscription → resource group → resource. Permissions
  granted at a scope inherit downward to everything beneath it. Least privilege
  means assigning at the narrowest scope that satisfies the need.
- **Role definition** -- either one of Azure's several hundred built-in roles
  (Reader, Contributor, Key Vault Secrets User, AcrPull, ...) or a custom role
  definition. A definition is referenced by name (resolved at the target scope)
  or by its fully-scoped definition ID.
- **Principal** -- the Azure AD (Entra ID) **object ID** of the grantee. This
  is not the application (client) ID: an application registration has both, and
  confusing them is the most common failure mode -- the assignment is accepted
  but grants nothing, because no directory object carries that object ID.

### Key properties

- **Immutable**: every property of an assignment is create-only. Changing the
  role, scope, or principal is a delete + create. The ARM object is an atomic
  grant record, which is why IaC engines mark every field ForceNew.
- **Identified by GUID**: the ARM resource name of an assignment is a GUID,
  generated at create time unless pinned. The full ARM ID is
  `{scope}/providers/Microsoft.Authorization/roleAssignments/{guid}`.
- **No tags**: `Microsoft.Authorization` resources do not support ARM tags.
- **Replication lag**: a freshly created service principal or managed identity
  may not have replicated through Azure AD when the assignment is submitted,
  producing `PrincipalNotFound`. The API offers an explicit escape hatch
  (`skip_service_principal_aad_check`) that bypasses the directory existence
  check -- only valid for service principals.
- **Limit**: 4,000 role assignments per subscription (a real constraint for
  automation-heavy organizations; group-based assignment is the mitigation).

## 2. Authorization Model Nuances

### Who can create assignments

Writing a role assignment requires
`Microsoft.Authorization/roleAssignments/write` at (or above) the target scope.
The built-in roles that carry it: **Owner**, **User Access Administrator**, and
**Role Based Access Control Administrator** (the least-privileged of the three,
designed for exactly this delegation). **Contributor explicitly excludes it** --
it manages resources, not authorization. Deploys failing with
`AuthorizationFailed` on assignment creation almost always trace to this.

### ABAC conditions

An assignment may carry a **condition** -- an attribute-based access control
expression that narrows when the role's permissions apply (e.g. "only blobs
tagged Project=Cascade"). Conditions have a **condition_version**; "2.0" is the
generally available syntax ("1.0" is legacy). Azure applies "2.0" when a
condition is supplied without a version. Conditions are supported on storage
data-plane roles and a growing set of others; an unsupported pairing fails at
create with a descriptive ARM error. Additionally, when the *creator* of an
assignment is itself constrained by an ABAC delegation condition that filters
on principal type, Azure requires `principal_type` to be declared explicitly on
the request -- which is why the spec models it.

### Cross-tenant delegation (Azure Lighthouse)

In managed-service-provider topologies, an identity in the managing tenant
creates assignments on resources in the customer tenant. The API models this
with `delegatedManagedIdentityResourceId` -- the ARM ID of the managed identity
in the delegated tenant. Single-tenant deployments (the overwhelming majority)
leave it unset.

## 3. Provider Surface (the completeness floor)

The Terraform AzureRM provider's `azurerm_role_assignment` resource defines the
authoritative field surface, all ForceNew:

| Provider field | Spec field | Notes |
|---|---|---|
| `scope` | `scope` | required; any ARM scope |
| `role_definition_id` XOR `role_definition_name` | same pair | exactly-one-of, enforced at validation |
| `principal_id` | `principal_id` | required; object ID |
| `principal_type` | `principal_type` | enum User/Group/ServicePrincipal |
| `description` | `description` | audit note |
| `condition` + `condition_version` | same | ABAC |
| `delegated_managed_identity_resource_id` | same | cross-tenant |
| `skip_service_principal_aad_check` | same | replication-lag escape hatch |
| `name` | `name` | optional pinned GUID |

Every field is modeled -- there are no skipped fields for this kind. The
provider's create path adds two behaviors both engines inherit: a role name is
resolved to its definition ID by listing definitions at the target scope, and
`PrincipalNotFound` responses are retried while directory replication catches
up.

## 4. Design Decisions

### First-class kind, not a bundled sub-resource

A grant connects two resources (principal, scope) without owning either.
Bundling assignments inside an identity component couples the identity's
lifecycle to its permissions, prevents granting roles to principals the
platform did not create (users, groups), and violates the composition principle
that a module never mutates a resource it merely references. Assignments are
many-per-principal and many-per-scope with independent lifecycles -- the
textbook profile of a standalone kind.

### Reference defaults chosen for the dominant composition

`scope` defaults to an `AzureResourceGroup` reference (its ARM ID output)
because the resource group is the most common grant boundary in composed
environments; `principal_id` defaults to an `AzureUserAssignedIdentity`
reference (its principal-ID output) because granting a workload identity access
to its dependencies is the dominant pattern. Both are genuinely polymorphic --
an explicit `kind`/`fieldPath` in the reference overrides the default, so any
resource's ID output can be a scope and any object ID a principal.

### Two role fields, not one

Azure exposes two equally-valid ways to identify a role: by name (ergonomic for
the several hundred built-ins) and by definition ID (exact, required for custom
roles created at arbitrary scopes). Collapsing them into one field would force
either ID-only friction on the common case or fragile string sniffing. The spec
keeps both with an exactly-one-of validation, mirroring the provider contract.

### Plain bool for the replication-lag flag

`skip_service_principal_aad_check` defaults to false in Azure; a plain proto3
bool matches that default exactly. The spec comment documents the one scenario
where setting it is correct (assignment immediately following principal
creation), because misuse on user/group principals fails the deploy.

## 5. Operational Notes

- **PrincipalNotFound on fresh identities**: when an assignment deploys in the
  same operation as the identity it targets, either rely on the provider's
  built-in retry or set `skip_service_principal_aad_check: true` (service
  principals only). Composed deployments that create identity + grant together
  should prefer the flag -- it removes minutes of retry latency.
- **Auditability**: `description` is the only free-text surface an operator
  sees in the portal's IAM blade; recording WHY a grant exists there pays off
  during access reviews.
- **Drift**: because assignments are immutable, drift takes exactly one form --
  the assignment exists or it does not. Re-applying recreates a deleted grant.
- **Deletion**: deleting an assignment revokes access within about a minute
  (token lifetime aside: already-issued tokens remain valid until expiry, so
  revocation is not instantaneous for data-plane access).

## 6. Verification

The component's live E2E provisions a resource group and a user-assigned
managed identity as fixtures, grants the identity the built-in **Reader** role
at the resource-group scope (by role name, exercising the name→ID resolution
path, with the replication-lag flag set as the fresh-principal best practice),
verifies the assignment through the Azure authorization API by its fully-scoped
ID, destroys it, and verifies it is gone -- on both IaC engines.
