# Overview

The **AzureDataFactory** component deploys an Azure Data Factory -- the workspace every other Data Factory resource lives inside. Pipelines (AzureDataFactoryPipeline), data flows, linked services, datasets, triggers, and integration runtimes are all created against a factory's ARM ID; the factory itself carries the workspace-level posture: managed identity, git integration, global parameters, customer-managed-key encryption, named credentials, a managed virtual network, and managed private endpoints.

## Purpose

- **One workspace, many pipelines**: the factory is the container -- teams create pipelines against it (AzureDataFactoryPipeline) without touching the workspace posture.
- **Identity done once**: the factory's managed identity (system-assigned, user-assigned, or both) is what linked services authenticate with; named credentials wrap specific identities for per-connection control.
- **Private by construction**: the managed virtual network runs integration inside Azure-managed private networking, and managed private endpoints give it private egress to data stores.

## Key Features

- Full azurerm v5 surface: identity (including the combined mode), GitHub/Azure DevOps repository binding, global parameters, managed virtual network, public network posture, Purview binding, inline customer-managed-key encryption, both credential flavors, and managed private endpoints -- five provider resources behind one spec.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup, identity and credential IDs to AzureUserAssignedIdentity, the CMK key to AzureKeyVaultKey; the `data_factory_id` output is the wiring edge every Data Factory child references.
- Composed children surface as name-keyed outputs (`credential_ids`, `managed_private_endpoint_ids`) for downstream wiring.

## Use Cases

- **The data-platform workspace**: one factory per environment; pipelines, datasets, and linked services onboard against it.
- **Locked-down ETL**: managed VNet enabled, public access off, managed private endpoints to the lake and warehouse, CMK encryption from Key Vault.
- **Git-integrated authoring**: bind the factory to GitHub or Azure DevOps so the Studio authors against branches and publishes deliberately.

## Future Enhancements

- Pipeline wiring lives in AzureDataFactoryPipeline -- point its `data_factory_id` at this component's ID output. Data flows, linked services, datasets, triggers, and integration runtimes follow as their own kinds.
