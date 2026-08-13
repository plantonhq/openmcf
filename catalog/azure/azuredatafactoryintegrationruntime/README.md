# Overview

The **AzureDataFactoryIntegrationRuntime** component deploys one integration runtime inside an Azure Data Factory (AzureDataFactory) -- the compute engine the factory's pipelines, data flows, and copy activities actually run on. The factory stores definitions; the integration runtime executes them. One kind covers all three engine flavors azurerm models: the managed data-flow compute, the managed SSIS package runtime, and the self-hosted agent registration.

## Purpose

- **Compute as configuration**: which engine a factory runs on -- its size, network posture, and setup -- lives in the manifest, reviewed and versioned like everything else, not clicked together in a portal.
- **One kind, 3 flavors**: the variant block declares the engine flavor; Azure stores every flavor in one factory-scoped namespace, and so does the catalog.
- **The private-network bridge**: the self-hosted flavor registers the agent you install on your own machines, and Azure issues the authorization keys it joins with -- surfaced as sensitive outputs, ready for chart wiring.

## Key Features

- Full azurerm v5 surface across ALL THREE provider integration runtime resources: the data-flow compute's sizing and time-to-live knobs, SSIS at full depth (SSISDB catalog, custom setup script, express custom setup with Key-Vault-sourced credentials, both virtual network injection forms, package stores, both compute-scale blocks, the on-premises proxy), and the self-hosted registration with RBAC-shared linking.
- Chart-ready: `data_factory_id` defaults its reference to AzureDataFactory's ID output; the SSIS proxy and the self-hosted RBAC link reference OTHER integration runtimes through this kind's own outputs; every linked service reference defaults to AzureDataFactoryLinkedService's name output; network fields default to AzureVirtualNetwork, AzureSubnet, and AzurePublicIp outputs.
- Exact contracts: validation mirrors the provider's own rules -- each flavor's name format, the data-flow core-count and TTL menus, the SSIS node-size menu, the vnet-addressing exactly-one rule, the catalog tier/pool conflict, and the express-setup at-least-one rule.
- Secure by default: the SSIS secrets (catalog password, SAS token, command-key passwords, component licenses) are marked sensitive, each with a Key-Vault-reference alternative; the self-hosted authorization keys are sensitive OUTPUTS even though Azure returns them readable.

## Use Cases

- **Mapping data flows**: the azure flavor provisions the serverless Spark that data flows transform on -- inside the factory's managed virtual network when flows must reach private endpoints.
- **Lift-and-shift SSIS**: the azure_ssis flavor runs existing SSIS projects unchanged, with the SSISDB catalog on an Azure SQL server and node-level custom setup for drivers and components.
- **On-premises and cross-network data**: the self-hosted flavor bridges Data Factory to machines Azure cannot reach directly -- install the agent, hand it a key, and copy activities flow through it.
- **Sharing one bridge across factories**: a linked self-hosted registration (RBAC authorization) reuses another factory's runtime instead of installing a second agent fleet.

## Future Enhancements

- The Data Factory family's reference graph is complete with this kind; linked services gain their integration runtime reference wiring as the family's deferred upgrades land.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
