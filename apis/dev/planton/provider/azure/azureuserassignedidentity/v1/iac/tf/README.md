# AzureUserAssignedIdentity Terraform Module

## Overview

This Terraform module provisions a user-assigned managed identity using the
`azurerm` provider. It creates a single `azurerm_user_assigned_identity`:
the standalone, credential-free Azure AD identity workloads authenticate as.

Tags and isolation scope update in place; the name, region, and resource
group are the identity's ARM identity, so changing any of them replaces it --
minting a NEW principal and client ID that invalidates existing grants and
federated trust rules. The identity is deliberately just the identity:
grants live in `AzureRoleAssignment` and keyless trust rules in
`AzureFederatedIdentityCredential`, both referencing this module's outputs.

## Resources Created

- `azurerm_user_assigned_identity.main` -- the managed identity

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | User-assigned managed identity specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region (the identity is a regional resource) |
| `resource_group` | yes | Resource group name |
| `name` | yes | Identity name, unique within the resource group (3-128 chars) |
| `isolation_scope` | no | `REGIONAL` restricts token issuance to the identity's region; unset means usable from any region |
| `tags` | no | User tags merged over the metadata-derived tags (user wins) |

## Outputs

| Output | Description |
|--------|-------------|
| `identity_id` | ARM ID of the identity (what consumers and federated credentials attach to) |
| `principal_id` | Service principal object ID (what role assignments grant to) |
| `client_id` | Client/application ID (what workloads present to authenticate) |
| `tenant_id` | Azure AD tenant ID |

## Usage

```hcl
module "user_assigned_identity" {
  source = "./iac/tf"

  metadata = {
    name = "payments-identity"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    region         = "eastus"
    resource_group = "platform-rg"
    name           = "payments-api"
    tags = {
      owner = "payments-team"
    }
  }
}
```

## Required Permissions

The deploying credential needs
`Microsoft.ManagedIdentity/userAssignedIdentities/write` in the target
resource group -- held via Managed Identity Contributor, Contributor, or
Owner.
