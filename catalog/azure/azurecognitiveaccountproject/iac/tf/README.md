# AzureCognitiveAccountProject Terraform Module

## Overview

Provisions an AI Foundry project (`azurerm_cognitive_account_project`) on an Azure AI services account -- the team workspace with its own managed identity, description, display name and tags.

## Resources Created

- `azurerm_cognitive_account_project` -- one project; an ARM child of the referenced account

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureCognitiveAccountProjectSpec fields; `cognitive_account_id` and `identity.identity_ids` arrive as resolved literal strings

## Outputs

- `project_id`, `project_name`
- `endpoints` -- the data-plane endpoints map as ARM reports it
- `is_default` -- whether ARM made this the account's default project
- `system_assigned_identity_principal_id`

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Identity is required** by the provider -- the module always emits the block; the type arrives as a proto enum NAME and maps to the wire value in `locals.tf`.
- **The empty-update quirk**: ARM cannot update `description` or `display_name` to an empty value -- the provider replaces the project when either is cleared (the spec documents it on the fields).
- **No resource group variable**: ARM derives the project's resource group through the account.

## Required Permissions

The deploying principal needs `Microsoft.CognitiveServices/accounts/projects/*` on the account's resource group (Cognitive Services Contributor).
