# AzureCognitiveAccount Terraform Module

## Overview

Provisions an Azure AI services account (`azurerm_cognitive_account`) plus its composed responsible-AI children: blocklists (`azurerm_cognitive_account_rai_blocklist`) and content-filter policies (`azurerm_cognitive_account_rai_policy`), one per spec entry, keyed by name.

The spec's CEL contracts enforce every kind-gated rule the provider checks at apply time (project management / network injection only on `AIServices`, bypass only on the AI kinds, the QnAMaker / TextAnalytics / MetricsAdvisor field gates), so the module maps the legal shape without re-validating it.

## Resources Created

- `azurerm_cognitive_account` -- the account (endpoint, keys, kind, SKU, identity, perimeter, encryption)
- `azurerm_cognitive_account_rai_blocklist` (for_each over `spec.rai_blocklists`, keyed by name)
- `azurerm_cognitive_account_rai_policy` (for_each over `spec.rai_policies`, keyed by name; depends on the blocklists so a policy can reference a blocklist defined in the same spec)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureCognitiveAccountSpec fields; `resource_group` and the other StringValueOrRef fields arrive as resolved literal strings

## Outputs

- `cognitive_account_id`, `cognitive_account_name`, `endpoint`
- `primary_access_key`, `secondary_access_key` (sensitive; empty when local auth is disabled)
- `system_assigned_identity_principal_id`
- `rai_blocklist_ids`, `rai_policy_ids` (maps keyed by the spec entries' names)

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **Soft delete**: deletion leaves a purgeable ghost holding the account name; the provider purges it on destroy by default (`purge_soft_delete_on_destroy`).
- **Enum wire maps**: identity types, network-ACL bypass, RAI severity levels and policy modes arrive as proto enum NAMES and are mapped to the provider's wire values in `locals.tf`; unspecified enums map to null so ARM defaults apply.
- **Optional-with-default booleans**: `local_auth_enabled` and `public_network_access_enabled` pass through as null when unset so the provider defaults (both true) apply.
- **Set-once subdomain**: the provider replaces the account only when CHANGING an existing `custom_subdomain_name`, not when adding one.
- **Agent network injection**: after the account deletes, ARM removes the subnet's service association link asynchronously -- the provider waits for it before finishing the destroy.

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege action manifest the deploying principal needs.
