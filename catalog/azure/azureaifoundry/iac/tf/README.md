# AzureAiFoundry Terraform Module

## Overview

Creates an Azure AI Foundry hub -- the shared foundation (security,
storage, network posture) a company's AI teams create their Foundry
projects in.

## Resources Created

- `azurerm_ai_foundry` -- the hub; ARM-wise an ML workspace of kind
  "Hub" (`Microsoft.MachineLearningServices/workspaces`)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureAiFoundrySpec fields; the resource group, key
  vault, storage, insights, registry, and identity references arrive
  as resolved literal values

## Outputs

- `ai_foundry_id` -- the hub's full ARM ID (what projects reference)
- `ai_foundry_name` -- the hub's name
- `workspace_guid` -- the hub's immutable GUID
- `discovery_url` -- the hub's regional discovery URL
- `system_assigned_identity_principal_id` -- for key vault / storage
  grants, when a system identity is enabled

## Usage

The module is executed by the Planton platform with a tfvars file
converted from the manifest. To run it standalone, provide `metadata`
and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **ForceNew surface**: name, region, resource group, key vault,
  storage, encryption, and the high-business-impact flag replace the
  hub; insights, registry, identity, public access, network mode,
  description, friendly name, and tags update in place.
- **Encryption key is VERSIONED**: the provider validates
  `encryption.key_id` as a versioned Key Vault key URL -- versionless
  URLs are rejected, and rotation does not auto-propagate (the hub's
  contract; differs from the classic ML workspace).
- **high_business_impact_enabled is sent only when true**: the
  service flips it true when encryption is on; pinning false would
  fight the read-back (the flag is ForceNew).
- **Soft delete**: a destroyed hub holds its name as a purgeable
  ghost until purged (the ML workspace class).

## Required Permissions

The deploying principal needs
`Microsoft.MachineLearningServices/workspaces/*` on the resource
group (Contributor covers it), plus read access on the referenced
key vault, storage account, and any insights/registry attachments.
