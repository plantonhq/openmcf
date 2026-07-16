# AzureUserAssignedIdentity

## Overview

`AzureUserAssignedIdentity` provisions a user-assigned managed identity --
Azure's credential-free identity for workloads. Applications, clusters, and
pipelines authenticate as the identity without any client secret to store,
rotate, or leak.

Unlike system-assigned identities (which are born and destroyed with a single
resource), user-assigned identities exist independently and can be shared:
the same identity can back an AKS cluster's kubelets, a Function App, a
Container App, and a VM at once, and it survives all of them. That
independent lifecycle is what makes it the right anchor for permissions.

## Why Identity Only?

The identity is deliberately just the identity. The two things that make it
useful are separate, first-class resources that attach to its outputs:

- **What it may do** -- `AzureRoleAssignment` resources granting roles to its
  `principal_id`; each grant is individually reviewable and removable
- **Who may act as it** -- `AzureFederatedIdentityCredential` resources
  trusting external OIDC workloads (GitHub Actions, AKS service accounts)
  against its `identity_id`

Keeping the three concerns decomposed means changing a grant never touches
the identity, revoking a CI pipeline's trust never touches the grants, and
the environment graph shows exactly who can do what.

## Key Features

- **Credential-free** -- no secret exists for this identity anywhere
- **Shareable** -- one identity can back many consuming resources
- **Regional isolation (opt-in)** -- restrict token issuance to the
  identity's own region for data-residency requirements
- **Governance tags** -- user tags merge over the Planton-derived resource
  tags (user tags win), feeding Azure Policy and Cost Management
- **Composable** -- the resource group is referenced (defaults to an
  `AzureResourceGroup` output); every downstream identity consumer references
  this kind's outputs

## When to Use

- Before deploying AKS clusters, Function Apps, Web Apps, Container Apps, or
  VMs that need to access Azure services (Key Vault, Storage, ACR, ...)
- As the anchor of keyless CI/CD (with a federated credential) or AKS
  workload identity
- When multiple resources need to share one set of permissions, or an
  identity must outlive the resources attached to it

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | Yes | Azure region (the identity is a regional resource) |
| `resource_group` | StringValueOrRef | Yes | Resource group (literal name or AzureResourceGroup reference) |
| `name` | string | Yes | Identity name, unique in the resource group (3-128 chars); renaming replaces the identity and mints a new principal |
| `isolation_scope` | enum | No | `REGIONAL` restricts token issuance to the identity's region; unset means usable from any region |
| `tags` | map | No | User tags merged over the metadata-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `identity_id` | ARM ID of the identity (what consumers and federated credentials attach to) |
| `principal_id` | Service principal object ID (what role assignments grant to) |
| `client_id` | Client/application ID (what workloads present to authenticate, e.g. `AZURE_CLIENT_ID`) |
| `tenant_id` | Azure AD tenant ID |

## Quick Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureUserAssignedIdentity
metadata:
  name: platform-identity
  org: mycompany
  env: production
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: platform-rg
  name: prod-platform-identity
  tags:
    cost-center: platform
```

Grant it permissions with an `AzureRoleAssignment`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRoleAssignment
metadata:
  name: platform-identity-kv-reader
spec:
  scope:
    valueFrom:
      kind: AzureKeyVault
      name: platform-kv
      fieldPath: status.outputs.key_vault_id
  roleDefinitionName: Key Vault Secrets User
  principalId:
    valueFrom:
      name: platform-identity
```

## Downstream Resources

Resources that reference this identity:

- **AzureRoleAssignment** -- `principal_id` for permission grants
- **AzureFederatedIdentityCredential** -- `identity_id` for keyless external trust
- **AzureAksCluster** -- `identity_id` for kubelet or control plane identity
- **AzureFunctionApp** / **AzureLinuxWebApp** / **AzureContainerApp** --
  `identity_id` for managed identity assignment

## References

- [Azure Managed Identities overview](https://learn.microsoft.com/en-us/entra/identity/managed-identities-azure-resources/overview)
- [Azure RBAC built-in roles](https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles)
- Research documentation: [docs/README.md](docs/README.md)
