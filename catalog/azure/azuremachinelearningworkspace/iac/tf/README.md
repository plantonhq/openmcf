# AzureMachineLearningWorkspace Terraform Module

## Overview

Provisions an Azure Machine Learning workspace (`azurerm_machine_learning_workspace`) with its identity, companion-service attachments, optional CMK encryption, managed virtual network, and serverless-compute settings -- plus the workspace's composed network outbound rules as ARM children.

## Resources Created

- `azurerm_machine_learning_workspace` -- the workspace
- `azurerm_machine_learning_workspace_network_outbound_rule_fqdn` -- one per `fqdn_outbound_rules` entry, keyed by name
- `azurerm_machine_learning_workspace_network_outbound_rule_private_endpoint` -- one per `private_endpoint_outbound_rules` entry, keyed by name
- `azurerm_machine_learning_workspace_network_outbound_rule_service_tag` -- one per `service_tag_outbound_rules` entry, keyed by name

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureMachineLearningWorkspaceSpec fields; every reference field (resource group, insights, vault, storage, registry, identities, subnet, CMK key) arrives as a resolved literal

## Outputs

- `machine_learning_workspace_id`, `machine_learning_workspace_name`
- `workspace_guid` -- the immutable GUID (distinct from the ARM ID)
- `discovery_url`, `system_assigned_identity_principal_id`
- `fqdn_outbound_rule_ids`, `private_endpoint_outbound_rule_ids`, `service_tag_outbound_rule_ids` -- name-keyed ARM ID maps

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Enum wire maps**: `kind`, `managed_network.isolation_mode`, `storage_account_access_type`, and the identity type arrive as proto enum NAMES and map to wire values in `locals.tf`; unspecified maps to null so provider/ARM defaults apply.
- **Omit-when-unset honesty**: `managed_network` is Optional+Computed on the provider -- when the block is absent, plans show it known-after-apply and the value is read back rather than defaulted.
- **Soft delete**: deleted workspaces hold their name until purged; the provider's `machine_learning.purge_soft_delete_on_destroy`-class features flag governs purge-on-destroy.
- **One rule namespace**: the three outbound-rule types share one ARM collection; the spec enforces cross-type name uniqueness before this module runs.
- **PARITY-EXCEPTION**: `storage_account_access_type` is Terraform-only -- the classic Pulumi SDK cannot express it; the Pulumi module fails loudly when it is set.

## Required Permissions

The deploying principal needs `Microsoft.MachineLearningServices/workspaces/*` plus read access on the referenced companion services (Contributor on the resource group covers it).
