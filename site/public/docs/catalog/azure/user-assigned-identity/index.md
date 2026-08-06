---
title: "User Assigned Identity"
description: "User Assigned Identity deployment documentation"
icon: "package"
order: 100
componentName: "azureuserassignedidentity"
---

# Azure User Assigned Identity

Deploys a user-assigned managed identity -- the standalone, credential-free Azure AD identity workloads authenticate as. The identity is deliberately just the identity: grant it permissions with `AzureRoleAssignment` resources and trust external workloads to act as it with `AzureFederatedIdentityCredential` resources, each attaching to this identity's outputs.

## What Gets Created

When you deploy an AzureUserAssignedIdentity resource, Planton provisions:

- **User-Assigned Managed Identity** — an `azurerm_user_assigned_identity` in the specified region and resource group
- **Azure Tags** — the Planton-derived resource tags (resource name, kind, organization, environment) with your spec tags merged over them (user tags win on key collision)

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **An Azure Resource Group** where the identity will be created (can reference an AzureResourceGroup resource)

## Quick Start

Create a file `identity.yaml`:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureUserAssignedIdentity
metadata:
  name: my-identity
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AzureUserAssignedIdentity.my-identity
spec:
  region: eastus
  resourceGroup:
    value: my-rg
  name: my-identity
```

Deploy:

```shell
planton apply -f identity.yaml
```

This creates the identity. It has no permissions until `AzureRoleAssignment` resources grant them, and no external workload can act as it until `AzureFederatedIdentityCredential` resources trust one.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | Azure region for the identity (e.g., `eastus`, `westeurope`). The identity is a regional resource. | Required, minimum length 1 |
| `resourceGroup` | `StringValueOrRef` | Azure Resource Group name. Can reference an AzureResourceGroup resource via `valueFrom`. | Required |
| `name` | `string` | Name of the identity, unique within the resource group. Renaming replaces the identity and mints a new principal, invalidating existing grants. | Required, 3-128 characters |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `isolationScope` | `enum` | unset | `REGIONAL` restricts token issuance to the identity's own region (a data-residency control). Unset means the identity is usable by resources in any region. Updates in place. |
| `tags` | `map<string,string>` | `{}` | User tags merged over the Planton-derived tags; a user tag with the same key wins. Updates in place. |

## Examples

### Identity for a Workload

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureUserAssignedIdentity
metadata:
  name: payments-identity
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureUserAssignedIdentity.payments-identity
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: platform-rg
  name: payments-api
  tags:
    owner: payments-team
```

Then grant what it needs -- for example, read access to Key Vault secrets:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureRoleAssignment
metadata:
  name: payments-kv-reader
spec:
  scope:
    valueFrom:
      kind: AzureKeyVault
      name: platform-kv
      fieldPath: status.outputs.key_vault_id
  roleDefinitionName: Key Vault Secrets User
  principalId:
    valueFrom:
      name: payments-identity
```

### Regionally Isolated, Governance-Tagged Identity

For regulated environments -- token issuance restricted to the identity's region and tags aligned with the org's Azure Policy conventions:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureUserAssignedIdentity
metadata:
  name: regulated-identity
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureUserAssignedIdentity.regulated-identity
spec:
  region: westeurope
  resourceGroup:
    valueFrom:
      name: regulated-rg
  name: regulated-workload
  isolationScope: REGIONAL
  tags:
    cost-center: compliance
    data-classification: restricted
```

### The Keyless-CI Composition

The identity anchors keyless CI/CD: an `AzureRoleAssignment` grants it deployment rights, and an `AzureFederatedIdentityCredential` lets a GitHub Actions workflow act as it -- no secret anywhere:

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureUserAssignedIdentity
metadata:
  name: ci-deployer
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureUserAssignedIdentity.ci-deployer
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: platform-rg
  name: ci-deployer
  tags:
    purpose: ci-deployment
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `identity_id` | `string` | ARM ID of the identity. Format: `/subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{name}`. What consuming resources and federated credentials attach to. |
| `principal_id` | `string` | Service principal object ID in Azure AD. What `AzureRoleAssignment` grants roles to. |
| `client_id` | `string` | Client (application) ID. What workloads present to authenticate as the identity (e.g., `AZURE_CLIENT_ID`). |
| `tenant_id` | `string` | Azure AD tenant ID the identity belongs to. |

## Related Components

- [AzureRoleAssignment](/docs/catalog/azure/role-assignment) -- grants this identity its permissions
- [AzureFederatedIdentityCredential](/docs/catalog/azure/federated-identity-credential) -- lets external OIDC workloads act as this identity
- [AzureRoleDefinition](/docs/catalog/azure/role-definition) -- custom roles for grants the built-ins don't express
- [AzureResourceGroup](/docs/catalog/azure/resource-group) -- provides the resource group where the identity is created
- [AzureAksCluster](/docs/catalog/azure/aks-cluster) -- AKS clusters can use user-assigned identities for kubelet or control plane authentication
