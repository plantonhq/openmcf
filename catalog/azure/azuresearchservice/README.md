# Overview

The **AzureSearchService** component creates an Azure AI Search service -- the managed search-and-retrieval engine AI applications use to index and query their own data (keyword, vector, and semantic search). AI Search is the standard retrieval companion to Azure OpenAI: the "R" in RAG.

## Purpose

- **Capacity as three explicit dials**: SKU (unit size and price class), partitions (storage and indexing), replicas (query throughput and availability) -- with the per-SKU caps validated at manifest time instead of failing at deploy.
- **Auth posture as configuration**: API keys only, RBAC alongside keys, or RBAC-only -- the provider's pairing rules front-loaded as validation.
- **Private data-source reach**: composed shared private links let indexers reach storage, SQL, or other sources sitting behind private endpoints.
- **Identity-first indexing**: the service's managed identity is how indexers reach data sources without connection strings.

## Key Features

- Full azurerm v5 surface: the seven-SKU vocabulary, replica/partition counts with every per-SKU cap as manifest-time CEL, high-density hosting (standard3), local-auth and Entra failure-mode pairing, CMK enforcement, public-access and IP-firewall controls with the AzureServices bypass, semantic ranking tiers, system/user-assigned identity.
- Composed `azurerm_search_shared_private_link_service` children, keyed by name into the `shared_private_link_service_ids` output, with the Pending-until-approved reality documented.
- The SKU update contract recorded where it bites: in-place upgrade ONLY along basic → standard → standard2 → standard3; every other change replaces the service.
- Service-minted credentials exported as sensitive outputs (admin keys, the built-in query key) -- masked by both engines.

## Use Cases

- **RAG retrieval**: the semantic-enabled search backend Azure OpenAI applications ground their answers on.
- **Application search**: product catalogs, documentation, in-app search.
- **Multi-tenant SaaS**: high-density hosting packs thousands of small indexes on one standard3 service.

## Future Enhancements

- Index/indexer management is data-plane (not ARM) -- stays with application tooling by design.
