# AzureAiFoundryProject Terraform Module

## Overview

Creates an Azure AI Foundry project -- the workspace one AI team
works in, created inside an AzureAiFoundry hub and inheriting its
posture.

## Resources Created

- `azurerm_ai_foundry_project` -- the project; ARM-wise an ML
  workspace of kind "Project"
  (`Microsoft.MachineLearningServices/workspaces`), placed in the
  HUB's resource group

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureAiFoundryProjectSpec fields; the hub and
  identity references arrive as resolved literal ARM IDs

## Outputs

- `ai_foundry_project_id` -- the project's full ARM ID
- `ai_foundry_project_name` -- the project's name
- `project_guid` -- the project's immutable GUID
- `system_assigned_identity_principal_id` -- for per-team data
  grants, when a system identity is enabled

## Usage

The module is executed by the Planton platform with a tfvars file
converted from the manifest. To run it standalone, provide `metadata`
and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **No resource group argument** -- the provider derives the group
  from the hub reference; the project always lands in the hub's group.
- **ForceNew surface**: name, region, the hub linkage, and the
  high-business-impact flag replace the project; identity, primary
  identity, description, friendly name, and tags update in place.
- **high_business_impact_enabled is sent only when true**: the
  service flips it true when hub encryption applies; pinning false
  would fight the read-back (the flag is ForceNew).
