# AzureDataProtectionResourceGuard Terraform Module

## Overview

Creates a Data Protection Resource Guard -- the approval gate behind Multi-User Authorization for Azure Backup vaults. The guard is a free governance object; its protection comes entirely from living in a scope a DIFFERENT administrator controls than the vaults it guards.

## Resources Created

- `azurerm_data_protection_resource_guard` -- the guard

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureDataProtectionResourceGuardSpec fields; the resource group reference arrives as a resolved literal

## Outputs

- `resource_guard_id` -- the guard's full ARM ID; what backup vaults reference to enable Multi-User Authorization
- `resource_guard_name` -- the guard's name

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **An empty exclusion list guards everything** -- the strongest posture, and the module sends null (not an empty list) so the provider's optional-argument semantics hold exactly.
- **The exclusion list updates in place**; name, resource group, and location replace the guard (the provider's ForceNew).
- **The guard deletes cleanly even while vaults reference it** -- and the approval requirement goes with it. Gate guard deletion with resource locks or pipeline approvals; Azure does not gate the guard with itself.

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege control-plane actions the deploying principal needs.
