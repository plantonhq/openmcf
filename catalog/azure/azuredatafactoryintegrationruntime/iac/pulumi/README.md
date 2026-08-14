# AzureDataFactoryIntegrationRuntime Pulumi Module

## Overview

Creates one integration runtime inside an Azure Data Factory -- the compute engine the factory's pipelines, data flows, and copy activities run on. The spec's variant block (exactly one of 3) selects the engine flavor, and the module dispatches to the matching Pulumi resource.

## Resources Created

Exactly one of:

- `datafactory.IntegrationRuntimeRule` -- the managed data-flow compute (`azure` variant)
- `datafactory.IntegrationRuntimeSsis` -- the managed SSIS package runtime (`azure_ssis` variant)
- `datafactory.IntegrationRuntimeSelfHosted` -- the self-hosted agent registration (`self_hosted` variant)

## Module Structure

- `module/main.go` -- provider setup and the 3-way variant dispatch
- `module/shared.go` -- the shared output contract and the send-only-when-set helpers
- `module/variants_azure.go` -- the managed data-flow compute builder
- `module/variants_ssis.go` -- the managed SSIS runtime builder (catalog, custom setup, network injection, package stores, compute scales, proxy)
- `module/variants_self_hosted.go` -- the self-hosted registration builder
- `module/outputs.go` -- exported output names

## Outputs

- `integration_runtime_id` -- the runtime's ARM resource ID (`{factory_id}/integrationRuntimes/{name}`, the same shape for all 3 flavors)
- `integration_runtime_name` -- the runtime's name (what linked services, data flow activities, and the SSIS proxy resolve against)
- `primary_authorization_key` / `secondary_authorization_key` -- the keys a self-hosted agent joins with. Populated only for a PRIMARY `self_hosted` runtime; exported as SECRETS (Azure returns them readable and the provider does not flag them -- the module does)

## Behavior Notes

- **All 3 flavors share one name namespace** inside the factory, so switching variant blocks replaces the runtime (a different resource is created at the same ARM address).
- **The azure flavor's interactive authoring TTL travels on a separate Azure operation** (enable-interactive-authoring, applied after the runtime is online), not on the create body -- the provider handles the extra call, and a live debug session bills while enabled.
- **`virtual_network_enabled` requires the factory's managed virtual network**: Azure rejects the create with its own error when the factory does not have it enabled.
- **Creating an SSIS runtime leaves it STOPPED and unbilled** -- node-hours bill only after the runtime is started (an operational action in Data Factory Studio or via the Data Factory API).
- **SSIS `pipeline_external_compute_scale.number_of_external_nodes` never reads back**: Azure's read API mirrors the pipeline-node count for it (a provider seam both engines inherit); the value IS sent correctly on create/update.
- **Inline secrets win over Key Vault references** on the SSIS express custom setup's command keys and components -- when both are set, Azure receives the inline value (the provider's own precedence). Prefer the Key Vault forms.
- **Omitted optional arguments fall back to the provider's own defaults**: 8 cores / General compute on the azure flavor; 1 node / 1 parallel execution / Standard edition / LicenseIncluded on SSIS.
- **No tags**: integration runtimes are ARM sub-resources of the factory and expose none.
