# AzureSearchService Terraform Module

## Overview

Creates an Azure AI Search service -- the managed search-and-retrieval
engine AI applications index and query their own data with -- plus its
composed shared private links.

## Resources Created

- `azurerm_search_service` -- the service
  (`Microsoft.Search/searchServices`)
- `azurerm_search_shared_private_link_service` -- one per
  `sharedPrivateLinkServices` entry
  (`.../sharedPrivateLinkResources/{name}`), keyed by name

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureSearchServiceSpec fields; the resource group and
  identity references arrive as resolved literal values

## Outputs

- `search_service_id`, `search_service_name`, `endpoint`
- `primary_key`, `secondary_key`, `default_query_key` -- sensitive
  (service-minted credentials; empty in the RBAC-only posture)
- `customer_managed_key_encryption_compliance_status`
- `system_assigned_identity_principal_id` -- for indexer data-source
  grants, when a system identity is enabled
- `shared_private_link_service_ids` -- per-link ARM IDs, keyed by the
  spec entry's name

## Usage

The module is executed by the Planton platform with a tfvars file
converted from the manifest. To run it standalone, provide `metadata`
and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **SKU update contract**: in-place upgrade ONLY along basic ->
  standard -> standard2 -> standard3 (the provider's CustomizeDiff);
  every other SKU change replaces the service. Counts, auth, network,
  semantic, identity, and tags update in place; name, region, group,
  and hosting mode are ForceNew.
- **Per-SKU caps are spec CEL** -- by the time this module runs, the
  shape is legal.
- **Shared private links sit "Pending"** until the target resource's
  owner approves them; creation never needs the target's consent. The
  request_message is the only in-place-updatable field on a link.
- **Sensitive outputs**: the admin/query keys are provider-Sensitive;
  the outputs carry `sensitive = true` (a plan without it fails).

## Required Permissions

The deploying principal needs `Microsoft.Search/searchServices/*` on
the resource group (Contributor covers it). Shared private links
additionally need `Microsoft.Search/searchServices/sharedPrivateLinkResources/*`
and read access on the target resource.
