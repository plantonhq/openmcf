# Overview

The **AzureAiFoundryProject** component creates an Azure AI Foundry project -- the workspace one AI team works in, created inside an AzureAiFoundry hub. The project inherits the hub's posture (key vault, storage, insights, registry, managed network, encryption) and carries only its own identity, naming, and description.

## Purpose

- **One project per team**: agents, evaluations, and model work organized per team, on a foundation the platform team governs once at the hub.
- **Per-team access scoping**: the project's own managed identity is what per-team data grants bind to.
- **Deliberately thin**: no companion-service fields -- everything infrastructural is inherited, so a team's manifest is a few lines.

## Key Features

- Full azurerm v5 surface: the typed hub reference (fixed at creation), optional system/user-assigned identity with the primary-identity selector (the provider's pairing enforced at manifest time), high-business-impact flag, descriptive surface.
- The resource-group inheritance recorded honestly: the project deploys into the HUB's resource group -- there is no resource-group field by design.

## Use Cases

- **Team workspaces**: fraud, support, and research teams each get a project inside the shared hub.
- **Environment separation**: dev and prod projects on differently-postured hubs.

## Future Enhancements

- Project connection resources (external data sources/models) as azurerm models them.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
