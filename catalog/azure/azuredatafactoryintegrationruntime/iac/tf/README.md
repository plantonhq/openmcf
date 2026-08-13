# AzureDataFactoryIntegrationRuntime Terraform Module

## Overview

Creates one integration runtime inside an Azure Data Factory -- the compute engine the factory's pipelines, data flows, and copy activities run on. The spec's variant block (exactly one of 3) selects the engine flavor; the matching resource below is count-armed on it.

## Resources Created

Exactly one of:

- `azurerm_data_factory_integration_runtime_azure` -- the managed data-flow compute (`azure` variant)
- `azurerm_data_factory_integration_runtime_azure_ssis` -- the managed SSIS package runtime (`azure_ssis` variant)
- `azurerm_data_factory_integration_runtime_self_hosted` -- the self-hosted agent registration (`self_hosted` variant)

## Inputs

- `metadata` -- resource metadata (name, id, org, env)
- `spec` -- the AzureDataFactoryIntegrationRuntime spec (see `variables.tf`, generated from the kind's proto contract)

## Outputs

- `integration_runtime_id` -- the runtime's ARM resource ID (`{factory_id}/integrationRuntimes/{name}`, the same shape for all 3 flavors)
- `integration_runtime_name` -- the runtime's name (what linked services, data flow activities, and the SSIS proxy resolve against)
- `primary_authorization_key` / `secondary_authorization_key` -- the keys a self-hosted agent joins with. Populated only for a PRIMARY `self_hosted` runtime; declared `sensitive` (Azure returns them readable and the provider does not flag them -- the module does)

## Behavior Notes

- **All 3 flavors share one name namespace** inside the factory, so switching variant blocks replaces the runtime (a different resource is created at the same ARM address).
- **The azure flavor's interactive authoring TTL travels on a separate Azure operation** (enable-interactive-authoring, applied after the runtime is online), not on the create body -- the provider handles the extra call, and a live debug session bills while enabled.
- **`virtual_network_enabled` requires the factory's managed virtual network**: Azure rejects the create with its own error when the factory does not have it enabled.
- **Creating an SSIS runtime leaves it STOPPED and unbilled** -- node-hours bill only after the runtime is started (an operational action in Data Factory Studio or via the Data Factory API).
- **SSIS `pipeline_external_compute_scale.number_of_external_nodes` never reads back**: Azure's read API mirrors the pipeline-node count for it (a provider seam both engines inherit); the value IS sent correctly on create/update.
- **Inline secrets win over Key Vault references** on the SSIS express custom setup's command keys and components -- when both are set, Azure receives the inline value (the provider's own precedence). Prefer the Key Vault forms.
- **Omitted optional arguments fall back to the provider's own defaults**: 8 cores / General compute on the azure flavor; 1 node / 1 parallel execution / Standard edition / LicenseIncluded on SSIS.
- **No tags**: integration runtimes are ARM sub-resources of the factory and expose none.
