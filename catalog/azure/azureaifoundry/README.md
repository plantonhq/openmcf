# Overview

The **AzureAiFoundry** component creates an Azure AI Foundry hub -- the shared foundation a company sets up ONCE for its AI teams. The hub owns the security and connectivity posture (identity, key vault, storage, optional customer-managed-key encryption, managed network, public access), and every AzureAiFoundryProject created inside it inherits that posture; teams then work in projects, not in the hub itself.

## Purpose

- **One governed foundation, many teams**: security decided once at the hub; each team gets a project without re-litigating vault, storage, network, or encryption.
- **Typed references end-to-end**: the resource group, key vault, storage account, insights, registry, and user-assigned identities all wire by reference -- chart-ready.
- **CMK and network isolation as first-class arms**: the regulated-estate posture (customer keys, approved-outbound managed network, private access) is declarative, not a portal afterthought.
- **Identity-first**: the hub reaches its companion services through its managed identity; grants compose before the hub exists when a user-assigned identity is used.

## Key Features

- Full azurerm v5 surface: required key vault + storage attachments (ForceNew), in-place-attachable Application Insights and container registry, system/user-assigned identity with the primary-identity selector, CMK encryption (versioned key URL -- the hub's own contract), managed-network isolation modes, public-access toggle, high-business-impact flag.
- The provider's contracts front-loaded as manifest-time validation: identity/identity-ids pairing, encryption block shape, the hub name's real code regex (underscores allowed despite the provider's error text).
- Soft-delete recorded where it bites: a deleted hub holds its name as a purgeable ghost (the ML workspace class).

## Use Cases

- **Company AI platform**: one hub; fraud, support, and research teams each get an AzureAiFoundryProject inside it.
- **Regulated estate**: the CMK-hardened hub with private access and approved-outbound isolation.
- **Azure OpenAI + Foundry**: the hub's projects consume model deployments from an AzureCognitiveAccount while AI Search carries retrieval.

## Future Enhancements

- Hub connection resources (external data sources/models) as azurerm models them.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
