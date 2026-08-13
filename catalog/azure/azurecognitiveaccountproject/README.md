# Overview

The **AzureCognitiveAccountProject** component deploys an AI Foundry project onto an Azure AI services account -- the workspace a team organizes its AI work in: agents, evaluations, files and data-plane assets live inside a project, isolated from sibling projects on the same account.

## Purpose

- **Team isolation on shared capacity**: one `AIServices` account, many projects -- each team gets its own workspace without its own account, keys, or quota pool.
- **An identity per project**: every project carries a managed identity (the provider requires it) -- what its agents and evaluations act as, and what data-source grants bind to.
- **Chart-ready wiring**: a typed reference to the parent account; the project's data-plane endpoints surface as outputs.

## Key Features

- Full azurerm v5 surface: name, location, required identity (system/user-assigned), description, display name, tags.
- The empty-update quirk recorded where it bites: ARM cannot clear `description`/`displayName` in place -- clearing replaces the project.
- Outputs include the ARM-reported endpoints map and the is-default marker (the first project on an account becomes its default).

## Use Cases

- **Per-team AI Foundry workspaces** on a platform-owned account.
- **Environment separation** (dev/staging projects) without multiplying accounts.
- **Agent workloads** that need a scoped identity for storage and search grants.

## Future Enhancements

- Project-scoped connections once the account connections surface lands (deferred with the account's connections arm).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
