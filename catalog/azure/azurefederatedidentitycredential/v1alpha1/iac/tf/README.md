# AzureFederatedIdentityCredential Terraform Module

## Overview

This Terraform module provisions a federated identity credential using the
`azurerm` provider. It creates a single
`azurerm_federated_identity_credential`: a keyless OIDC trust rule on a
user-assigned managed identity that lets one external workload exchange its
own token for the identity's Azure credentials.

Issuer, subject, and audience update in place; the name and the parent
identity are the credential's ARM identity, so changing either replaces it.
The provider serializes credential writes per parent identity (ARM rejects
concurrent writes on one identity), so several credentials on the same
identity apply sequentially. Federated credentials carry no ARM tags (they
are untagged child resources of the identity). The resource group is derived
from the parent identity's ARM ID -- the module never asks for state that is
already encoded in the parent reference.

## Resources Created

- `azurerm_federated_identity_credential.main` -- the trust rule on the identity

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Federated identity credential specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `name` | yes | The credential's name under the parent identity (3-120 chars) |
| `user_assigned_identity` | yes | ARM ID of the parent user-assigned identity |
| `issuer` | yes | OIDC issuer URL the token's `iss` claim must equal |
| `subject` | yes | Workload identifier the token's `sub` claim must equal |
| `audience` | no | The token's required `aud` claim; defaults to `api://AzureADTokenExchange` |

## Outputs

| Output | Description |
|--------|-------------|
| `federated_identity_credential_id` | Full ARM ID of the credential |
| `name` | The credential's name as deployed |
| `user_assigned_identity_id` | ARM ID of the parent identity |
| `issuer` | The trusted issuer as deployed |
| `subject` | The trusted subject as deployed |
| `audience` | The required audience as deployed |

## Usage

```hcl
module "federated_identity_credential" {
  source = "./iac/tf"

  metadata = {
    name = "ci-deployer-trust"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    name                   = "github-main-branch"
    user_assigned_identity = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.ManagedIdentity/userAssignedIdentities/ci-deployer"
    issuer                 = "https://token.actions.githubusercontent.com"
    subject                = "repo:mycompany/platform:ref:refs/heads/main"
  }
}
```

## Required Permissions

The deploying credential needs
`Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials/write`
on the parent identity -- held via Managed Identity Contributor, Contributor,
or Owner on the identity's scope.
