# AzureCognitiveDeployment Terraform Module

## Overview

Provisions a model deployment (`azurerm_cognitive_deployment`) on an Azure AI services account -- the model (format/name/version), the SKU (throughput class + capacity), the version-upgrade policy, and the responsible-AI policy selection.

## Resources Created

- `azurerm_cognitive_deployment` -- one deployment; an ARM child of the referenced account

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureCognitiveDeploymentSpec fields; `cognitive_account_id` arrives as a resolved literal ARM ID

## Outputs

- `deployment_id`, `deployment_name`
- `model_version` -- the version ARM resolved (useful when the spec tracks the default)

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **No region / resource group / tags**: the deployment is an ARM child of its account -- the provider's schema carries none of the three (ARM derives them through the account).
- **Enum wire maps**: `sku.tier` and `version_upgrade_option` arrive as proto enum NAMES and map to wire values in `locals.tf`; unspecified maps to null (ARM derives the tier; the provider defaults the upgrade option to `OnceNewDefaultVersionAvailable`).
- **Omit-when-unset honesty**: `rai_policy_name` and `model.version` are Optional+Computed on the provider -- the module emits null when the spec leaves them empty so ARM's defaults apply without read drift.
- **Capacity**: unset applies the provider default of 1; updates in place.

## Required Permissions

The deploying principal needs `Microsoft.CognitiveServices/accounts/deployments/*` on the account's resource group (Cognitive Services Contributor).
