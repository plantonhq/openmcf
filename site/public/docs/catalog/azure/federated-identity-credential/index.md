---
title: "Federated Identity Credential"
description: "Federated Identity Credential deployment documentation"
icon: "package"
order: 100
componentName: "azurefederatedidentitycredential"
---

# Azure Federated Identity Credential

Creates a federated identity credential -- a keyless OIDC trust rule on a user-assigned managed identity that lets an external workload (a GitHub Actions workflow, a Kubernetes service account, any OIDC-issuing system) exchange its own token for the identity's Azure credentials. No client secret is created, stored, or rotated.

## What Gets Created

When you deploy an AzureFederatedIdentityCredential resource, Planton provisions:

- **Federated Identity Credential** — an `azurerm_federated_identity_credential` on the referenced user-assigned identity, encoding one issuer + subject + audience trust triple

Federated credentials carry no Azure tags -- ARM models them as untagged child resources of the identity.

## Prerequisites

- **Azure credentials** configured via environment variables or Planton provider config
- **A user-assigned managed identity** to write the credential on (an `AzureUserAssignedIdentity` in composed environments)
- **Write rights on the identity**: `Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/write` (Managed Identity Contributor, Contributor, or Owner)

## Quick Start

Create a file `federated-credential.yaml`:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFederatedIdentityCredential
metadata:
  name: ci-deployer-trust
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFederatedIdentityCredential.ci-deployer-trust
spec:
  name: github-main-branch
  userAssignedIdentity:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ci-deployer
  issuer: https://token.actions.githubusercontent.com
  subject: repo:my-org/platform:ref:refs/heads/main
```

Deploy:

```shell
planton apply -f federated-credential.yaml
```

This lets the `main`-branch workflows of `my-org/platform` authenticate to Azure as the `ci-deployer` identity -- with no secret anywhere.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `name` | `string` | The credential's name under the parent identity; name it after the workload it trusts. | Required, 3-120 characters |
| `userAssignedIdentity` | `StringValueOrRef` | ARM ID of the parent identity. Defaults to referencing an `AzureUserAssignedIdentity`'s `identity_id` output. | Required |
| `issuer` | `string` | OIDC issuer URL the incoming token's `iss` claim must equal exactly. | Required, valid URI |
| `subject` | `string` | Workload identifier the token's `sub` claim must equal exactly (no wildcards). | Required, non-empty |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `audience` | `string` | The token's required `aud` claim. Defaults to `api://AzureADTokenExchange` -- the value every standard client requests; override only for providers that mint a different audience. |

## Examples

### GitHub Actions: Deploy from a Protected Environment

Trust only jobs running against the `production` environment (with its required reviewers), not every push to a branch:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFederatedIdentityCredential
metadata:
  name: prod-deploy-trust
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFederatedIdentityCredential.prod-deploy-trust
spec:
  name: github-production-environment
  userAssignedIdentity:
    valueFrom:
      name: prod-deployer-identity
  issuer: https://token.actions.githubusercontent.com
  subject: repo:my-org/platform:environment:production
```

### AKS Workload Identity for a Service Account

Let pods running as the `payments-api` service account authenticate as a managed identity. The issuer is the cluster's OIDC issuer URL:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFederatedIdentityCredential
metadata:
  name: payments-workload-trust
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.AzureFederatedIdentityCredential.payments-workload-trust
spec:
  name: aks-payments-serviceaccount
  userAssignedIdentity:
    valueFrom:
      name: payments-identity
  issuer: https://eastus.oic.prod-aks.azure.com/00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111/
  subject: system:serviceaccount:payments:payments-api
```

### One Identity, Several Trusted Workloads

An identity accumulates one credential per external workload (up to 20). Deploy several AzureFederatedIdentityCredential resources referencing the same identity -- one for `main`-branch CI, one for the release tags, one for a cluster service account -- each an independent, individually removable trust rule.

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `federated_identity_credential_id` | `string` | Full ARM ID of the credential (`{identity-id}/federatedIdentityCredentials/{name}`) |
| `name` | `string` | The credential's name as deployed |
| `user_assigned_identity_id` | `string` | ARM ID of the parent identity |
| `issuer` | `string` | The trusted issuer as deployed |
| `subject` | `string` | The trusted subject as deployed |
| `audience` | `string` | The required audience as deployed |

## Related Components

- [AzureUserAssignedIdentity](/docs/catalog/azure/user-assigned-identity) — the parent identity the external workload acts as
- [AzureRoleAssignment](/docs/catalog/azure/role-assignment) — grants the identity the permissions the external workload needs
- [AzureRoleDefinition](/docs/catalog/azure/role-definition) — custom roles for grants the built-ins don't express
