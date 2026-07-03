# AzureUserAssignedIdentity — Research Notes

Design notes for the user-assigned managed identity component: what the
resource is, why it is modeled as identity-only, and how the spec maps to
the provider surface.

## What a User-Assigned Managed Identity Is

A managed identity is an Azure AD service principal that Azure itself
manages -- no password, certificate, or key ever exists for it. Workloads
authenticate as the identity through Azure's platform channels (IMDS on
compute, workload-identity token exchange elsewhere) and receive OAuth
tokens for any Azure AD-protected service.

Azure offers two flavors:

- **System-assigned**: born and destroyed with a single resource, never
  shared. Enabled ON a resource, not deployed as one.
- **User-assigned**: a standalone ARM resource with an independent
  lifecycle, attachable to many resources at once, surviving all of them.
  This is what this component models.

The user-assigned form is the enterprise grain: it survives resource
recreation and blue/green swaps, permissions can be provisioned before the
consuming resources exist, and one identity can back every replica of a
horizontally-scaled workload.

The identity's three identifiers matter to different consumers:

| Identifier | Who uses it |
|------------|-------------|
| `identity_id` (ARM ID) | Resources that attach the identity (AKS, apps, VMs) and federated credentials parented on it |
| `principal_id` (directory object ID) | Role assignments -- RBAC grants target the principal |
| `client_id` (application ID) | The workload itself, to say WHICH identity to authenticate as (`AZURE_CLIENT_ID`) |

## Why Identity-Only (the decomposition)

The identity is modeled without any embedded grants. Permissions
(`AzureRoleAssignment`) and external trust (`AzureFederatedIdentityCredential`)
are separate first-class kinds referencing the identity's outputs, because:

- **A module must not own what it merely references.** A grant binds an
  identity to some OTHER resource's scope; embedding grants in the identity
  makes the identity's module mutate resources it doesn't own, and two
  identities granting at the same scope become order-dependent.
- **Grants and trust rules have their own lifecycles.** Permissions change
  as workloads evolve; trust rules change as pipelines come and go. Each
  should be addable and revocable without touching (or risking) the
  identity -- whose replacement would mint a new principal and invalidate
  everything.
- **Full grant surface, once.** The standalone role-assignment kind carries
  ABAC conditions, principal types, custom-role binding, and cross-tenant
  fields; an embedded mini-grant would either duplicate that surface or
  permanently lag it.

## Provider Surface Map (azurerm v4)

`azurerm_user_assigned_identity`
(`internal/services/managedidentity/user_assigned_identity_resource.go`):

| azurerm argument | Spec field | Notes |
|------------------|-----------|-------|
| `name` (required, ForceNew) | `name` | Unique per resource group; replacement mints a NEW principal |
| `resource_group_name` (required, ForceNew) | `resource_group` | FK to AzureResourceGroup's name output |
| `location` (required, ForceNew) | `region` | The catalog's Azure-region field name |
| `isolation_scope` (optional) | `isolation_scope` (enum) | Only `Regional` is exposed by the provider; ARM's default (`None`) is the unspecified enum value. In-place update |
| `tags` (optional) | `tags` | Merged over the metadata-derived tags; user wins. In-place update |

Attributes → outputs: `id` → `identity_id`, `principal_id`, `client_id`,
`tenant_id` -- the complete attribute set the provider exports.

Engine parity notes:

- `isolation_scope` reached the classic Pulumi provider in the same bridged
  line as azurerm's own addition; both modules send it only when set to
  `Regional`, so an unspecified spec and ARM's default deploy identically.
- Both modules build the provider through the shared Azure builder (static
  client secret, keyless web identity, or ambient chain).

### Deliberately not modeled (recorded reasons)

- **Embedded role assignments** -- decomposed to `AzureRoleAssignment` (see
  "Why Identity-Only" above).
- **System-assigned identities** -- not a deployable resource; they are a
  property of the consuming resource and are modeled there.

## Operational Notes

- **Replacement invalidates everything.** Renaming or moving the identity
  replaces it, minting a new `principal_id`/`client_id`. Grants and
  federated credentials in the same composed environment re-resolve on the
  next deploy; externally-configured consumers (an `AZURE_CLIENT_ID` in some
  app's config) must be updated by hand. Name identities durably.
- **AAD replication lag.** A grant created immediately after the identity
  can hit "principal not found" in Azure AD; the role-assignment kind's
  `skip_service_principal_aad_check` exists for exactly this
  identity-plus-grant-in-one-deployment case.
- **Region and usability.** The identity is a regional resource, but by
  default it is usable by resources in ANY region; `REGIONAL` isolation is
  the opt-in restriction for data-residency requirements.
- **Free.** Managed identities cost nothing; prefer one identity per
  workload/duty over shared broad-grant identities.
